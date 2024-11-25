package processor

import (
	"context"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	lru "github.com/hashicorp/golang-lru"
	"github.com/metis-seq/themis/bridge/setu/rpc"
	"github.com/metis-seq/themis/bridge/setu/util"
	"github.com/metis-seq/themis/bridge/setu/util/sqlite"
	"github.com/metis-seq/themis/helper"
	metisTypes "github.com/metis-seq/themis/metis/types"
	themisTypes "github.com/metis-seq/themis/types"
)

const listenerMetisBlockKey = "listener-metischain-last-block"

// RpcTxProcessor
type RpcTxProcessor struct {
	BaseProcessor
	cancelService context.CancelFunc
	// cache
	cache *lru.Cache

	currentL2Block   uint64
	currentSequencer string
	currentEpoch     int64
	currentEpochInfo *epochInfo

	epochRotateExecuting bool
	epochRotated         int64
}

// NewRpcTxProcessor
func NewRpcTxProcessor() *RpcTxProcessor {
	rp := &RpcTxProcessor{}

	cache, _ := lru.New(10000)
	rp.cache = cache

	// set init rotated epoch to -1, so epoch 0 can rotate correctly
	rp.epochRotated = -1

	return rp
}

// Start starts sendRpcTx subscription
func (rp *RpcTxProcessor) Start() error {
	rp.Logger.Info("Starting")

	// create cancellable context
	rpCtx, cancelSpanService := context.WithCancel(context.Background())
	rp.cancelService = cancelSpanService

	// start polling for rpctx
	rp.Logger.Info("Start polling for broadcast tx to metis", "pollInterval", helper.GetConfig().MetisPollInterval)

	rp.fresh()

	go rp.startSendTxToThemisPolling(rpCtx)
	go rp.startSendTxToMetisPolling(rpCtx)
	go rp.freshEpochInfo(rpCtx)
	go rp.startSendRetryTx(rpCtx)

	return nil
}

func (rp *RpcTxProcessor) freshEpochInfo(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	// stop ticker when everything done
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rp.fresh()
		case <-ctx.Done():
			rp.Logger.Info("fresh epoch info stopped")
			ticker.Stop()
			return
		}
	}
}

func (rp *RpcTxProcessor) fresh() {
	var err error
	// check current sequencer
	rp.currentL2Block, err = getCurrentChildBlock(rp.contractConnector)
	if err != nil {
		rp.Logger.Error("Error getCurrentChildBlock", "error", err)
		return
	}

	detectSequencer, detectEpoch, err := getCurrentChildSequencer(rp.contractConnector, rp.cliCtx, uint64(rp.currentL2Block+1))
	if err != nil {
		rp.Logger.Error("Error getCurrentChildSequencer", "error", err)
		return
	}

	rp.Logger.Debug("rpInfo", "currentL2Block", rp.currentL2Block,
		"currentEpoch", detectEpoch, "currentSequencer", detectSequencer, "localSequencer", helper.GetAddressStr())

	// query epoch info
	detectEpochInfo, err := getChildEpochInfo(rp.contractConnector, rp.cliCtx, detectEpoch)
	if err != nil {
		rp.Logger.Error("Error getChildEpochInfo", "error", err)
		return
	}

	// check local handled l2 block
	listenerMetisBlockBytes, err := rp.storageClient.Get([]byte(listenerMetisBlockKey), nil)
	if err != nil {
		rp.Logger.Error("Error while fetching last block bytes from storage", "error", err)
		return
	}

	listenerMetisBlock, _ := strconv.ParseUint(string(listenerMetisBlockBytes), 10, 64)
	if listenerMetisBlock < uint64(detectEpochInfo.StartBlock)-1 {
		rp.Logger.Info("wait metis block to be handled to latest", "listenerMetisBlock", listenerMetisBlock, "epochStartBlock", detectEpochInfo.StartBlock)
		return
	}

	// sync lock with rpc_tx.go chan write
	util.MetisTxLock.Lock()
	defer util.MetisTxLock.Unlock()

	// update epoch and sequencer after validation and logic
	// this makes sure to remove old txs before broadcasting txs of db to metis chain
	rp.currentSequencer = detectSequencer
	rp.currentEpoch = detectEpoch
	rp.currentEpochInfo = detectEpochInfo

	// set old txs push true at the first time the seq switches
	if rp.currentL2Block+1 == uint64(detectEpochInfo.StartBlock) {
		if !rp.epochRotateExecuting && detectEpoch > rp.epochRotated {
			rp.Logger.Info("epoch rotate start", "unsavedChan", len(util.MetisTxChan), "currentL2Block", rp.currentL2Block,
				"currentEpoch", detectEpoch, "currentSequencer", detectSequencer, "localSequencer", helper.GetAddressStr())
			rp.epochRotateExecuting = true
			rp.SetOldTxsPushed(rp.currentL2Block)
			rp.epochRotateExecuting = false
			rp.epochRotated = detectEpoch
		} else {
			rp.Logger.Info("epoch rotate but should not execute again", "executing", rp.epochRotateExecuting, "detectEpoch", detectEpoch, "executedEpochRotated", rp.epochRotated)
		}
	}
}

// RegisterTasks - Registers rpcTx tasks with machinery
func (rp *RpcTxProcessor) RegisterTasks() {
}

func (rp *RpcTxProcessor) broadcastToThemis(logBytes string) error {
	encodedTx, err := hexutil.Decode(logBytes)
	if err != nil {
		rp.Logger.Error("Error while parsing rpcTx", "error", err)
		return err
	}

	tx := new(types.Transaction)
	if err := rlp.DecodeBytes(encodedTx, tx); err != nil {
		rp.Logger.Error("Error while parsing rpcTx", "error", err)
		return err
	}
	rp.Logger.Info(
		"✅ Received task to send rpcTx to themis",
		"txHash", tx.Hash().Hex(),
	)

	// broadcast to themis
	msg := metisTypes.NewMsgMetisTx(
		themisTypes.BytesToThemisAddress(helper.GetAddress()),
		themisTypes.HexToThemisHash(tx.Hash().Hex()),
		0,
		logBytes,
	)

	// broadcast to themis
	if err := rp.txBroadcaster.BroadcastToThemis(msg, nil); err != nil {
		rp.Logger.Error("Error while broadcasting rpcTx to themis", "txHash", tx.Hash().Hex(), "error", err)
		return err
	}

	return nil
}

func (rp *RpcTxProcessor) startSendTxToThemisPolling(ctx context.Context) {
	for {
		select {
		case txCache := <-util.MetisTxCacheChan:
			err := rp.broadcastToThemis(txCache.TxData)
			if err != nil {
				rp.Logger.Error("broadcastToThemis", "error", err)
			}
		case <-ctx.Done():
			rp.Logger.Info("Polling stopped")
			return
		}
	}
}

func (rp *RpcTxProcessor) startSendTxToMetisPolling(ctx context.Context) {
	for {
		select {
		case metisTx := <-util.MetisTxChan:
			rp.Logger.Info("consume MetisTxChan", "chanLen", len(util.MetisTxChan), "txHash", metisTx.TxHash)

			// insert to db
			saveID, _ := rp.sqlClient.BridgeSqliteMetisTx.Insert(metisTx)
			rp.Logger.Info("SendTxToMetis insert to db", "txHash", metisTx.TxHash, "recordID", saveID)

			if !strings.EqualFold(rp.currentSequencer, helper.GetAddressStr()) || !strings.EqualFold(rp.currentEpochInfo.Signer, helper.GetAddressStr()) {
				rp.Logger.Debug("not current sequencer, ignore it")
				rpc.HealthValue.IsCurrentSequencer = 0
				break
			}
			rpc.HealthValue.IsCurrentSequencer = 1

			// check local handled l2 block
			listenerMetisBlockBytes, err := rp.storageClient.Get([]byte(listenerMetisBlockKey), nil)
			if err != nil {
				rp.Logger.Error("Error while fetching last block bytes from storage", "error", err)
				break
			}

			listenerMetisBlock, _ := strconv.ParseUint(string(listenerMetisBlockBytes), 10, 64)
			if listenerMetisBlock < uint64(rp.currentEpochInfo.StartBlock)-1 {
				rp.Logger.Info("wait metis block to be handled to latest", "listenerMetisBlock", listenerMetisBlock, "epochStartBlock", rp.currentEpochInfo.StartBlock)
				break
			}

			go func(mTx *sqlite.MetisTx) {
				rp.Logger.Info("SendTxToMetis go routine get metis block ", "txHash", mTx.TxHash, "listenerMetisBlock", listenerMetisBlock)

				if metisTx.TxData != "" {
					rp.Logger.Info("SendTxToMetis go routine build tx ", "txHash", mTx.TxHash)
					// build tx
					encodedTx, err := hexutil.Decode(metisTx.TxData)
					if err != nil {
						rp.Logger.Error("Error broadcastToMetis while parsing rpcTx", "txHash", metisTx.TxHash, "txData", metisTx.TxData, "error", err)
						return
					}

					signedTx := new(types.Transaction)
					if err := rlp.DecodeBytes(encodedTx, signedTx); err != nil {
						rp.Logger.Error("Error broadcastToMetis while parsing rpcTx", "error", err)
						return
					}

					rp.Logger.Info("SendTxToMetis go routine build tx success", "txHash", mTx.TxHash)

					// set health value
					rpc.HealthValue.Writer.TxHash = signedTx.Hash().Hex()
					rpc.HealthValue.Writer.Timestamp = time.Now().Unix()

					// broadcast to metis
					rp.Logger.Info("SendTxToMetis go routine build start BroadcastRawTxToMetis", "txHash", mTx.TxHash)

					if err := rp.txBroadcaster.BroadcastRawTxToMetis(signedTx); err != nil {
						rpc.HealthValue.Writer.Ret = err.Error()
						rp.Logger.Error("Error broadcastToMetis while broadcasting rpcTx to metis", "txHash", signedTx.Hash().Hex(), "error", err)

						// need resend tx
						rp.sqlClient.BridgeSqliteMetisTxRetry.Insert(metisTx.TxData)

						return
					}

					rp.Logger.Info("✅ Success broadcastToMetis tx to metis", "txHash", signedTx.Hash().Hex())
				}
			}(metisTx)
		case <-ctx.Done():
			rp.Logger.Info("Polling stopped")
			return
		}
	}
}

func (rp *RpcTxProcessor) SetOldTxsPushed(height uint64) {
	rp.Logger.Info("SetOldTxsPushed in", "block_height", height)

	client := helper.GetMetisClient()
	// find 100 blocks with multiple transactions
	maxFound := 100
	var foundID uint64

	for i := height; i > 0; i-- {
		block, err := client.BlockByNumber(context.TODO(), big.NewInt(int64(i)))
		if err != nil {
			rp.Logger.Error("SetOldTxsPushed BlockByNumber", "Error", err)
			continue
		}

		transactions := block.Transactions()
		if len(transactions) <= 0 {
			continue
		}
		// L2 block contains multiple transactions
		for j := len(transactions) - 1; j >= 0; j-- {
			txHash := transactions[j].Hash().Hex()
			txInDB, _ := rp.sqlClient.BridgeSqliteMetisTx.GetMetisTxByTxHash(txHash)
			if txInDB != nil {
				foundID = txInDB.ID
				rp.Logger.Info("SetOldTxsPushed found tx", "foundTxHash", txHash, "foundID", foundID)
				break // stop for range
			}
		}
		if foundID > 0 {
			break
		}

		// reach max found
		if maxFound == 0 {
			break
		}
		maxFound--
	}

	if foundID > 0 {
		// write all new txs
		metisTxs, err := rp.sqlClient.BridgeSqliteMetisTx.GetAllWaitPushMetisTxsByStartID(int64(foundID))
		if err != nil {
			rp.Logger.Error("SetOldTxsPushed GetAllWaitPushMetisTxsByStartID", "err", err)
			return
		}
		for _, metisTx := range metisTxs {
			// delete old tx
			rp.sqlClient.BridgeSqliteMetisTx.DeleteByID(metisTx.ID)
			// resend tx
			util.MetisTxChan <- metisTx
			rp.Logger.Info("SetOldTxsPushed load history metis tx", "txHash", metisTx.TxHash, "recordID", metisTx.ID)
		}

		// delete
		rp.sqlClient.BridgeSqliteMetisTx.DeleteExpiredDataByID(foundID)
	}
}

func (rp *RpcTxProcessor) startSendRetryTx(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	// stop ticker when everything done
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			retryTxs, err := rp.sqlClient.BridgeSqliteMetisTxRetry.GetAllWaitPushMetisTxRetrys(100, 0)
			if err != nil {
				rp.Logger.Info("GetAllWaitPushMetisTxRetrys", "err", err)
				break
			}

			for _, retryTx := range retryTxs {
				err = rp.broadcastToThemis(retryTx.TxData)
				if err != nil {
					rp.Logger.Info("startSendRetryTx broadcastToThemis", "err", err)
					continue
				}
				rp.sqlClient.BridgeSqliteMetisTxRetry.Delete(retryTx.ID)
			}
		case <-ctx.Done():
			rp.Logger.Info("fresh epoch info stopped")
			ticker.Stop()
			return
		}
	}
}

// Stop stops all necessary go routines
func (rp *RpcTxProcessor) Stop() {
	// cancel  polling
	rp.cancelService()
}

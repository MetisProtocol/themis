package listener

import (
	"context"
	"encoding/hex"
	"math/big"
	"strconv"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	jsoniter "github.com/json-iterator/go"

	"github.com/metis-seq/themis/bridge/setu/util"
	"github.com/metis-seq/themis/bridge/setu/util/sqlite"
	"github.com/metis-seq/themis/helper"
	metisTypes "github.com/metis-seq/themis/metis/types"
	mpcTypes "github.com/metis-seq/themis/mpc/types"
	"github.com/metis-seq/themis/types"
)

const (
	themisLastBlockKey = "listener-themis-last-block" // storage key
)

// ThemisListener - Listens to and process events from themis
type ThemisListener struct {
	BaseListener
}

// NewThemisListener - constructor func
func NewThemisListener() *ThemisListener {
	return &ThemisListener{}
}

// Start starts new block subscription
func (hl *ThemisListener) Start() error {
	hl.Logger.Info("Starting")

	// create cancellable context
	headerCtx, cancelHeaderProcess := context.WithCancel(context.Background())
	hl.cancelHeaderProcess = cancelHeaderProcess

	// Themis pollIntervall = (minimal pollInterval of themis chain)
	pollInterval := helper.GetConfig().ThemisPollInterval

	hl.Logger.Info("Start polling for themis events", "pollInterval", pollInterval)
	hl.StartPolling(headerCtx, pollInterval, nil)

	return nil
}

// ProcessHeader -
func (hl *ThemisListener) ProcessHeader(_ *blockHeader) {

}

// StartPolling - starts polling for themis events
func (hl *ThemisListener) StartPolling(ctx context.Context, pollInterval time.Duration, _ *big.Int) {
	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(pollInterval)

	// ticker := time.NewTicker(500 * time.Millisecond)
	// start listening
	for {
		select {
		case <-ticker.C:
			begin := time.Now()
			fromBlock, toBlock, err := hl.fetchFromAndToBlock()

			if err != nil {
				hl.Logger.Error("Error fetching from and toBlock, skipping events query", "fromBlock", fromBlock, "toBlock", toBlock, "error", err)
			} else if fromBlock < toBlock {
				hl.Logger.Info("StartPolling Fetching new events between", "fromBlock", fromBlock, "toBlock", toBlock)

				// Querying and processing Begin events
				for i := fromBlock; i <= toBlock; i++ {
					events, err := helper.GetBeginBlockEvents(hl.httpClient, int64(i))
					if err != nil {
						hl.Logger.Error("Error fetching begin block events", "error", err)
					}

					for _, event := range events {
						hl.ProcessBlockEvent(sdk.StringifyEvent(event), int64(i))
					}

					// parse deliver tx event
					// hl.ProcessBlockDeliverTxEvent(int64(i))
				}

				// set last block to storage
				if err := hl.storageClient.Put([]byte(themisLastBlockKey), []byte(strconv.FormatUint(toBlock, 10)), nil); err != nil {
					hl.Logger.Error("hl.storageClient.Put", "Error", err)
				}
				hl.Logger.Debug("StartPolling put themisLastBlock", "cost", time.Since(begin).String())
			}
		case <-ctx.Done():
			hl.Logger.Info("Polling stopped")
			ticker.Stop()

			return
		}
	}
}

func (hl *ThemisListener) fetchFromAndToBlock() (uint64, uint64, error) {
	// toBlock - get latest blockheight from themis node
	fromBlock := uint64(0)
	toBlock := uint64(0)

	nodeStatus, err := helper.GetNodeStatus(hl.cliCtx)
	if err != nil {
		hl.Logger.Error("Error while fetching themis node status", "error", err)
		return fromBlock, toBlock, err
	}

	toBlock = uint64(nodeStatus.SyncInfo.LatestBlockHeight)
	fromBlock = uint64(nodeStatus.SyncInfo.LatestBlockHeight) - 1000

	// fromBlock - get last block from storage
	hasLastBlock, _ := hl.storageClient.Has([]byte(themisLastBlockKey), nil)
	if hasLastBlock {
		lastBlockBytes, err := hl.storageClient.Get([]byte(themisLastBlockKey), nil)
		if err != nil {
			hl.Logger.Info("Error while fetching last block bytes from storage", "error", err)
			return fromBlock, toBlock, err
		}

		if result, err := strconv.ParseUint(string(lastBlockBytes), 10, 64); err == nil {
			hl.Logger.Debug("Got themis last block from bridge storage", "lastBlock", result)
			fromBlock = result + 1
		} else {
			hl.Logger.Info("Error parsing last block bytes from storage", "error", err)
			toBlock = 0

			return fromBlock, toBlock, err
		}
	}

	return fromBlock, toBlock, err
}

// ProcessBlockEvent - process Blockevents (BeginBlock, EndBlock events) from themis.
func (hl *ThemisListener) ProcessBlockEvent(event sdk.StringEvent, blockHeight int64) {
	hl.Logger.Info("Received block event from Themis", "eventType", event.Type)

	eventBytes, err := jsoniter.ConfigFastest.Marshal(event)
	if err != nil {
		hl.Logger.Error("Error while parsing block event", "eventType", event.Type, "error", err)
		return
	}
	hl.Logger.Debug("Event bytes", "eventType", event.Type, "data", eventBytes)

	switch event.Type {
	case mpcTypes.EventTypeProposeMpcSign, metisTypes.EventTypeProposeSpan, metisTypes.EventTypeReProposeSpan:
		hl.insertToDb(event.Type, string(eventBytes))
	case metisTypes.EventTypeMetisTx:
		hl.saveToMetisTxEvent(string(eventBytes))

	default:
		hl.Logger.Debug("BlockEvent Type mismatch", "eventType", event.Type)
	}
}

func (hl *ThemisListener) ProcessBlockDeliverTxEvent(blockHeight int64) error {
	hl.Logger.Info("ProcessBlockDeliverTxEvent", "height", blockHeight)

	block, err := helper.GetBlock(hl.cliCtx, blockHeight)
	if err != nil {
		return err
	}

	for _, tx := range block.Block.Txs {
		txHash := types.BytesToThemisHash(tx.Hash()).Hex()
		hl.Logger.Info("ProcessBlockDeliverTxEvent", "tx_hash", txHash)

		queryTxHash := strings.TrimPrefix(txHash, "0x")
		txRes, err := helper.QueryTx(hl.cliCtx, queryTxHash)
		if err != nil {
			hl.Logger.Error("ProcessBlockDeliverTxEvent QueryTx", "err", err)
			return err
		}

		for _, event := range txRes.Events {
			eventBytes, err := jsoniter.ConfigFastest.Marshal(event)
			if err != nil {
				hl.Logger.Error("Error while parsing block event", "eventType", event.Type, "error", err)
				return err
			}
			hl.Logger.Debug("Event bytes", "eventType", event.Type, "data", eventBytes)
			hl.insertToDb(event.Type, hex.EncodeToString(eventBytes))
		}
	}
	return nil
}

func (hl *ThemisListener) insertToDb(eventName, eventLog string) error {
	_, err := hl.sqlClient.BridgeSqliteThemisEvent.Insert(eventName, eventLog)
	return err
}

func (hl *ThemisListener) saveToMetisTxEvent(eventBytes string) error {
	hl.Logger.Debug("Received broadcastToMetis request", "eventBytes", eventBytes)

	var event sdk.StringEvent
	if err := jsoniter.ConfigFastest.Unmarshal([]byte(eventBytes), &event); err != nil {
		hl.Logger.Error("Error broadcastToMetis unmarshalling event from themis", "error", err)
		return err
	}
	hl.Logger.Debug("processing broadcastToMetis confirmation event", "eventtype", event.Type)

	var txHash string
	var txData string

	for _, attr := range event.Attributes {
		if attr.Key == metisTypes.AttributeKeyMetisTxHash {
			txHash = attr.Value
		}

		if attr.Key == metisTypes.AttributeKeyMetisTxData {
			txData = attr.Value
		}
	}

	hl.Logger.Info(
		"✅ Received broadcastToMetis task to send rpcTx to metis",
		"txHash", txHash,
	)

	if txHash == "" || txData == "" {
		hl.Logger.Info("processing broadcastToMetis invalid tx info, ignore it", "txHash", txHash)
		return nil
	}

	// sync lock with rpc_tx.go chan write
	util.MetisTxLock.Lock()
	defer util.MetisTxLock.Unlock()

	// check tx record
	record, _ := hl.sqlClient.BridgeSqliteMetisTx.GetMetisTxByTxHash(txHash)
	if record != nil {
		hl.sqlClient.BridgeSqliteMetisTx.DeleteByID(record.ID)
		hl.Logger.Info("saveToMetisTxEvent tx already exist, delete it", "txHash", txHash)
	}

	// send to chan
	util.MetisTxChan <- &sqlite.MetisTx{
		TxHash: txHash,
		TxData: txData,
		Pushed: false,
		Mined:  false,
	}

	hl.Logger.Info("Success send metis tx to chan", "txHash", txHash, "chanLen", len(util.MetisTxChan))
	return nil
}

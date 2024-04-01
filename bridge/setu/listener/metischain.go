package listener

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/metis-seq/themis/bridge/setu/rpc"
	"github.com/metis-seq/themis/helper"
)

// MetisListener - Listens to and process events from metischain
type MetisListener struct {
	BaseListener
}

const (
	lastMetisBlockKey = "listener-metischain-last-block" // storage key
)

// NewMetisListener - constructor func
func NewMetisListener() *MetisListener {
	return &MetisListener{}
}

// Start starts new block subscription
func (ml *MetisListener) Start() error {
	ml.Logger.Info("Starting")

	// create cancellable context
	ctx, cancelSubscription := context.WithCancel(context.Background())
	ml.cancelSubscription = cancelSubscription

	// start go routine to poll for new header using client object
	ml.Logger.Info("Start polling for metischain blocks", "pollInterval", helper.GetConfig().MetisPollInterval)

	// start polling blocks in metis chain (available post-merge)
	go ml.StartPolling(ctx, helper.GetConfig().MetisPollInterval, nil)
	return nil
}

// ProcessHeader - process header block from metis chain
func (ml *MetisListener) ProcessHeader(newHeader *blockHeader) {
}

// StartPolling - starts polling for metis events
func (ml *MetisListener) StartPolling(ctx context.Context, pollInterval time.Duration, _ *big.Int) {
	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(pollInterval)

	// start listening
	for {
		select {
		case <-ticker.C:
			// get last block from metis chain
			chainHeight, err := ml.chainClient.BlockNumber(context.TODO())
			if err != nil {
				ml.Logger.Error("Error while fetching BlockNumber", "error", err)
				break
			}
			ml.Logger.Debug("Got metis last block from chain", "chainHeight", chainHeight)

			// get last block from storage
			hasLastBlock, _ := ml.storageClient.Has([]byte(lastMetisBlockKey), nil)
			if hasLastBlock {
				lastBlockBytes, err := ml.storageClient.Get([]byte(lastMetisBlockKey), nil)
				if err != nil {
					ml.Logger.Error("Error while fetching last block bytes from storage", "error", err)
					break
				}
				ml.Logger.Debug("Got metis last block from bridge storage", "lastBlock", string(lastBlockBytes))

				dbHeight, _ := strconv.ParseUint(string(lastBlockBytes), 10, 64)
				if dbHeight >= chainHeight {
					ml.Logger.Debug("ml handle metis block, wait new block", "dbHeight", dbHeight, "chainHeight", chainHeight)
					break
				}

				// handle block
				ml.handleBlock(big.NewInt(int64(dbHeight+1)), big.NewInt(int64(chainHeight)))
			}

			// set health value
			rpc.HealthValue.L2.BlockNumber = chainHeight
			rpc.HealthValue.L2.Timestamp = time.Now().Unix()

			// Set last block to storage
			if err := ml.storageClient.Put([]byte(lastMetisBlockKey), []byte(fmt.Sprintf("%v", chainHeight)), nil); err != nil {
				ml.Logger.Error("ml.storageClient.Put", "Error", err)
			}
		case <-ctx.Done():
			ml.Logger.Info("Polling stopped")
			ticker.Stop()
			return
		}
	}
}

func (ml *MetisListener) handleBlock(from, to *big.Int) {
	for i := from; i.Cmp(to) <= 0; {
		ml.Logger.Debug("ml.handleBlock", "block", i.String())

		rootCtx := context.Background()
		ctx, cancel := context.WithTimeout(rootCtx, 10*time.Second)
		defer cancel()

		// get tx
		block, err := ml.chainClient.BlockByNumber(ctx, i)
		if err != nil {
			ml.Logger.Error("ml.BlockByNumber", "Error", err)
			return
		}

		for _, tx := range block.Transactions() {
			ml.Logger.Debug("ml received metis tx", "tx_hash", tx.Hash().Hex())
			// record, _ := ml.sqlClient.BridgeSqliteMetisTx.GetMetisTxByTxHash(tx.Hash().Hex())
			// if record != nil {
			// 	ml.Logger.Info("ml delete metis tx id", "tx_hash", tx.Hash().Hex(), "id", record.ID)
			// 	// ml.sqlClient.BridgeSqliteMetisTx.DeleteExpiredDataByID(record.ID)
			// }
		}

		i = big.NewInt(0).Add(i, big.NewInt(1)) // i++
	}
}

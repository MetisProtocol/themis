package processor

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"

	"github.com/metis-seq/themis/bridge/setu/util"
	"github.com/metis-seq/themis/helper"
	metisTypes "github.com/metis-seq/themis/metis/types"
	mpcTypes "github.com/metis-seq/themis/mpc/types"
	"github.com/metis-seq/themis/types"
)

const (
	lastMetisBlockKey   = "processor-metis-last-block"     // storage key
	localSpanProposeKey = "processor-local-span-propose"   // storage key
	localRespanCountKey = "processor-local-count-times-%v" // storage key
)

// SpanProcessor - process span related events
type SpanProcessor struct {
	BaseProcessor

	// header listener subscription
	cancelSpanService context.CancelFunc

	caller helper.ContractCaller
	lock   sync.Mutex

	recoverOnce sync.Once
}

// Start starts new block subscription
func (sp *SpanProcessor) Start() error {
	sp.Logger.Info("Starting")

	// create cancellable context
	spanCtx, cancelSpanService := context.WithCancel(context.Background())

	sp.cancelSpanService = cancelSpanService

	// start polling for span
	sp.Logger.Info("Start polling for span", "pollInterval", helper.GetConfig().SpanPollInterval)

	go sp.startPolling(spanCtx, helper.GetConfig().SpanPollInterval)

	// start watching for re-span
	sp.Logger.Info("Start checking for re-span", "pollInterval", helper.GetConfig().ReSpanPollInterval, "delayTime", helper.GetConfig().ReSpanDelayTime)
	go sp.startCheckRespan(spanCtx, helper.GetConfig().ReSpanPollInterval, helper.GetConfig().ReSpanDelayTime)

	go sp.startRecoverPrevSpanTimer()

	return nil
}

func (sp *SpanProcessor) startRecoverPrevSpanTimer() {
	time.AfterFunc(5*time.Minute, func() {
		sp.recoverOnce.Do(func() {
			sp.recoverPrevSpan()
		})
	})
}

func (sp *SpanProcessor) recoverPrevSpan() {
	sp.Logger.Info("Span recover prev span check")
	lastSpan, err := sp.getLastSpan()
	if err != nil {
		sp.Logger.Error("Unable to fetch last span", "error", err)
		return
	}

	if lastSpan == nil {
		sp.Logger.Error("Last span is nil")
		return
	}
	// get span info from L2
	l2EpochNumber, _, l2EpochEnd, err := getChildLatestEpoch(sp.contractConnector, sp.cliCtx)
	if err != nil {
		sp.Logger.Error("Unable to fetch l2 span", "error", err)
		return
	}
	if l2EpochNumber+1 == int64(lastSpan.ID) && uint64(l2EpochEnd+1) < lastSpan.StartBlock && lastSpan.SelectedProducers[0].Signer.EthAddress() == common.HexToAddress(helper.GetAddressStr()) {
		// l2 epoch large than pos span,recover span list from L2
		prevSpan, err := getChildEpochInfo(sp.contractConnector, sp.cliCtx, l2EpochNumber)
		if err != nil {
			sp.Logger.Error("Unable to getChildEpochInfo for prev span", "l2EpochNumber", l2EpochNumber, "error", err)
			return
		}
		sp.proposeRecoverSpan(prevSpan, lastSpan.ChainID)
		sp.Logger.Info("Propose recover prev span", "l2EpochNumber", l2EpochNumber, "l2EpochEnd", l2EpochEnd, "lastSpan.ID", lastSpan.ID, "lastSpan.StartBlock", lastSpan.StartBlock)
	} else {
		sp.Logger.Info("Not propose recover prev span", "l2EpochNumber", l2EpochNumber, "l2EpochEnd", l2EpochEnd, "lastSpan.ID", lastSpan.ID, "lastSpan.StartBlock", lastSpan.StartBlock)
	}
}

// RegisterTasks - nil
func (sp *SpanProcessor) RegisterTasks() {
}

// startPolling - polls themis and checks if new span needs to be proposed
func (sp *SpanProcessor) startPolling(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	dbTicker := time.NewTicker(10 * time.Millisecond)
	// stop ticker when everything done
	defer ticker.Stop()
	defer dbTicker.Stop()

	for {
		select {
		case <-dbTicker.C:
			sp.checkEventForProposeSpan()
			sp.checkEventForProposeReSpan()
		case <-ticker.C:
			sp.checkAndPropose()
		case <-ctx.Done():
			sp.Logger.Info("Polling stopped")
			ticker.Stop()

			return
		}
	}
}

// startCheckRespan - polls themis and checks if new span needs to be proposed
func (sp *SpanProcessor) startCheckRespan(ctx context.Context, interval, delayTime time.Duration) {
	ticker := time.NewTicker(interval) // metis block time

	// stop ticker when everything done
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sp.checkReSpan(delayTime)
		case <-ctx.Done():
			sp.Logger.Info("Polling stopped")
			ticker.Stop()

			return
		}
	}
}

// checkAndPropose - will check if current user is span proposer and proposes the span
func (sp *SpanProcessor) checkAndPropose() {
	lastSpan, err := sp.getLastSpan()
	if err != nil {
		sp.Logger.Error("Unable to fetch last span", "error", err)
		return
	}

	if lastSpan == nil {
		return
	}
	sp.Logger.Debug("span checkAndPropose Found last span", "lastSpan", lastSpan.ID, "startBlock", lastSpan.StartBlock, "endBlock", lastSpan.EndBlock)

	epochLength, err := getSequencerEpochLength(sp.contractConnector, sp.cliCtx)
	if err != nil {
		sp.Logger.Error("Unable to fetch current epoch length", "error", err)
		return
	}

	nextSpanMsg, err := sp.fetchNextSpanDetails(lastSpan.ID+1, lastSpan.EndBlock+1, uint64(epochLength))
	if err != nil {
		sp.Logger.Error("Unable to fetch next span details", "error", err, "lastSpanId", lastSpan.ID)
		return
	}
	sp.Logger.Debug("span checkAndPropose Fetch next span msg", "nextSpan", nextSpanMsg.String())

	// get span info from L2
	l2EpochNumber, _, _, err := getChildLatestEpoch(sp.contractConnector, sp.cliCtx)
	if err != nil {
		sp.Logger.Error("Unable to fetch l2 span", "error", err)
		return
	}
	sp.Logger.Debug("span propose", "l2EpochNumber", l2EpochNumber)

	if l2EpochNumber > int64(lastSpan.ID) {
		// l2 epoch large than pos span,recover span list from L2
		newSpan, err := getChildEpochInfo(sp.contractConnector, sp.cliCtx, int64(lastSpan.ID)+1)
		if err != nil {
			sp.Logger.Error("Unable to getChildEpochInfo", "error", err)
			return
		}
		sp.proposeRecoverSpan(newSpan, lastSpan.ChainID)
		return
	}

	util.RecoverSpanFinished = true
	sp.propose(lastSpan, nextSpanMsg)
}

// propose producers for next span if needed
func (sp *SpanProcessor) propose(lastSpan *types.Span, nextSpanMsg *types.Span) {
	// call with last span on record + new span duration and see if it has been proposed
	currentBlock, err := getCurrentChildBlock(sp.contractConnector)
	if err != nil {
		sp.Logger.Error("Unable to fetch current block", "error", err)
		return
	}

	// get current sequencer epoch
	_, currentL2Epoch, err := getCurrentChildSequencer(sp.contractConnector, sp.cliCtx, currentBlock)
	if err != nil {
		sp.Logger.Error("Error while getCurrentChildSequencer", "error", err)
		return
	}

	if currentL2Epoch+1 != int64(nextSpanMsg.ID) {
		sp.Logger.Info("Unable to propose new span", "nextSpanMsg.ID", nextSpanMsg.ID, "currentL2Epoch", currentL2Epoch)
		return
	}

	// Delay the span's propose time to ensure that the first initialization has enough time to complete
	spanBlock := lastSpan.StartBlock + (lastSpan.EndBlock-lastSpan.StartBlock+1)/10
	sp.Logger.Info("span propose", "lastSpanStartBlock", lastSpan.StartBlock, "lastSpanEndBlock", lastSpan.EndBlock, "currentBlock", currentBlock, "spanBlock", spanBlock)

	if spanBlock <= currentBlock && currentBlock <= lastSpan.EndBlock {
		sp.lock.Lock()
		defer sp.lock.Unlock()

		// log new span
		sp.Logger.Info("✅ Proposing new span", "spanId", nextSpanMsg.ID, "startBlock", nextSpanMsg.StartBlock, "endBlock", nextSpanMsg.EndBlock)

		seed, err := sp.fetchNextSpanSeed()
		if err != nil {
			sp.Logger.Info("Error while fetching next span seed from ThemisServer", "err", err)
			return
		}

		// broadcast to themis
		msg := metisTypes.MsgProposeSpan{
			ID:              nextSpanMsg.ID,
			Proposer:        types.BytesToThemisAddress(helper.GetAddress()),
			CurrentL2Height: currentBlock,
			StartBlock:      nextSpanMsg.StartBlock,
			EndBlock:        nextSpanMsg.EndBlock,
			ChainID:         nextSpanMsg.ChainID,
			Seed:            seed,
			BlackList:       nil,
		}
		sp.Logger.Info("new span", "msg", msg)

		// return broadcast to themis
		if err := sp.txBroadcaster.BroadcastToThemis(msg, nil); err != nil {
			sp.Logger.Error("Error while broadcasting span to themis", "spanId", nextSpanMsg.ID, "startBlock", nextSpanMsg.StartBlock, "endBlock", nextSpanMsg.EndBlock, "error", err)
			return
		}
	}
}

func (sp *SpanProcessor) proposeRecoverSpan(newSpan *epochInfo, chainID string) {
	// log new span
	sp.Logger.Info("✅ Proposing recover new span", "spanId", newSpan.ID, "startBlock", newSpan.StartBlock, "endBlock", newSpan.EndBlock)

	seed, err := sp.fetchNextSpanSeed()
	if err != nil {
		sp.Logger.Info("Error while fetching next span seed from ThemisServer", "err", err)
		return
	}

	// broadcast to themis
	msg := metisTypes.MsgProposeSpan{
		ID:            uint64(newSpan.ID),
		Proposer:      types.BytesToThemisAddress(helper.GetAddress()),
		StartBlock:    uint64(newSpan.StartBlock),
		EndBlock:      uint64(newSpan.EndBlock),
		ChainID:       chainID,
		Seed:          seed,
		IsRecover:     true,
		RecoverSigner: types.HexToThemisAddress(newSpan.Signer),
	}

	// return broadcast to themis
	if err := sp.txBroadcaster.BroadcastToThemis(msg, nil); err != nil {
		sp.Logger.Error("Error while broadcasting span to themis", "spanId", newSpan.ID, "startBlock", newSpan.StartBlock, "endBlock", newSpan.EndBlock, "error", err)
		return
	}

	sp.setLocalSpanPropose(uint64(newSpan.ID))
}

func (sp *SpanProcessor) checkEventForProposeSpan() {
	// need to broadcast tx to metis
	allEvents, _ := sp.sqlClient.BridgeSqliteThemisEvent.GetAllWaitPushThemisEventsByType(metisTypes.EventTypeProposeSpan, 100, 0)
	for _, event := range allEvents {
		err := sp.l2CommitSpan(event.EventLog)
		if err == nil {
			sp.sqlClient.BridgeSqliteThemisEvent.Delete(event.ID)
		}
	}
}

func (sp *SpanProcessor) checkEventForProposeReSpan() {
	// need to broadcast tx to metis
	allEvents, _ := sp.sqlClient.BridgeSqliteThemisEvent.GetAllWaitPushThemisEventsByType(metisTypes.EventTypeReProposeSpan, 100, 0)
	for _, event := range allEvents {
		err := sp.l2RecommitSpan(event.EventLog)
		if err == nil {
			sp.sqlClient.BridgeSqliteThemisEvent.Delete(event.ID)
		}
	}
}

func (sp *SpanProcessor) setLocalSpanPropose(spanID uint64) {
	if err := sp.storageClient.Put([]byte(localSpanProposeKey+fmt.Sprintf("%v", spanID)), []byte(fmt.Sprintf("%v", spanID)), nil); err != nil {
		sp.Logger.Error("rl.storageClient.Put", "Error", err)
	}
}

// checks span status
func (sp *SpanProcessor) getLastSpan() (*types.Span, error) {
	// fetch latest start block from themis via rest query
	result, err := helper.FetchFromAPI(sp.cliCtx, helper.GetThemisServerEndpoint(util.LatestSpanURL))
	if err != nil {
		sp.Logger.Error("Error while fetching latest span")
		return nil, err
	}

	var lastSpan types.Span
	if err = jsoniter.ConfigFastest.Unmarshal(result.Result, &lastSpan); err != nil {
		sp.Logger.Error("Error unmarshalling span", "error", err)
		return nil, err
	}
	return &lastSpan, nil
}

func (sp *SpanProcessor) getSpan(id uint64) (*types.Span, error) {
	// fetch latest start block from themis via rest query
	result, err := helper.FetchFromAPI(sp.cliCtx, helper.GetThemisServerEndpoint(fmt.Sprintf(util.SpanByIdURL, id)))
	if err != nil {
		sp.Logger.Error("Error while fetching latest span")
		return nil, err
	}

	var lastSpan types.Span
	if err = jsoniter.ConfigFastest.Unmarshal(result.Result, &lastSpan); err != nil {
		sp.Logger.Error("Error unmarshalling span", "error", err)
		return nil, err
	}
	return &lastSpan, nil
}

// isValidator checks if current user is a validator
func (sp *SpanProcessor) isValidator() bool {
	valSet, err := util.GetValidatorSet(sp.cliCtx)
	if err != nil {
		sp.Logger.Error("Error GetValidatorSet", "error", err)
		return false
	}

	for _, val := range valSet.Validators {
		if val.Signer.EthAddress().Hex() == common.HexToAddress(helper.GetAddressStr()).Hex() {
			return true
		}
	}

	return false
}

// fetch next span details from themis.
func (sp *SpanProcessor) fetchNextSpanDetails(id, start, epochLength uint64) (*types.Span, error) {
	req, err := http.NewRequest("GET", helper.GetThemisServerEndpoint(util.NextSpanInfoURL), nil)
	if err != nil {
		sp.Logger.Error("Error creating a new request", "error", err)
		return nil, err
	}

	configParams, err := util.GetChainmanagerParams(sp.cliCtx)
	if err != nil {
		sp.Logger.Error("Error while fetching chainmanager params", "error", err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("span_id", strconv.FormatUint(id, 10))
	q.Add("start_block", strconv.FormatUint(start, 10))
	q.Add("chain_id", configParams.ChainParams.MetisChainID)
	q.Add("proposer", helper.GetFromAddress(sp.cliCtx).String())
	req.URL.RawQuery = q.Encode()

	// fetch next span details
	result, err := helper.FetchFromAPI(sp.cliCtx, req.URL.String())
	if err != nil {
		sp.Logger.Error("Error fetching proposers", "error", err)
		return nil, err
	}

	var msg types.Span
	if err = jsoniter.ConfigFastest.Unmarshal(result.Result, &msg); err != nil {
		sp.Logger.Error("Error unmarshalling propose tx msg ", "error", err)
		return nil, err
	}

	msg.EndBlock = msg.StartBlock + epochLength - 1
	sp.Logger.Debug("Generated proposer span msg", "msg", msg.String())

	return &msg, nil
}

// fetchNextSpanSeed - fetches seed for next span
func (sp *SpanProcessor) fetchNextSpanSeed() (nextSpanSeed common.Hash, err error) {
	sp.Logger.Info("Sending Rest call to Get Seed for next span")

	response, err := helper.FetchFromAPI(sp.cliCtx, helper.GetThemisServerEndpoint(util.NextSpanSeedURL))
	if err != nil {
		sp.Logger.Error("Error Fetching nextspanseed from ThemisServer ", "error", err)
		return nextSpanSeed, err
	}

	sp.Logger.Debug("Next span seed fetched")

	if err = jsoniter.ConfigFastest.Unmarshal(response.Result, &nextSpanSeed); err != nil {
		sp.Logger.Error("Error unmarshalling nextSpanSeed received from Themis Server", "error", err)
		return nextSpanSeed, err
	}

	return nextSpanSeed, nil
}

func (sp *SpanProcessor) checkReSpan(delayTime time.Duration) {
	if !util.RecoverSpanFinished {
		sp.Logger.Info("checkReSpan before recover span finished")
		return
	}

	// l2 chain block height
	chainBlockHeight, err := getCurrentChildBlock(sp.contractConnector)
	if err != nil {
		sp.Logger.Error("Error while getCurrentChildBlock", "error", err)
		return
	}
	sp.Logger.Debug("startCheckRespan get current child block", "blockHeight", chainBlockHeight)
	checkBlockHeight := chainBlockHeight + 1

	// get current sequencer
	currentSequencer, currentL2Epoch, err := getCurrentChildSequencer(sp.contractConnector, sp.cliCtx, checkBlockHeight)
	if err != nil {
		sp.Logger.Error("Error while getCurrentChildSequencer", "error", err)
		return
	}

	l2LatestEPoch, l2LatestEPochStart, _, err := getChildLatestEpoch(sp.contractConnector, sp.cliCtx)
	if err != nil {
		sp.Logger.Error("Error while getChildLatestEpoch", "error", err)
		return
	}
	if l2LatestEPoch < currentL2Epoch || currentSequencer == "" || currentSequencer == "0x0000000000000000000000000000000000000000" {
		currentSequencer, currentL2Epoch, _ = getCurrentChildSequencer(sp.contractConnector, sp.cliCtx, chainBlockHeight)
	}
	sp.Logger.Debug("startCheckRespan get current child sequencer", "currentSequencer", currentSequencer, "currentL2Epoch", currentL2Epoch)

	// query latest span
	latestSpan, err := sp.getLastSpan()
	if err != nil {
		sp.Logger.Error("Unable to fetch last span", "error", err)
		return
	}
	sp.Logger.Debug("startCheckRespan Found last span", "latestSpan", latestSpan.ID, "startBlock", latestSpan.StartBlock, "endBlock", latestSpan.EndBlock)

	if latestSpan.ID < uint64(currentL2Epoch) {
		sp.Logger.Info("startCheckRespan latest span too low wait recover span")
		return
	}

	// set newSpan
	newSpan := new(types.Span)
	if l2LatestEPochStart == int64(chainBlockHeight)+1 {
		newSpan, err = sp.getSpan(uint64(currentL2Epoch))
		if err != nil {
			sp.Logger.Error("Unable to fetch last span", "error", err)
			return
		}
		currentL2Epoch = currentL2Epoch - 1
	} else if latestSpan.ID == uint64(currentL2Epoch) {
		newSpan = latestSpan
		newSpan.ID = uint64(currentL2Epoch) + 1
	} else {
		newSpan, err = sp.getSpan(uint64(currentL2Epoch + 1))
		if err != nil {
			sp.Logger.Error("Unable to fetch last span", "error", err)
			return
		}
	}

	if newSpan == nil {
		sp.Logger.Error("Unable to fetch new span is nil")
		return
	}
	sp.Logger.Debug("startCheckRespan Found last span", "newSpan", newSpan.ID, "startBlock", newSpan.StartBlock, "endBlock", newSpan.EndBlock)

	// get last block from db
	hasLastBlock, _ := sp.storageClient.Has([]byte(lastMetisBlockKey), nil)
	if hasLastBlock {
		lastBlockBytes, err := sp.storageClient.Get([]byte(lastMetisBlockKey), nil)
		if err != nil {
			sp.Logger.Error("Error while fetching metis last block bytes from storage", "error", err)
			return
		}
		sp.Logger.Debug("startCheckRespan Got metis last block from bridge storage", "lastBlock", string(lastBlockBytes))

		dbBlockHeight, _ := strconv.ParseUint(string(lastBlockBytes), 10, 64)
		// block does not grow
		if chainBlockHeight == dbBlockHeight {
			sp.Logger.Debug("metis block does not grow", "dbBlockHeight", dbBlockHeight, "chainBlockHeight", chainBlockHeight)

			// Wait for the delay time to end or for a new block to be generated
			delayTicker := time.NewTicker(delayTime)
			chainBlockRefreshTicker := time.NewTicker(5 * time.Second)

			wg := sync.WaitGroup{}
			wg.Add(1)

			go func() {
			DELAY:
				for {
					select {
					case <-delayTicker.C:
						wg.Done()
						break DELAY
					case <-chainBlockRefreshTicker.C:
						chainBlockHeight, err := getCurrentChildBlock(sp.contractConnector)
						if err != nil {
							sp.Logger.Error("Error while getCurrentChildBlock", "error", err)
							continue
						}
						if chainBlockHeight > dbBlockHeight {
							wg.Done()
							break DELAY
						}
					}
				}
			}()
			wg.Wait()

			// poll metis number again
			chainBlockHeight, err = getCurrentChildBlock(sp.contractConnector)
			if err != nil {
				sp.Logger.Error("Error while getCurrentChildBlock", "error", err)
				return
			}

			// blocks are still not growing
			if chainBlockHeight == dbBlockHeight {
				sp.lock.Lock()
				defer sp.lock.Unlock()

				sp.Logger.Debug("metis block does not grow after delay", "dbBlockHeight", dbBlockHeight, "chainBlockHeight", chainBlockHeight)
				sp.Logger.Info("✅ Proposing new re-span", "dbBlockHeight", dbBlockHeight, "chainBlockHeight", chainBlockHeight)

				// generate new sequencer
				currentBatch, err := util.GetCurrentBatch(sp.cliCtx)
				if err != nil {
					sp.Logger.Error("Error GetCurrentBatch", "error", err)
					return
				}

				valSet, err := util.GetValidatorSet(sp.cliCtx)
				if err != nil {
					sp.Logger.Error("Error GetValidatorSet", "error", err)
					return
				}

				// query reSpan count
				reSpanCount := sp.getRespanCount(currentL2Epoch)
				seed := common.Hash{}
				seed.SetBytes(big.NewInt(int64(reSpanCount)).Bytes())

				newSequencer := helper.CalcMetisSequencerWithSeed(currentSequencer, chainBlockHeight, currentBatch, valSet, seed)
				sp.Logger.Info("generate new sequencer", "currentSequencer", currentSequencer, "newSequencer", newSequencer)
				if newSequencer == "" {
					sp.Logger.Error("invalid newSequencer", newSequencer)
					return
				}

				// notify metis respan start
				err = helper.NotifyRespanStart(currentSequencer, newSequencer, chainBlockHeight+1)
				if err != nil {
					sp.Logger.Error("Error while NotifyRespanStart", "error", err)
					return
				}

				// increase 100 block every respan
				reSpanIncrease := 100

				// get epoch length getSequencerEpochLength
				epochLength, err := getSequencerEpochLength(sp.contractConnector, sp.cliCtx)
				if err != nil {
					sp.Logger.Error("Error while getSequencerEpochLength", "error", err)
					return
				}
				endBlock := newSpan.EndBlock + uint64(reSpanIncrease)
				if endBlock-chainBlockHeight < uint64(epochLength) {
					endBlock += uint64(epochLength)
				}

				// broadcast reselect span msg to themis
				msg := metisTypes.MsgReProposeSpan{
					ID:              newSpan.ID,
					Proposer:        types.BytesToThemisAddress(helper.GetAddress()),
					CurrentProducer: types.HexToThemisAddress(currentSequencer),
					NextProducer:    types.HexToThemisAddress(newSequencer),
					CurrentL2Height: chainBlockHeight,
					CurrentL2Epoch:  uint64(currentL2Epoch),
					StartBlock:      chainBlockHeight + 1,
					EndBlock:        endBlock,
					ChainID:         newSpan.ChainID,
					Seed:            seed,
				}
				sp.Logger.Info("new span", "msg", msg)

				if sp.isValidator() {
					// return broadcast to themis
					if err := sp.txBroadcaster.BroadcastToThemis(msg, nil); err != nil {
						sp.Logger.Error("Error while broadcasting span to themis", "spanId", newSpan.ID, "startBlock", dbBlockHeight+1, "error", err)
						return
					}
				}

				// set respan count
				sp.setRespanCount(uint64(currentL2Epoch), reSpanCount+1)
			}

			sp.Logger.Info("metis block grows", "dbBlockHeight", dbBlockHeight, "chainBlockHeight", chainBlockHeight)
		}
	}

	// Set last block to storage
	if err = sp.storageClient.Put([]byte(lastMetisBlockKey), []byte(fmt.Sprintf("%v", chainBlockHeight)), nil); err != nil {
		sp.Logger.Error("rl.storageClient.Put", "Error", err)
	}
}

// Stop stops all necessary go routines
func (sp *SpanProcessor) Stop() {
	// cancel span polling
	sp.cancelSpanService()
}

func (sp *SpanProcessor) l2CommitSpan(eventBytes string) error {
	sp.Logger.Info("Received CommitEpoch request", "eventBytes", eventBytes)

	var event sdk.StringEvent
	if err := jsoniter.ConfigFastest.Unmarshal([]byte(eventBytes), &event); err != nil {
		sp.Logger.Error("Error unmarshalling event from themis", "error", err)
		return err
	}
	sp.Logger.Debug("processing CommitSpan confirmation event", "eventtype", event.Type)

	var spanID int64
	var isSpanRecover bool
	for _, attr := range event.Attributes {
		if attr.Key == metisTypes.AttributeKeySpanID {
			spanID, _ = strconv.ParseInt(attr.Value, 10, 64)
		}

		if attr.Key == metisTypes.AttributeKeySpanRecover {
			if attr.Value == "true" {
				isSpanRecover = true
			}
		}
	}

	if isSpanRecover {
		sp.Logger.Debug("processing CommitSpan is span cover,ignore it")
		return nil
	}

	if spanID == 0 {
		sp.Logger.Debug("processing CommitSpan invalid span id,ignore it")
		return nil
	}

	// l2 chain block height
	chainBlockHeight, err := getCurrentChildBlock(sp.contractConnector)
	if err != nil {
		sp.Logger.Error("Error while getCurrentChildBlock", "error", err)
		return err
	}
	sp.Logger.Info("startCheckRespan get current child block", "blockHeight", chainBlockHeight)

	// get current sequencer
	currentSequencer, _, err := getCurrentChildSequencer(sp.contractConnector, sp.cliCtx, chainBlockHeight+1)
	if err != nil {
		sp.Logger.Error("Error while getCurrentChildSequencer", "error", err)
		return err
	}

	span, err := sp.getSpan(uint64(spanID))
	if err != nil {
		sp.Logger.Error("Unable to fetch last span", "error", err)
		return err
	}

	if chainBlockHeight >= span.StartBlock {
		sp.Logger.Info("processing CommitSpan invalid start block")
		return nil
	}

	// check if current user is among span producers
	if strings.EqualFold(helper.GetAddressStr(), currentSequencer) {
		sp.Logger.Info("Received CommitEpoch request", "eventBytes", eventBytes)

		input, err := sp.caller.ValidatorSetABI.Pack("commitEpoch",
			big.NewInt(spanID),
			big.NewInt(int64(span.StartBlock)),
			big.NewInt(int64(span.EndBlock)),
			span.SelectedProducers[0].Signer.EthAddress())
		if err != nil {
			sp.Logger.Error("Pack tx input data", "error", err)
			return err
		}

		txId, err := sp.broadcastTx(input, types.CommitEpochToMetis, "", "", 0)
		if err != nil {
			sp.Logger.Error("l2CommitSpan broadcastTx", "error", err)
			return err
		}

		sp.Logger.Info("CommitEpoch to metis successfully", "txId", txId)
	}
	return nil
}

func (sp *SpanProcessor) l2RecommitSpan(eventBytes string) error {
	sp.Logger.Debug("Received ReCommitEpoch request", "eventBytes", eventBytes)

	var event sdk.StringEvent
	if err := jsoniter.ConfigFastest.Unmarshal([]byte(eventBytes), &event); err != nil {
		sp.Logger.Error("Error unmarshalling event from themis", "error", err)
		return err
	}

	var spanID int64
	var oldSpanID int64
	for _, attr := range event.Attributes {
		if attr.Key == metisTypes.AttributeKeySpanID {
			spanID, _ = strconv.ParseInt(attr.Value, 10, 64)
		}

		if attr.Key == metisTypes.AttributeKeyOldSpanID {
			oldSpanID, _ = strconv.ParseInt(attr.Value, 10, 64)
		}
	}

	if spanID == 0 {
		sp.Logger.Info("processing l2RecommitSpan invalid span info,ignore it")
		return nil
	}

	span, err := sp.getSpan(uint64(spanID))
	if err != nil {
		sp.Logger.Error("Unable to fetch last span", "error", err)
		return err
	}

	// l2 chain block height
	chainBlockHeight, err := getCurrentChildBlock(sp.contractConnector)
	if err != nil {
		sp.Logger.Error("Error while getCurrentChildBlock", "error", err)
		return err
	}
	if chainBlockHeight >= span.StartBlock {
		sp.Logger.Info("processing l2RecommitSpan invalid start block")
		return nil
	}

	newSequencer := span.SelectedProducers[0].Signer.EthAddress()

	oldSpan, err := sp.getSpan(uint64(oldSpanID))
	if err != nil {
		sp.Logger.Error("Unable to fetch old span", "error", err)
		return err
	}
	oldSpanNewSequencer := oldSpan.SelectedProducers[0].Signer.EthAddress()
	sp.Logger.Info("Received ReCommitEpoch request", "oldSpanNewSequencer", oldSpanNewSequencer, "newSequencer", newSequencer)

	// check if current user is newSequencer
	if strings.EqualFold(helper.GetAddressStr(), newSequencer.Hex()) {
		sp.Logger.Info("ReCommitEpoch new sequencer, start broadcastTx")

		input, err := sp.caller.ValidatorSetABI.Pack("recommitEpoch",
			big.NewInt(oldSpanID),
			big.NewInt(spanID),
			big.NewInt(int64(span.StartBlock)),
			big.NewInt(int64(span.EndBlock)),
			newSequencer)
		if err != nil {
			sp.Logger.Error("Pack tx input data", "error", err)
			return err
		}

		txId, err := sp.broadcastTx(input, types.ReCommitEpochToMetis, oldSpanNewSequencer.Hex(), newSequencer.Hex(), span.StartBlock)
		if err != nil {
			sp.Logger.Error("l2RecommitSpan broadcastTx", "error", err)
			return err
		}

		sp.Logger.Info("ReCommitEpoch to metis successfully", "txId", txId)
	}
	return nil
}

func (sp *SpanProcessor) broadcastTx(input []byte, signType types.SignType, oldSequencer, newSequencer string, startBlock uint64) (txId string, err error) {
	sp.Logger.Debug("broadcastTx", "input", hex.EncodeToString(input))
	// assembly transaction
	configParams, err := util.GetChainmanagerParams(sp.cliCtx)
	if err != nil {
		sp.Logger.Error("broadcastTx error while fetching chainmanager params", "error", err)
		return
	}
	sp.Logger.Debug("broadcastTx", "ValidatorSetAddress", configParams.ChainParams.ValidatorSetAddress)

	// get mpc address
	latestMpc, err := sp.getLastMpc()
	if err != nil {
		sp.Logger.Error("broadcastTx getLastMpc", "error", err)
		return
	}
	sp.Logger.Debug("broadcastTx", "latestMpc", latestMpc.MpcAddress.String())
	if latestMpc.ID == "" || latestMpc.MpcAddress.String() == "0x0000000000000000000000000000000000000000" {
		sp.Logger.Info("broadcastTx mpc not found, wait mpc generated")
		return "", errors.New("mpc not found")
	}

	nonce, err := sp.caller.GetMetisNonce(latestMpc.MpcAddress.EthAddress())
	if err != nil {
		sp.Logger.Error("broadcastTx GetMetisNonce", "error", err)
		return
	}

	gasPrice, err := sp.caller.GetMetisGasprice()
	if err != nil {
		sp.Logger.Error("broadcastTx GetMetisGasprice", "error", err)
		return
	}

	metisChainID, err := sp.caller.GetMetisChainID()
	if err != nil {
		sp.Logger.Error("broadcastTx GetMetisGasprice", "error", err)
		return
	}
	sp.Logger.Debug("broadcastTx tx params", "nonce", nonce, "gasPrice", gasPrice.String(), "chainID", metisChainID)

	gasLimit := helper.DefaultMetischainGasLimit
	tx := ethTypes.NewTransaction(nonce, configParams.ChainParams.ValidatorSetAddress.EthAddress(), big.NewInt(0), gasLimit, gasPrice, input)

	// calc sig hash
	txSigner := ethTypes.NewEIP155Signer(big.NewInt(int64(metisChainID)))
	sp.Logger.Debug("broadcastTx tx params", "signer chainID", txSigner.ChainID().Int64())

	signMsg := txSigner.Hash(tx).Bytes()

	// make propose mpc sign request
	signID := uuid.New().String()
	signData, _ := tx.MarshalBinary()
	msg := mpcTypes.MsgProposeMpcSign{
		ID:       signID,
		MpcID:    latestMpc.ID,
		SignType: signType,
		SignData: signData,
		SignMsg:  signMsg,
		Proposer: types.BytesToThemisAddress(helper.GetAddress()),
	}
	sp.Logger.Debug("broadcastTx", "signID", signID)

	// return broadcast to themis
	if err = sp.txBroadcaster.BroadcastToThemis(msg, nil); err != nil {
		sp.Logger.Error("broadcastTx error while broadcasting mpc to themis", "signID", signID, "error", err)
		return
	}

	// wait for sign finish
	ctx := context.Background()
	ctxTimeout, cancel := context.WithTimeout(ctx, 300*time.Second)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	var signature []byte
	signedTx := new(ethTypes.Transaction)
	var signedTxData []byte

SIGNATURE_EXIT:
	for {
		select {
		case <-ticker.C:
			mpcSignResult, _ := sp.getMpcSign(signID)
			if mpcSignResult != nil && mpcSignResult.Signature != nil {
				signature = mpcSignResult.Signature
				signedTxData = mpcSignResult.SignedTx
				err = signedTx.UnmarshalBinary(mpcSignResult.SignedTx)
				if err != nil {
					sp.Logger.Error("broadcastTx UnmarshalBinary", "err", err)
					return
				}

				break SIGNATURE_EXIT
			}
		case <-ctxTimeout.Done():
			sp.Logger.Debug("broadcastTx sign timeout")
			if signType == types.ReCommitEpochToMetis {
				sp.notifyRespanEnd()
			}
			return "", errors.New("wait sign timeout")

		}
	}
	sp.Logger.Info("broadcastTx sign success", "signature", hex.EncodeToString(signature), "len", len(signature))
	sp.Logger.Debug("broadcastTx signed tx", "data", hexutil.Encode(signedTxData))

	// notify metis respan start
	if signType == types.ReCommitEpochToMetis {
		err = helper.NotifyRespanStart(oldSequencer, newSequencer, startBlock)
		if err != nil {
			sp.Logger.Error("Error while NotifyRespanStart", "error", err)
			return
		}
	}

	txId = signedTx.Hash().Hex()
	err = helper.GetMetisClient().SendTransaction(context.Background(), signedTx)
	if err != nil {
		sp.Logger.Error("broadcastTx SendTransaction", "txId", txId, "error", err)
		return
	}

	return txId, err
}

func (sp *SpanProcessor) notifyRespanEnd() {
	// notify metis respan end
	zeroAddress := common.BigToAddress(big.NewInt(0))
	err := helper.NotifyRespanStart(zeroAddress.Hex(), zeroAddress.Hex(), 0)
	if err != nil {
		sp.Logger.Error("Error while notifyRespanEnd", "error", err)
		return
	}
}

func (sp *SpanProcessor) getLastMpc() (*types.Mpc, error) {
	// fetch latest start block from themis via rest query
	result, err := helper.FetchFromAPI(sp.cliCtx, helper.GetThemisServerEndpoint(fmt.Sprintf(util.LatestMpcURL, types.CommonMpcType)))
	if err != nil {
		sp.Logger.Error("Error while fetching latest mpc", "err", err)
		return nil, err
	}

	var lastMpc types.Mpc
	if err = jsoniter.ConfigFastest.Unmarshal(result.Result, &lastMpc); err != nil {
		sp.Logger.Error("Error unmarshalling mpc", "error", err)
		return nil, err
	}
	return &lastMpc, nil
}

func (sp *SpanProcessor) getMpcSign(signId string) (*types.MpcSign, error) {
	// fetch latest start block from themis via rest query
	result, err := helper.FetchFromAPI(sp.cliCtx, helper.GetThemisServerEndpoint(fmt.Sprintf(util.MpcSignByIdURL, signId)))
	if err != nil {
		sp.Logger.Error("Error while fetching latest mpc", "err", err)
		return nil, err
	}

	var sign types.MpcSign
	if err = jsoniter.ConfigFastest.Unmarshal(result.Result, &sign); err != nil {
		sp.Logger.Error("Error unmarshalling mpc", "error", err)
		return nil, err
	}
	return &sign, nil
}

func (sp *SpanProcessor) setRespanCount(currentL2Epoch, count uint64) {
	key := []byte(fmt.Sprintf(localRespanCountKey, currentL2Epoch))
	sp.storageClient.Put(key, []byte(fmt.Sprintf("%v", count)), nil)
}

func (sp *SpanProcessor) getRespanCount(currentL2Epoch int64) uint64 {
	key := []byte(fmt.Sprintf(localRespanCountKey, currentL2Epoch))
	lastBlockBytes, _ := sp.storageClient.Get(key, nil)
	reSpanCount, _ := strconv.ParseUint(string(lastBlockBytes), 10, 64)
	return reSpanCount
}

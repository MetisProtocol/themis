package util

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	mLog "github.com/RichardKnop/machinery/v1/log"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	httpClient "github.com/tendermint/tendermint/rpc/client"
	tmTypes "github.com/tendermint/tendermint/types"

	authTypes "github.com/metis-seq/themis/auth/types"
	"github.com/metis-seq/themis/bridge/setu/util/sqlite"
	chainManagerTypes "github.com/metis-seq/themis/chainmanager/types"
	"github.com/metis-seq/themis/helper"
	"github.com/metis-seq/themis/types"
	hmtypes "github.com/metis-seq/themis/types"
)

type BridgeEvent string

const (
	AccountDetailsURL      = "/auth/accounts/%v"
	ChainManagerParamsURL  = "/chainmanager/params"
	ProposersURL           = "/staking/proposer/%v"
	CurrentProposerURL     = "/staking/current-proposer"
	LatestSpanURL          = "/metis/latest-span"
	SpanByIdURL            = "/metis/span/%v"
	NextSpanInfoURL        = "/metis/prepare-next-span"
	NextSpanSeedURL        = "/metis/next-span-seed"
	ValidatorURL           = "/staking/validator/%v"
	CurrentValidatorSetURL = "staking/validator-set"
	CurrentL1BatchURL      = "staking/current-batch"
	StakingTxStatusURL     = "/staking/isoldtx"
	MpcSetURL              = "/mpc/set"
	LatestMpcURL           = "/mpc/latest/%v"
	MpcSignByIdURL         = "/mpc/sign/%v"

	TendermintUnconfirmedTxsURL      = "/unconfirmed_txs"
	TendermintUnconfirmedTxsCountURL = "/num_unconfirmed_txs"

	TransactionTimeout      = 1 * time.Minute
	CommitTimeout           = 2 * time.Minute
	TaskDelayBetweenEachVal = 10 * time.Second
	RetryTaskDelay          = 12 * time.Second
	RetryStateSyncTaskDelay = 24 * time.Second

	mempoolTxnCountDivisor = 1000

	// Bridge event types
	StakingEvent BridgeEvent = "staking"

	BridgeDBFlag       = "bridge-db"
	BridgeSqliteDBFlag = "bridge-sqlite-db"
)

var MpcCommonGenerated bool
var MpcStateCommitGenerated bool
var MpcRewardCommitGenerated bool
var MpcBlobCommitGenerated bool
var RecoverSpanFinished bool
var MetisTxChan chan *sqlite.MetisTx
var MetisTxLock sync.Mutex
var MetisTxCacheChan chan *sqlite.MetisTxCache

var logger log.Logger
var loggerOnce sync.Once

func init() {
	MetisTxChan = make(chan *sqlite.MetisTx, 8192)
	MetisTxCacheChan = make(chan *sqlite.MetisTxCache, 2048)
}

// Logger returns logger singleton instance
func Logger() log.Logger {
	loggerOnce.Do(func() {
		defaultLevel := "debug"
		logsWriter := helper.GetLogsWriter(helper.GetConfig().LogsWriterFile)
		logger = log.NewTMLogger(log.NewSyncWriter(logsWriter))

		logLevel := os.Getenv("LOG_LEVEL")
		option, err := log.AllowLevel(logLevel)
		if err != nil {
			// cosmos sdk is using different style of log format
			// and levels don't map well, config.toml
			// see: https://github.com/cosmos/cosmos-sdk/pull/8072
			logger.Error("Unable to parse logging level", "Error", err)
			option, err = log.AllowLevel(defaultLevel)
			if err != nil {
				logger.Error("failed to allow default log level", "Level", defaultLevel, "Error", err)
			}
		}

		logger = log.NewFilter(logger, option)

		// set no-op logger if log level is not debug for machinery
		if viper.GetString("log_level") != "debug" {
			mLog.SetDebug(NoopLogger{})
		}
	})

	return logger
}

// IsProposer  checks if we are proposer
func IsProposer(cliCtx cliContext.CLIContext) (bool, error) {
	var (
		proposers []hmtypes.Validator
		count     = uint64(1)
	)

	result, err := helper.FetchFromAPI(cliCtx,
		helper.GetThemisServerEndpoint(fmt.Sprintf(ProposersURL, strconv.FormatUint(count, 10))),
	)
	if err != nil {
		logger.Error("Error fetching proposers", "url", ProposersURL, "error", err)
		return false, err
	}

	err = jsoniter.ConfigFastest.Unmarshal(result.Result, &proposers)
	if err != nil {
		logger.Error("error unmarshalling proposer slice", "error", err)
		return false, err
	}

	if bytes.Equal(proposers[0].Signer.Bytes(), helper.GetAddress()) {
		return true, nil
	}

	return false, nil
}

// IsInProposerList checks if we are in current proposer
func IsInProposerList(cliCtx cliContext.CLIContext, count uint64) (bool, error) {
	logger.Debug("Skipping proposers", "count", strconv.FormatUint(count, 10))

	response, err := helper.FetchFromAPI(
		cliCtx,
		helper.GetThemisServerEndpoint(fmt.Sprintf(ProposersURL, strconv.FormatUint(count, 10))),
	)
	if err != nil {
		logger.Error("Unable to send request for next proposers", "url", ProposersURL, "error", err)
		return false, err
	}

	// unmarshall data from buffer
	var proposers []hmtypes.Validator
	if err := jsoniter.ConfigFastest.Unmarshal(response.Result, &proposers); err != nil {
		logger.Error("Error unmarshalling validator data ", "error", err)
		return false, err
	}

	logger.Debug("Fetched proposers list", "numberOfProposers", count)

	for _, proposer := range proposers {
		if bytes.Equal(proposer.Signer.Bytes(), helper.GetAddress()) {
			return true, nil
		}
	}

	return false, nil
}

// CalculateTaskDelay calculates delay required for current validator to propose the tx
// It solves for multiple validators sending same transaction.
func CalculateTaskDelay(cliCtx cliContext.CLIContext, event interface{}) (bool, time.Duration) {
	// defer LogElapsedTimeForStateSyncedEvent(event, "CalculateTaskDelay", time.Now())
	// calculate validator position
	valPosition := 0
	isCurrentValidator := false

	validatorSet, err := GetValidatorSet(cliCtx)
	if err != nil {
		logger.Error("Error getting current validatorset data ", "error", err)
		return isCurrentValidator, 0
	}

	logger.Info("Fetched current validatorset list", "currentValidatorcount", len(validatorSet.Validators))

	for i, validator := range validatorSet.Validators {
		if bytes.Equal(validator.Signer.Bytes(), helper.GetAddress()) {
			valPosition = i + 1
			isCurrentValidator = true

			break
		}
	}

	// Change calculation later as per the discussion
	// Currently it will multiply delay for every 1000 unconfirmed txns in mempool
	// For example if the current default delay is 12 Seconds
	// Then for upto 1000 txns it will stay as 12 only
	// For 1000-2000 It will be 24 seconds
	// For 2000-3000 it will be 36 seconds
	// Basically for every 1000 txns it will increase the factor by 1.

	mempoolFactor := GetUnconfirmedTxnCount(event) / mempoolTxnCountDivisor

	// calculate delay
	taskDelay := time.Duration(valPosition) * TaskDelayBetweenEachVal * time.Duration(mempoolFactor+1)

	return isCurrentValidator, taskDelay
}

// IsCurrentProposer checks if we are current proposer
func IsCurrentProposer(cliCtx cliContext.CLIContext) (bool, error) {
	var proposer hmtypes.Validator

	result, err := helper.FetchFromAPI(cliCtx, helper.GetThemisServerEndpoint(CurrentProposerURL))
	if err != nil {
		logger.Error("Error fetching proposers", "error", err)
		return false, err
	}

	if err = jsoniter.ConfigFastest.Unmarshal(result.Result, &proposer); err != nil {
		logger.Error("error unmarshalling validator", "error", err)
		return false, err
	}

	logger.Debug("Current proposer fetched", "validator", proposer.String())

	if bytes.Equal(proposer.Signer.Bytes(), helper.GetAddress()) {
		return true, nil
	}

	logger.Debug("We are not the current proposer")

	return false, nil
}

// IsEventSender check if we are the EventSender
func IsEventSender(cliCtx cliContext.CLIContext, validatorID uint64) bool {
	var validator hmtypes.Validator

	result, err := helper.FetchFromAPI(cliCtx,
		helper.GetThemisServerEndpoint(fmt.Sprintf(ValidatorURL, strconv.FormatUint(validatorID, 10))),
	)
	if err != nil {
		logger.Error("Error fetching proposers", "error", err)
		return false
	}

	if err = jsoniter.ConfigFastest.Unmarshal(result.Result, &validator); err != nil {
		logger.Error("error unmarshalling proposer slice", "error", err)
		return false
	}

	logger.Debug("Current event sender received", "validator", validator.String())

	return bytes.Equal(validator.Signer.Bytes(), helper.GetAddress())
}

// CreateURLWithQuery receives the uri and parameters in key value form
// it will return the new url with the given query from the parameter
func CreateURLWithQuery(uri string, param map[string]interface{}) (string, error) {
	urlObj, err := url.Parse(uri)
	if err != nil {
		return uri, err
	}

	query := urlObj.Query()
	for k, v := range param {
		query.Set(k, fmt.Sprintf("%v", v))
	}

	urlObj.RawQuery = query.Encode()

	return urlObj.String(), nil
}

// WaitForOneEvent subscribes to a websocket event for the given
// event time and returns upon receiving it one time, or
// when the timeout duration has expired.
//
// This handles subscribing and unsubscribing under the hood
func WaitForOneEvent(tx tmTypes.Tx, client *httpClient.HTTP) (tmTypes.TMEventData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CommitTimeout)
	defer cancel()

	// subscriber
	subscriber := hex.EncodeToString(tx.Hash())

	// query
	query := tmTypes.EventQueryTxFor(tx).String()

	// register for the next event of this type
	eventCh, err := client.Subscribe(ctx, subscriber, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to subscribe")
	}

	// make sure to unregister after the test is over
	defer func() {
		if err := client.UnsubscribeAll(ctx, subscriber); err != nil {
			logger.Error("WaitForOneEvent | UnsubscribeAll", "Error", err)
		}
	}()

	select {
	case event := <-eventCh:
		return event.Data, nil
	case <-ctx.Done():
		return nil, errors.New("timed out waiting for event")
	}
}

// IsCatchingUp checks if the themis node you are connected to is fully synced or not
// returns true when synced
func IsCatchingUp(cliCtx cliContext.CLIContext) bool {
	resp, err := helper.GetNodeStatus(cliCtx)
	if err != nil {
		logger.Error("Can not get node status", "err", err)
		return true
	}

	logger.Info("get node status ", resp.SyncInfo.CatchingUp)
	return resp.SyncInfo.CatchingUp
}

// GetAccount returns themis auth account
func GetAccount(cliCtx cliContext.CLIContext, address types.ThemisAddress) (account authTypes.Account, err error) {
	url := helper.GetThemisServerEndpoint(fmt.Sprintf(AccountDetailsURL, address))

	// call account rest api
	response, err := helper.FetchFromAPI(cliCtx, url)
	if err != nil {
		return
	}

	if err = cliCtx.Codec.UnmarshalJSON(response.Result, &account); err != nil {
		logger.Error("Error unmarshalling account details", "url", url)
		return
	}

	return
}

// GetChainmanagerParams return chain manager params
func GetChainmanagerParams(cliCtx cliContext.CLIContext) (*chainManagerTypes.Params, error) {
	response, err := helper.FetchFromAPI(
		cliCtx,
		helper.GetThemisServerEndpoint(ChainManagerParamsURL),
	)
	if err != nil {
		logger.Error("Error fetching chainmanager params", "err", err)
		return nil, err
	}

	var params chainManagerTypes.Params
	if err = jsoniter.ConfigFastest.Unmarshal(response.Result, &params); err != nil {
		logger.Error("Error unmarshalling chainmanager params", "url", ChainManagerParamsURL, "err", err)
		return nil, err
	}

	return &params, nil
}

// AppendPrefix returns publickey in uncompressed format
func AppendPrefix(signerPubKey []byte) []byte {
	// append prefix - "0x04" as themis uses publickey in uncompressed format. Refer below link
	// https://superuser.com/questions/1465455/what-is-the-size-of-public-key-for-ecdsa-spec256r1
	prefix := make([]byte, 1)
	prefix[0] = byte(0x04)
	signerPubKey = append(prefix[:], signerPubKey[:]...)

	return signerPubKey
}

// GetValidatorNonce fetches validator nonce and height
func GetValidatorNonce(cliCtx cliContext.CLIContext, validatorID uint64) (uint64, int64, error) {
	var validator hmtypes.Validator

	result, err := helper.FetchFromAPI(cliCtx,
		helper.GetThemisServerEndpoint(fmt.Sprintf(ValidatorURL, strconv.FormatUint(validatorID, 10))),
	)

	if err != nil {
		logger.Error("Error fetching validator data", "error", err)
		return 0, 0, err
	}

	if err = jsoniter.ConfigFastest.Unmarshal(result.Result, &validator); err != nil {
		logger.Error("error unmarshalling validator data", "error", err)
		return 0, 0, err
	}

	logger.Debug("Validator data received ", "validator", validator.String())

	return validator.Nonce, result.Height, nil
}

// GetValidatorSet fetches the current validator set
func GetValidatorSet(cliCtx cliContext.CLIContext) (*hmtypes.ValidatorSet, error) {
	response, err := helper.FetchFromAPI(cliCtx, helper.GetThemisServerEndpoint(CurrentValidatorSetURL))
	if err != nil {
		logger.Error("Unable to send request for current validatorset", "url", CurrentValidatorSetURL, "error", err)
		return nil, err
	}

	var validatorSet hmtypes.ValidatorSet
	if err = jsoniter.ConfigFastest.Unmarshal(response.Result, &validatorSet); err != nil {
		logger.Error("Error unmarshalling current validatorset data ", "error", err)
		return nil, err
	}

	return &validatorSet, nil
}

type CurrentBatch struct {
	Batch uint64 `json:"batch"`
}

// GetCurrentBatch fetches the current validator set
func GetCurrentBatch(cliCtx cliContext.CLIContext) (uint64, error) {
	response, err := helper.FetchFromAPI(cliCtx, helper.GetThemisServerEndpoint(CurrentL1BatchURL))
	if err != nil {
		logger.Error("Unable to send request for current l1 batch", "url", CurrentL1BatchURL, "error", err)
		return 0, err
	}

	var currentL1Batch CurrentBatch
	if err = jsoniter.ConfigFastest.Unmarshal(response.Result, &currentL1Batch); err != nil {
		logger.Error("Error unmarshalling current l1 batch data ", "error", err)
		return 0, err
	}

	return currentL1Batch.Batch, nil
}

func GetUnconfirmedTxnCount(event interface{}) int {
	// defer LogElapsedTimeForStateSyncedEvent(event, "GetUnconfirmedTxnCount", time.Now())

	endpoint := helper.GetConfig().TendermintRPCUrl + TendermintUnconfirmedTxsCountURL

	resp, err := helper.Client.Get(endpoint)
	if err != nil || resp.StatusCode != http.StatusOK {
		logger.Error("Error fetching mempool txs count", "url", endpoint, "error", err)
		return 0
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		logger.Error("Error fetching mempool txs count", "error", err)
		return 0
	}

	// a minimal response of the unconfirmed txs
	var response TendermintUnconfirmedTxs

	err = jsoniter.ConfigFastest.Unmarshal(body, &response)
	if err != nil {
		logger.Error("Error unmarshalling response received from Themis Server", "error", err)
		return 0
	}

	count, _ := strconv.Atoi(response.Result.Total)

	return count
}

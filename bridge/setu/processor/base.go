package processor

import (
	"github.com/cosmos/cosmos-sdk/client"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/log"
	httpClient "github.com/tendermint/tendermint/rpc/client"

	"github.com/metis-seq/themis/bridge/setu/broadcaster"
	"github.com/metis-seq/themis/bridge/setu/util"
	"github.com/metis-seq/themis/bridge/setu/util/sqlite"
	"github.com/metis-seq/themis/helper"
)

// Processor defines a block header listerner for Rootchain, Metischain, Themis
type Processor interface {
	Start() error

	RegisterTasks()

	String() string

	Stop()
}

type BaseProcessor struct {
	Logger log.Logger
	name   string
	quit   chan struct{}

	// tx broadcaster
	txBroadcaster *broadcaster.TxBroadcaster

	// The "subclass" of BaseProcessor
	impl Processor

	// cli context
	cliCtx cliContext.CLIContext

	// contract caller
	contractConnector helper.ContractCaller

	// http client to subscribe to
	httpClient *httpClient.HTTP

	// storage client
	storageClient *leveldb.DB

	// sql client
	sqlClient *sqlite.SqliteDB
}

// NewBaseProcessor creates a new BaseProcessor.
func NewBaseProcessor(cdc *codec.Codec, httpClient *httpClient.HTTP, txBroadcaster *broadcaster.TxBroadcaster, name string, impl Processor) *BaseProcessor {
	logger := util.Logger().With("service", "processor", "module", name)

	cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	// cliCtx.BroadcastMode = client.BroadcastAsync
	cliCtx.BroadcastMode = client.BroadcastBlock
	cliCtx.TrustNode = true

	contractCaller, err := helper.NewContractCaller()
	if err != nil {
		logger.Error("Error while getting root chain instance", "error", err)
		panic(err)
	}

	if logger == nil {
		logger = log.NewNopLogger()
	}

	// creating syncer object
	return &BaseProcessor{
		Logger: logger,
		name:   name,
		quit:   make(chan struct{}),
		impl:   impl,

		cliCtx:            cliCtx,
		contractConnector: contractCaller,
		txBroadcaster:     txBroadcaster,
		httpClient:        httpClient,
		storageClient:     util.GetBridgeDBInstance(viper.GetString(util.BridgeDBFlag)),
		sqlClient:         sqlite.GetBridgeSqlDBInstance(viper.GetString(util.BridgeSqliteDBFlag)),
	}
}

// String implements Service by returning a string representation of the service.
func (bp *BaseProcessor) String() string {
	return bp.name
}

// OnStop stops all necessary go routines
func (bp *BaseProcessor) Stop() {
	// override to stop any go-routines in individual processors
}

// isOldTx checks if the transaction already exists in the chain or not
// It is a generic function, which is consumed in all processors
func (bp *BaseProcessor) isOldTx(cliCtx cliContext.CLIContext, txHash string, logIndex uint64, eventType util.BridgeEvent, event interface{}) (bool, error) {
	// defer util.LogElapsedTimeForStateSyncedEvent(event, "isOldTx", time.Now())

	queryParam := map[string]interface{}{
		"txhash":   txHash,
		"logindex": logIndex,
	}

	// define the endpoint based on the type of event
	var endpoint string

	switch eventType {
	case util.StakingEvent:
		endpoint = helper.GetThemisServerEndpoint(util.StakingTxStatusURL)
	}

	url, err := util.CreateURLWithQuery(endpoint, queryParam)
	if err != nil {
		bp.Logger.Error("Error in creating url", "endpoint", endpoint, "error", err)
		return false, err
	}

	res, err := helper.FetchFromAPI(bp.cliCtx, url)
	if err != nil {
		// bp.Logger.Error("Error fetching tx status", "url", url, "error", err)
		return false, err
	}

	var status bool
	if err := jsoniter.ConfigFastest.Unmarshal(res.Result, &status); err != nil {
		bp.Logger.Error("Error unmarshalling tx status received from Themis Server", "error", err)
		return false, err
	}

	return status, nil
}

package processor

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/common"
	httpClient "github.com/tendermint/tendermint/rpc/client"

	"github.com/metis-seq/themis/bridge/setu/broadcaster"
	"github.com/metis-seq/themis/bridge/setu/util"
	"github.com/metis-seq/themis/helper"
)

const (
	processorServiceStr = "processor-service"
)

// ProcessorService starts and stops all event processors
type ProcessorService struct {
	// Base service
	common.BaseService

	processors []Processor
}

// NewProcessorService returns new service object for processing queue msg
func NewProcessorService(
	cdc *codec.Codec,
	httpClient *httpClient.HTTP,
	txBroadcaster *broadcaster.TxBroadcaster,
) *ProcessorService {
	var logger = util.Logger().With("module", processorServiceStr)
	// creating processor object
	processorService := &ProcessorService{}

	contractCaller, err := helper.NewContractCaller()
	if err != nil {
		panic(err)
	}

	processorService.BaseService = *common.NewBaseService(logger, processorServiceStr, processorService)

	//
	// Intitialize processors
	//

	// initialize staking processor
	stakingProcessor := NewStakingProcessor(&contractCaller.StakingInfoABI)
	stakingProcessor.BaseProcessor = *NewBaseProcessor(cdc, httpClient, txBroadcaster, "staking", stakingProcessor)

	// initialize span processor
	spanProcessor := &SpanProcessor{
		caller: contractCaller,
	}
	spanProcessor.BaseProcessor = *NewBaseProcessor(cdc, httpClient, txBroadcaster, "span", spanProcessor)

	// init mpc processor
	mpcProcessor := &MpcProcessor{}
	mpcProcessor.BaseProcessor = *NewBaseProcessor(cdc, httpClient, txBroadcaster, "mpc", mpcProcessor)

	// init rpcTx processor
	rpcTxProcessor := NewRpcTxProcessor()
	rpcTxProcessor.BaseProcessor = *NewBaseProcessor(cdc, httpClient, txBroadcaster, "rpcTx", rpcTxProcessor)

	//
	// Select processors
	//

	// add into processor list
	startAll := viper.GetBool("all")
	onlyServices := viper.GetStringSlice("only")

	if startAll {
		processorService.processors = append(processorService.processors,
			stakingProcessor,
			spanProcessor,
			mpcProcessor,
			rpcTxProcessor,
		)
	} else {
		for _, service := range onlyServices {
			switch service {
			case "staking":
				processorService.processors = append(processorService.processors, stakingProcessor)
			case "span":
				processorService.processors = append(processorService.processors, spanProcessor)
			case "mpc":
				processorService.processors = append(processorService.processors, mpcProcessor)
			case "rpcTx":
				processorService.processors = append(processorService.processors, rpcTxProcessor)
			}
		}
	}

	if len(processorService.processors) == 0 {
		panic("No processors selected. Use --all or --only <coma-seprated processors>")
	}

	return processorService
}

// OnStart starts new block subscription
func (processorService *ProcessorService) OnStart() error {
	if err := processorService.BaseService.OnStart(); err != nil {
		processorService.Logger.Error("OnStart | OnStart", "Error", err)
	} // Always call the overridden method.

	// start processors
	for _, processor := range processorService.processors {
		processor.RegisterTasks()

		go func(processor Processor) {
			if err := processor.Start(); err != nil {
				processorService.Logger.Error("OnStart | processor.Start", "Error", err)
			}
		}(processor)
	}

	processorService.Logger.Info("all processors Started")

	return nil
}

// OnStop stops all necessary go routines
func (processorService *ProcessorService) OnStop() {
	processorService.BaseService.OnStop() // Always call the overridden method.
	// start chain listeners
	for _, processor := range processorService.processors {
		processor.Stop()
	}

	processorService.Logger.Info("all processors stopped")
}

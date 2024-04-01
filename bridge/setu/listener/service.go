package listener

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/metis-seq/themis/bridge/setu/util"
	"github.com/metis-seq/themis/helper"
	"github.com/tendermint/tendermint/libs/common"
	httpClient "github.com/tendermint/tendermint/rpc/client"
)

const (
	ListenerServiceStr = "listener"

	RootChainListenerStr  = "rootchain"
	ThemisListenerStr     = "themis"
	MetisChainListenerStr = "metischain"
)

// var logger = util.Logger().With("service", ListenerServiceStr)

// ListenerService starts and stops all chain event listeners
type ListenerService struct {
	// Base service
	common.BaseService
	listeners []Listener
}

// NewListenerService returns new service object for listneing to events
func NewListenerService(cdc *codec.Codec, httpClient *httpClient.HTTP) *ListenerService {
	var logger = util.Logger().With("service", ListenerServiceStr)

	// creating listener object
	listenerService := &ListenerService{}

	listenerService.BaseService = *common.NewBaseService(logger, ListenerServiceStr, listenerService)

	rootchainListener := NewRootChainListener()
	rootchainListener.BaseListener = *NewBaseListener(cdc, httpClient, helper.GetMainClient(), RootChainListenerStr, rootchainListener)
	listenerService.listeners = append(listenerService.listeners, rootchainListener)

	metisListener := NewMetisListener()
	metisListener.BaseListener = *NewBaseListener(cdc, httpClient, helper.GetMetisClient(), MetisChainListenerStr, metisListener)
	listenerService.listeners = append(listenerService.listeners, metisListener)

	themisListener := &ThemisListener{}
	themisListener.BaseListener = *NewBaseListener(cdc, httpClient, nil, ThemisListenerStr, themisListener)
	listenerService.listeners = append(listenerService.listeners, themisListener)

	return listenerService
}

// OnStart starts new block subscription
func (listenerService *ListenerService) OnStart() error {
	if err := listenerService.BaseService.OnStart(); err != nil {
		listenerService.Logger.Error("OnStart | OnStart", "Error", err)
	} // Always call the overridden method.

	// start chain listeners
	for _, listener := range listenerService.listeners {
		if err := listener.Start(); err != nil {
			listenerService.Logger.Error("OnStart | Start", "Error", err)
		}
	}

	listenerService.Logger.Info("all listeners Started")

	return nil
}

// OnStop stops all necessary go routines
func (listenerService *ListenerService) OnStop() {
	listenerService.BaseService.OnStop() // Always call the overridden method.

	// start chain listeners
	for _, listener := range listenerService.listeners {
		listener.Stop()
	}

	listenerService.Logger.Info("all listeners stopped")
}

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/common"
	httpClient "github.com/tendermint/tendermint/rpc/client"
	"golang.org/x/sync/errgroup"

	"github.com/metis-seq/themis/app"
	"github.com/metis-seq/themis/bridge/setu/broadcaster"
	"github.com/metis-seq/themis/bridge/setu/listener"
	"github.com/metis-seq/themis/bridge/setu/processor"
	"github.com/metis-seq/themis/bridge/setu/rpc"
	"github.com/metis-seq/themis/bridge/setu/util"
	"github.com/metis-seq/themis/bridge/setu/util/sqlite"

	"github.com/spf13/viper"

	"github.com/metis-seq/themis/helper"
)

const (
	waitDuration = 10 * time.Second
)

// StartBridgeWithCtx starts bridge service and is able to shutdow gracefully
// returns service errors, if any
func StartBridgeWithCtx(cmd *cobra.Command, shutdownCtx context.Context) error {
	logger.Info("start bridge helper.Config", helper.GetConfig())

	// create codec
	cdc := app.MakeCodec()

	_txBroadcaster := broadcaster.NewTxBroadcaster(cdc)
	_httpClient := httpClient.NewHTTP(helper.GetConfig().TendermintRPCUrl, "/websocket")

	rpcListenAddr := cmd.Flag(rpcServerFlag).Value.String()

	// selected services to start
	services := []common.Service{}
	services = append(services,
		listener.NewListenerService(cdc, _httpClient),
		processor.NewProcessorService(cdc, _httpClient, _txBroadcaster),
		rpc.NewMetisEthService(rpcListenAddr, cdc),
	)

	// Start http client
	err := _httpClient.Start()
	if err != nil {
		logger.Error("Error connecting to server: %v", err)
		return err
	}

	// cli context
	cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	cliCtx.BroadcastMode = client.BroadcastAsync
	cliCtx.TrustNode = true

	// start bridge services only when node fully synced
	loop := true
	for loop {
		select {
		case <-shutdownCtx.Done():
			return nil
		case <-time.After(waitDuration):
			if !util.IsCatchingUp(cliCtx) {
				logger.Info("Node up to date, starting bridge services")

				loop = false
			} else {
				logger.Info("Waiting for themis to be synced")
			}
		}
	}

	// start services
	var g errgroup.Group

	for _, service := range services {
		// loop variable must be captured
		srv := service

		g.Go(func() error {
			if err := srv.Start(); err != nil {
				logger.Error("GetStartCmd | serv.Start", "Error", err)
				return err
			}
			<-srv.Quit()
			return nil
		})
	}

	// shutdown phase
	g.Go(func() error {
		// wait for interrupt and start shut down
		<-shutdownCtx.Done()

		logger.Info("Received stop signal - Stopping all themis bridge services")
		for _, service := range services {
			srv := service
			if srv.IsRunning() {
				if err := srv.Stop(); err != nil {
					logger.Error("GetStartCmd | service.Stop", "Error", err)
					return err
				}
			}
		}
		// stop http client
		if err := _httpClient.Stop(); err != nil {
			logger.Error("GetStartCmd | _httpClient.Stop", "Error", err)
			return err
		}
		// stop db instance
		util.CloseBridgeDBInstance()
		sqlite.CloseBridgeSqlDBInstance()

		return nil
	})

	// wait for all routines to finish and log error
	if err := g.Wait(); err != nil {
		logger.Error("Bridge stopped", "err", err)
		return err
	}

	return nil
}

// StartBridge starts bridge service, isStandAlone prevents os.Exit if the bridge started as side service
func StartBridge(rpcListenAddr string, isStandAlone bool) {
	// create codec
	cdc := app.MakeCodec()

	_txBroadcaster := broadcaster.NewTxBroadcaster(cdc)
	_httpClient := httpClient.NewHTTP(helper.GetConfig().TendermintRPCUrl, "/websocket")

	// selected services to start
	services := []common.Service{}
	services = append(services,
		listener.NewListenerService(cdc, _httpClient),
		processor.NewProcessorService(cdc, _httpClient, _txBroadcaster),
		rpc.NewMetisEthService(rpcListenAddr, cdc),
	)

	// sync group
	var wg sync.WaitGroup

	// go routine to catch signal
	catchSignal := make(chan os.Signal, 1)
	signal.Notify(catchSignal, os.Interrupt, syscall.SIGTERM)

	go func() {
		// sig is a ^C, handle it
		for range catchSignal {
			// stop processes
			logger.Info("Received stop signal - Stopping all services")

			for _, service := range services {
				if err := service.Stop(); err != nil {
					logger.Error("GetStartCmd | service.Stop", "Error", err)
				}
			}

			// stop http client
			if err := _httpClient.Stop(); err != nil {
				logger.Error("GetStartCmd | _httpClient.Stop", "Error", err)
			}

			// stop db instance
			util.CloseBridgeDBInstance()
			sqlite.CloseBridgeSqlDBInstance()

			// exit
			if isStandAlone {
				os.Exit(1)
			}
		}
	}()

	// Start http client
	err := _httpClient.Start()
	if err != nil {
		panic(fmt.Sprintf("Error connecting to server %v", err))
	}

	// cli context
	cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	cliCtx.BroadcastMode = client.BroadcastAsync
	cliCtx.TrustNode = true

	// start bridge services only when node fully synced
	for {
		if !util.IsCatchingUp(cliCtx) {
			logger.Info("Node upto date, starting bridge services")
			break
		} else {
			logger.Info("Waiting for themis to be synced")
		}

		time.Sleep(waitDuration)
	}

	// start all processes
	for _, service := range services {
		go func(serv common.Service) {
			defer wg.Done()
			// TODO handle error while starting service
			if err := serv.Start(); err != nil {
				logger.Error("GetStartCmd | serv.Start", "Error", err)
			}

			<-serv.Quit()
		}(service)
	}

	// wait for all processes
	wg.Add(len(services))
	wg.Wait()
}

// GetStartCmd returns the start command to start bridge
func GetStartCmd() *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start bridge server",
		Run: func(cmd *cobra.Command, args []string) {
			rpcListenAddr := cmd.Flag(rpcServerFlag).Value.String()
			StartBridge(rpcListenAddr, true)
		}}

	// log level
	startCmd.Flags().String(helper.LogLevel, "info", "Log level for bridge")

	if err := viper.BindPFlag(helper.LogLevel, startCmd.Flags().Lookup(helper.LogLevel)); err != nil {
		logger.Error("GetStartCmd | BindPFlag | logLevel", "Error", err)
	}

	startCmd.Flags().Bool("all", false, "start all bridge services")

	if err := viper.BindPFlag("all", startCmd.Flags().Lookup("all")); err != nil {
		logger.Error("GetStartCmd | BindPFlag | all", "Error", err)
	}

	startCmd.Flags().StringSlice("only", []string{}, "comma separated bridge services to start")

	if err := viper.BindPFlag("only", startCmd.Flags().Lookup("only")); err != nil {
		logger.Error("GetStartCmd | BindPFlag | only", "Error", err)
	}

	return startCmd
}

func init() {
	rootCmd.AddCommand(GetStartCmd())
}

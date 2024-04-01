package cmd

import (
	"strconv"

	"github.com/metis-seq/themis/bridge/setu/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	lastRootBlockKey = "rootchain-last-block" // storage key
)

// updateRootHeightCmd represents the start command
var updateRootHeightCmd = &cobra.Command{
	Use:   "update-root-height",
	Short: "Update root chain handled height",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := util.Logger().With("cmd", "update_root_height")

		// command line height
		height := args[0]

		// storage init
		storageClient := util.GetBridgeDBInstance(viper.GetString(util.BridgeDBFlag))

		// query current height
		beforeUpdateHeight, _ := queryCurrentRootHeight(storageClient)
		logger.Info("current storage height before update", "height", beforeUpdateHeight)

		// Set last block to storage
		if err := storageClient.Put([]byte(lastRootBlockKey), []byte(height), nil); err != nil {
			logger.Error("rl.storageClient.Put", "Error", err)
			return err
		}

		// query current height
		afterUpdateHeight, _ := queryCurrentRootHeight(storageClient)
		logger.Info("current storage height after update", "height", afterUpdateHeight)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateRootHeightCmd)
}

func queryCurrentRootHeight(storageClient *leveldb.DB) (uint64, error) {
	// get current height
	lastBlockBytes, err := storageClient.Get([]byte(lastRootBlockKey), nil)
	if err != nil {
		logger.Info("Error while fetching last block bytes from storage", "error", err)
		return 0, err
	}
	logger.Debug("Got rootchain last block from bridge storage", "lastBlock", string(lastBlockBytes))

	result, err := strconv.ParseUint(string(lastBlockBytes), 10, 64)
	if err != nil {
		logger.Info("Error while parse lastBlockBytes", "error", err)
		return 0, err
	}
	return result, nil
}

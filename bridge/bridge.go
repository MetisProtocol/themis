package main

import (
	"fmt"
	"os"

	"github.com/metis-seq/themis/bridge/cmd"
	"github.com/metis-seq/themis/helper"
	"github.com/spf13/viper"
)

func main() {
	var logger = helper.Logger.With("module", "bridge/cmd/")
	rootCmd := cmd.BridgeCommands(viper.GetViper(), logger, "bridge-main")

	// add themis flags
	helper.DecorateWithThemisFlags(rootCmd, viper.GetViper(), logger, "bridge-main")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

package helper

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"

	cfg "github.com/tendermint/tendermint/config"
)

// Test - to check themis config
func TestThemisConfig(t *testing.T) {
	t.Parallel()

	// cli context
	tendermintNode := "tcp://localhost:26657"
	viper.Set(TendermintNodeFlag, tendermintNode)
	viper.Set("log_level", "info")
	// cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	// cliCtx.BroadcastMode = client.BroadcastSync
	// cliCtx.TrustNode = true

	InitThemisConfig(os.ExpandEnv("$HOME/.themisd"))

	fmt.Println("Address", GetAddress())

	pubKey := GetPubKey()

	fmt.Println("PublicKey", pubKey.String())
}

func TestThemisConfigNewSelectionAlgoHeight(t *testing.T) {
	t.Parallel()

	data := map[string]bool{"testnet": false, "mainnet": false, "local": true}
	for chain, shouldBeZero := range data {
		conf.MetisRPCUrl = "" // allow config to be loaded again

		viper.Set("chain", chain)

		InitThemisConfig(os.ExpandEnv("$HOME/.themisd"))

		nsah := GetNewSelectionAlgoHeight()
		if nsah == 0 && !shouldBeZero || nsah != 0 && shouldBeZero {
			t.Errorf("Invalid GetNewSelectionAlgoHeight = %d for chain %s", nsah, chain)
		}
	}
}

func TestThemisConfigUpdateTendermintConfig(t *testing.T) {
	t.Parallel()

	type teststruct struct {
		chain string
		viper string
		def   string
		value string
	}

	data := []teststruct{
		{chain: "testnet", viper: "viper", def: "default", value: "viper"},
		{chain: "testnet", viper: "viper", def: "", value: "viper"},
		{chain: "testnet", viper: "", def: "default", value: "default"},
		{chain: "testnet", viper: "", def: "", value: DefaultTestnetSeeds},
		{chain: "mainnet", viper: "viper", def: "default", value: "viper"},
		{chain: "mainnet", viper: "viper", def: "", value: "viper"},
		{chain: "mainnet", viper: "", def: "default", value: "default"},
		{chain: "mainnet", viper: "", def: "", value: DefaultMainnetSeeds},
		{chain: "local", viper: "viper", def: "default", value: "viper"},
		{chain: "local", viper: "viper", def: "", value: "viper"},
		{chain: "local", viper: "", def: "default", value: "default"},
		{chain: "local", viper: "", def: "", value: ""},
	}

	oldConf := conf.Chain
	viperObj := viper.New()
	tendermintConfig := cfg.DefaultConfig()

	for _, ts := range data {
		conf.Chain = ts.chain
		tendermintConfig.P2P.Seeds = ts.def
		viperObj.Set(SeedsFlag, ts.viper)
		UpdateTendermintConfig(tendermintConfig, viperObj)

		if tendermintConfig.P2P.Seeds != ts.value {
			t.Errorf("Invalid UpdateTendermintConfig, tendermintConfig.P2P.Seeds not set correctly")
		}
	}

	conf.Chain = oldConf
}

func TestThemisConfigWithEnv(t *testing.T) {
	t.Parallel()

	// cli context
	tendermintNode := "tcp://localhost:26657"
	viper.Set(TendermintNodeFlag, tendermintNode)
	viper.Set("log_level", "info")

	os.Setenv("ETH_RPC_URL", "http://127.0.0.1:38545")
	os.Setenv("METIS_RPC_URL", "http://l2geth:38545")

	InitThemisConfig(os.ExpandEnv("$HOME/.themisd"))

	pubKey := GetPubKey()
	fmt.Println("PublicKey", pubKey.String())

	fmt.Println("ethRpcUrl", conf.EthRPCUrl)
	fmt.Println("metisRpcUrl", conf.MetisRPCUrl)
}

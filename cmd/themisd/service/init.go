package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/metis-seq/themis/app"
	authTypes "github.com/metis-seq/themis/auth/types"
	"github.com/metis-seq/themis/helper"
	metisTypes "github.com/metis-seq/themis/metis/types"
	mpcTypes "github.com/metis-seq/themis/mpc/types"
	stakingcli "github.com/metis-seq/themis/staking/client/cli"
	stakingTypes "github.com/metis-seq/themis/staking/types"

	hmTypes "github.com/metis-seq/themis/types"
)

type initThemisConfig struct {
	clientHome  string
	chainID     string
	validatorID int64
	chain       string
	forceInit   bool
}

func themisInit(_ *server.Context, cdc *codec.Codec, initConfig *initThemisConfig, config *cfg.Config) error {
	conf := helper.GetDefaultThemisConfig()
	conf.Chain = initConfig.chain
	WriteDefaultThemisConfig(filepath.Join(config.RootDir, "config/themis-config.toml"), conf)

	nodeID, valPubKey, valPriKey, err := InitializeNodeValidatorFiles(config)
	if err != nil {
		return err
	}

	// do not execute init if forceInit is false and genesis.json already exists (or we do not have permission to write to file)
	writeGenesis := initConfig.forceInit

	if !writeGenesis {
		// When not forcing, check if genesis file exists
		_, err := os.Stat(config.GenesisFile())
		if err != nil && errors.Is(err, os.ErrNotExist) {
			logger.Info(fmt.Sprintf("Genesis file %v not found, writing genesis file\n", config.GenesisFile()))

			writeGenesis = true
		} else if err == nil {
			logger.Info(fmt.Sprintf("Found genesis file %v, skipping writing genesis file\n", config.GenesisFile()))
		} else {
			logger.Error(fmt.Sprintf("Error checking if genesis file %v exists: %v\n", config.GenesisFile(), err))
			return err
		}
	} else {
		logger.Info(fmt.Sprintf("Force writing genesis file to %v\n", config.GenesisFile()))
	}

	if writeGenesis {
		genesisCreated, err := helper.WriteGenesisFile(initConfig.chain, config.GenesisFile(), cdc)
		if err != nil {
			return err
		} else if genesisCreated {
			return nil
		}
	} else {
		return nil
	}

	// create chain id
	chainID := initConfig.chainID
	if chainID == "" {
		chainID = fmt.Sprintf("themis-%v", common.RandStr(6))
	}

	// get pubkey
	newPubkey := CryptoKeyToPubkey(valPubKey)

	// create validator account
	validator := hmTypes.NewValidator(hmTypes.NewValidatorID(uint64(initConfig.validatorID)),
		0, 0, 1, 1, newPubkey,
		hmTypes.BytesToThemisAddress(valPubKey.Address().Bytes()))

	vals := []*hmTypes.Validator{validator}
	validatorSet := hmTypes.NewValidatorSet(vals)

	// create genesis state
	appStateBytes := app.NewDefaultGenesisState()

	// auth state change
	appStateBytes, err = authTypes.SetGenesisStateToAppState(
		appStateBytes,
		[]authTypes.GenesisAccount{getGenesisAccount(validator.Signer.Bytes())},
	)
	if err != nil {
		return err
	}

	// staking state change
	appStateBytes, err = stakingTypes.SetGenesisStateToAppState(appStateBytes, vals, *validatorSet)
	if err != nil {
		return err
	}

	// chain manager state change
	appStateBytes, err = stakingTypes.SetGenesisStateToAppState(appStateBytes, vals, *validatorSet)
	if err != nil {
		return err
	}

	// metis state change
	appStateBytes, err = metisTypes.SetGenesisStateToAppState(appStateBytes, *validatorSet, 0, metisTypes.DefaultSpanDuration-1, metisTypes.DefaultFirstSpanDuration)
	if err != nil {
		return err
	}

	// mpc state change
	parties, err := hmTypes.NewMpcPartyIDFromPrivatekey(cdc, []crypto.PrivKey{valPriKey})
	if err != nil {
		return err
	}
	appStateBytes, err = mpcTypes.SetGenesisStateToAppState(appStateBytes, parties)
	if err != nil {
		return err
	}

	// app state json
	appStateJSON, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(appStateBytes)
	if err != nil {
		return err
	}

	toPrint := struct {
		ChainID string `json:"chain_id"`
		NodeID  string `json:"node_id"`
	}{
		chainID,
		nodeID,
	}

	out, err := codec.MarshalJSONIndent(cdc, toPrint)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "%s\n", string(out))

	return writeGenesisFile(tmtime.Now(), config.GenesisFile(), chainID, appStateJSON)
}

// InitCmd initialises files required to start themis
func initCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize genesis config, priv-validator file, and p2p-node file",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			initConfig := &initThemisConfig{
				chainID:     viper.GetString(client.FlagChainID),
				chain:       viper.GetString(helper.ChainFlag),
				validatorID: viper.GetInt64(stakingcli.FlagValidatorID),
				clientHome:  viper.GetString(helper.FlagClientHome),
				forceInit:   viper.GetBool(helper.OverwriteGenesisFlag),
			}
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			return themisInit(ctx, cdc, initConfig, config)
		},
	}

	cmd.Flags().String(cli.HomeFlag, helper.DefaultNodeHome, "node's home directory")
	cmd.Flags().String(helper.FlagClientHome, helper.DefaultCLIHome, "client's home directory")
	cmd.Flags().String(client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(helper.ChainFlag, helper.DefaultChain, "genesis file chain")
	cmd.Flags().Int(stakingcli.FlagValidatorID, 1, "--id=<validator ID here>, if left blank will be assigned 1")
	cmd.Flags().Bool(helper.OverwriteGenesisFlag, false, "overwrite the genesis.json file if it exists")

	return cmd
}

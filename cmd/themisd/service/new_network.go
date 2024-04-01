package service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/common"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/metis-seq/themis/app"
	authTypes "github.com/metis-seq/themis/auth/types"
	"github.com/metis-seq/themis/helper"
	metisTypes "github.com/metis-seq/themis/metis/types"
	mpcTypes "github.com/metis-seq/themis/mpc/types"
	stakingcli "github.com/metis-seq/themis/staking/client/cli"
	stakingTypes "github.com/metis-seq/themis/staking/types"

	chainManagerTypes "github.com/metis-seq/themis/chainmanager/types"
	hmTypes "github.com/metis-seq/themis/types"
)

// TestnetCmd initialises files required to start themis network
func newNetworkCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-network",
		Short: "Initialize files for a Themis network",
		Long: `network will create "v" + "n" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.
Optionally, it will fill in persistent_peers list in config file using either hostnames or IPs.

Example:
network --v 4 --n 8 --output-dir ./output 
`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			config.Consensus = defaultConsensusConfig()

			// update config consensus
			config.Consensus.CreateEmptyBlocks = true

			outDir := viper.GetString(flagOutputDir)

			// create chain id
			chainID := viper.GetString(client.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("themis-%v", common.RandStr(6))
			}

			// create chain id
			chain := viper.GetString(helper.ChainFlag)
			if chain == "" {
				chain = "mainnet"
			}

			// num of validators = validators in genesis files
			numValidators := viper.GetInt(flagNumValidators)

			// get total number of validators to be generated
			totalValidators := totalValidators()

			// first validators start ID
			startID := viper.GetInt64(stakingcli.FlagValidatorID)
			if startID == 0 {
				startID = 1
			}

			// signers data to dump in the signer-dump file
			signers := make([]hmTypes.ValidatorAccountFormatter, totalValidators)

			// Initialise variables for all validators
			nodeIDs := make([]string, totalValidators)
			valPubKeys := make([]crypto.PubKey, totalValidators)
			privKeys := make([]crypto.PrivKey, totalValidators)
			validators := make([]*hmTypes.Validator, numValidators)

			genFiles := make([]string, totalValidators)
			var err error

			nodeDaemonHomeName := viper.GetString(flagNodeDaemonHome)
			nodeCliHomeName := viper.GetString(flagNodeCliHome)

			// get genesis time
			genesisTime := tmtime.Now()

			for i := 0; i < totalValidators; i++ {
				config.Moniker = fmt.Sprintf("node-%d", i)

				// get node dir name = PREFIX+INDEX
				nodeDirName := fmt.Sprintf("%s%d", viper.GetString(flagNodeDirPrefix), i)

				// generate node and client dir
				nodeDir := filepath.Join(outDir, nodeDirName, nodeDaemonHomeName)
				clientDir := filepath.Join(outDir, nodeDirName, nodeCliHomeName)

				// set root in config
				config.SetRoot(nodeDir)

				// create config folder
				err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm)
				if err != nil {
					_ = os.RemoveAll(outDir)
					return err
				}

				err = os.MkdirAll(clientDir, nodeDirPerm)
				if err != nil {
					_ = os.RemoveAll(outDir)
					return err
				}

				nodeIDs[i], valPubKeys[i], privKeys[i], err = InitializeNodeValidatorFiles(config)
				if err != nil {
					return err
				}

				genFiles[i] = config.GenesisFile()
				newPubkey := CryptoKeyToPubkey(valPubKeys[i])

				if i < numValidators {
					// create validator account
					validators[i] = hmTypes.NewValidator(
						hmTypes.NewValidatorID(uint64(startID+int64(i))),
						0,
						0,
						1,
						20000,
						newPubkey,
						hmTypes.BytesToThemisAddress(valPubKeys[i].Address().Bytes()),
					)
				}

				signers[i] = GetSignerInfo(valPubKeys[i], privKeys[i].Bytes(), cdc)

				defaultThemisConfig := helper.GetDefaultThemisConfig()
				defaultThemisConfig.Chain = chain
				WriteDefaultThemisConfig(filepath.Join(config.RootDir, "config/themis-config.toml"), defaultThemisConfig)
			}

			// other data
			accounts := make([]authTypes.GenesisAccount, totalValidators)
			for i := 0; i < totalValidators; i++ {
				populatePersistentPeersInConfigAndWriteIt(config)
				// genesis account
				accounts[i] = getGenesisAccount(valPubKeys[i].Address().Bytes())
			}
			validatorSet := hmTypes.NewValidatorSet(validators)

			// new app state
			appStateBytes := app.NewDefaultGenesisState()

			// auth state change
			appStateBytes, err = authTypes.SetGenesisStateToAppState(appStateBytes, accounts)
			if err != nil {
				return err
			}

			// staking state change
			appStateBytes, err = stakingTypes.SetGenesisStateToAppState(appStateBytes, validators, *validatorSet)
			if err != nil {
				return err
			}

			// mpc state change
			parties, err := hmTypes.NewMpcPartyIDFromPrivatekey(cdc, privKeys)
			if err != nil {
				return err
			}
			appStateBytes, err = mpcTypes.SetGenesisStateToAppState(appStateBytes, parties)
			if err != nil {
				return err
			}

			// chain manager config
			appStateBytes, err = chainManagerTypes.SetGenesisStateToAppState(appStateBytes, chainManagerTypes.Params{
				MainchainTxConfirmations:  viper.GetUint64(flagL1Confirmation),
				MetischainTxConfirmations: viper.GetUint64(flagL2Confirmation),
				ChainParams: chainManagerTypes.ChainParams{
					MainChainID:           viper.GetString(flagL1ChainId),
					MetisChainID:          viper.GetString(flagL2ChainId),
					MetisTokenAddress:     hmTypes.HexToThemisAddress(viper.GetString(flagMetisAddr)),
					StakingManagerAddress: hmTypes.HexToThemisAddress(viper.GetString(flagLockingPool)),
					StakingInfoAddress:    hmTypes.HexToThemisAddress(viper.GetString(flagLockingInfo)),
					ValidatorSetAddress:   hmTypes.HexToThemisAddress(viper.GetString(flagSeqSet)),
				},
			})
			if err != nil {
				return err
			}

			// metis state change
			appStateBytes, err = metisTypes.SetGenesisStateToAppState(appStateBytes, *validatorSet,
				viper.GetUint64(flagEpoch0Start),
				viper.GetUint64(flagEpoch0End),
				viper.GetUint64(flagEpochLength),
			)
			if err != nil {
				return err
			}

			appStateJSON, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(appStateBytes)
			if err != nil {
				return err
			}

			for i := 0; i < totalValidators; i++ {
				if err = writeGenesisFile(genesisTime, genFiles[i], chainID, appStateJSON); err != nil {
					return err
				}
			}

			// dump signer information in a json file
			// TODO move to const string flag
			dump := viper.GetBool("signer-dump")
			if dump {
				signerJSON, err := jsoniter.ConfigFastest.MarshalIndent(signers, "", "  ")
				if err != nil {
					return err
				}

				if err := common.WriteFileAtomic(filepath.Join(outDir, "signer-dump.json"), signerJSON, 0600); err != nil {
					fmt.Println("Error writing signer-dump", err)
					return err
				}
			}

			fmt.Printf("Successfully initialized %d node directories\n", totalValidators)
			return nil
		},
	}

	cmd.Flags().Int(flagNumValidators, 4,
		"Number of validators to initialize the network with",
	)

	cmd.Flags().Int(flagNumNonValidators, 8,
		"Number of non-validators to initialize the network with",
	)

	cmd.Flags().StringP(flagOutputDir, "o", "./init-net",
		"Directory to store initialization data for the network",
	)

	cmd.Flags().String(flagNodeDirPrefix, "node",
		"Prefix the directory name for each node with (node results in node0, node1, ...)",
	)

	cmd.Flags().String(flagNodeDaemonHome, "themis",
		"Home directory of the node's daemon configuration",
	)

	cmd.Flags().String(flagNodeCliHome, "themiscli",
		"Home directory of the node's cli configuration",
	)

	cmd.Flags().String(flagNodeHostPrefix, "node",
		"Hostname prefix (node results in persistent peers list ID0@node0:26656, ID1@node1:26656, ...)")

	cmd.Flags().String(client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(helper.ChainFlag, helper.DefaultChain, "themis config chain")
	cmd.Flags().Bool("signer-dump", true, "dumps all signer information in a json file")
	cmd.Flags().Uint64(flagL1Confirmation, 12, "l1 confirmation number")
	cmd.Flags().Uint64(flagL2Confirmation, 1, "l2 confirmation number")
	cmd.Flags().String(flagL1ChainId, "1", "l1 chain id")
	cmd.Flags().String(flagL2ChainId, "1088", "l2 chain id")
	cmd.Flags().String(flagMetisAddr, "0x9e32b13ce7f2e80a01932b42553652e053d6ed8e", "the metis token address")
	cmd.Flags().String(flagLockingPool, "", "the LockingPool contract")
	cmd.Flags().String(flagLockingInfo, "", "the LockingInfo contract")
	cmd.Flags().String(flagSeqSet, "", "the SequencerSet contract")
	cmd.Flags().Uint64(flagEpoch0Start, 0, "the start block of epoch 0")
	cmd.Flags().Uint64(flagEpoch0End, metisTypes.DefaultSpanDuration-1, "the start block of epoch 0")
	cmd.Flags().Uint64(flagEpochLength, metisTypes.DefaultSpanDuration, "the start block of epoch 0")

	return cmd
}

// defaultConsensusConfig returns a default configuration for the consensus service
func defaultConsensusConfig() *config.ConsensusConfig {
	return &config.ConsensusConfig{
		WalPath:                     filepath.Join("data", "cs.wal", "wal"),
		TimeoutPropose:              500 * time.Millisecond,
		TimeoutProposeDelta:         300 * time.Millisecond,
		TimeoutPrevote:              500 * time.Millisecond,
		TimeoutPrevoteDelta:         500 * time.Millisecond,
		TimeoutPrecommit:            500 * time.Millisecond,
		TimeoutPrecommitDelta:       500 * time.Millisecond,
		TimeoutCommit:               500 * time.Millisecond,
		SkipTimeoutCommit:           false,
		CreateEmptyBlocks:           true,
		CreateEmptyBlocksInterval:   0 * time.Second,
		PeerGossipSleepDuration:     100 * time.Millisecond,
		PeerQueryMaj23SleepDuration: 2000 * time.Millisecond,
	}
}

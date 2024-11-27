package helper

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	logger "github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/privval"

	"github.com/metis-seq/themis/file"
	"github.com/metis-seq/themis/types"

	cfg "github.com/tendermint/tendermint/config"
	tmTypes "github.com/tendermint/tendermint/types"
)

const (
	TendermintNodeFlag   = "node"
	WithThemisConfigFlag = "themis-config"
	HomeFlag             = "home"
	FlagClientHome       = "home-client"
	OverwriteGenesisFlag = "overwrite-genesis"
	RestServerFlag       = "rest-server"
	BridgeFlag           = "bridge"
	LogLevel             = "log_level"
	LogsWriterFileFlag   = "logs_writer_file"
	SeedsFlag            = "seeds"

	MainChain  = "mainnet"
	TestChain  = "testnet"
	LocalChain = "local"

	// themis-config flags
	MainRPCUrlFlag           = "eth_rpc_url"
	MetisRPCUrlFlag          = "metis_rpc_url"
	TendermintNodeURLFlag    = "tendermint_rpc_url"
	ThemisServerURLFlag      = "themis_rest_server"
	MpcServerURLFlag         = "mpc_rpc_server"
	RootPollIntervalFlag     = "root_poll_interval"
	ThemisPollIntervalFlag   = "themis_poll_interval"
	MetisPollIntervalFlag    = "metis_poll_interval"
	SpanPollIntervalFlag     = "span_poll_interval"
	ReSpanPollIntervalFlag   = "respan_poll_interval"
	ReSpanDelayTimeFlag      = "respan_delay_time"
	MpcPollIntervalFlag      = "mpc_poll_interval"
	MainchainGasLimitFlag    = "main_chain_gas_limit"
	MainchainMaxGasPriceFlag = "main_chain_max_gas_price"
	ChainFlag                = "chain"

	// ---
	// TODO Move these to common client flags
	// BroadcastBlock defines a tx broadcasting mode where the client waits for
	// the tx to be committed in a block.
	BroadcastBlock = "block"

	// BroadcastSync defines a tx broadcasting mode where the client waits for
	// a CheckTx execution response only.
	BroadcastSync = "sync"

	// BroadcastAsync defines a tx broadcasting mode where the client returns
	// immediately.
	BroadcastAsync = "async"
	// --

	// RPC Endpoints
	DefaultMainRPCUrl  = "http://localhost:9545"
	DefaultMetisRPCUrl = "http://localhost:8545"

	// RPC Timeouts
	DefaultEthRPCTimeout   = 10 * time.Second
	DefaultMetisRPCTimeout = 10 * time.Second

	// Services
	DefaultThemisServerURL   = "http://0.0.0.0:1317"
	DefaultMpcServerURL      = "tcp://0.0.0.0:9091"
	DefaultTendermintNodeURL = "http://0.0.0.0:26657"

	DefaultRootSyncerPollInterval   = 15 * time.Second
	DefaultThemisSyncerPollInterval = 1 * time.Second
	DefaultMetisSyncerPollInterval  = 1 * time.Second
	DefaultSpanPollInterval         = 10 * time.Second
	DefaultReSpanPollInterval       = 10 * time.Second
	DefaultReSpanDelayTime          = 3 * time.Minute
	DefaultMpcPollInterval          = 10 * time.Second
	DefaultEnableSH                 = false
	DefaultSHStateSyncedInterval    = 15 * time.Minute
	DefaultSHStakeUpdateInterval    = 3 * time.Hour
	DefaultSHMaxDepthDuration       = time.Hour

	DefaultMainchainGasLimit  = uint64(5000000)
	DefaultMetischainGasLimit = uint64(5000000)

	MainchainBuildGasLimit  = uint64(30_000_000)
	MainnetNetwork          = "andromeda"
	MainnetForkMaxGasHeight = 695000
	MainnetForkMaxGasLimit  = uint64(30_000_000)

	DefaultMainchainMaxGasPrice = 400000000000 // 400 Gwei

	DefaultMainChainID  = "1"
	DefaultMetisChainID = "1088"

	DefaultLogsType = "json"
	DefaultChain    = MainChain

	DefaultMetricsListenAddr = ":2112"
	DefaultRPCListenAddr     = ":8646"

	DefaultTendermintNode = "tcp://localhost:26657"

	DefaultMainnetSeeds = ""

	DefaultTestnetSeeds = ""

	secretFilePerm = 0600

	// Legacy value - DO NOT CHANGE
	// Maximum allowed event record data size
	LegacyMaxStateSyncSize = 100000

	// New max state sync size after hardfork
	MaxStateSyncSize = 30000

	// Default Open Collector Endpoint
	DefaultOpenCollectorEndpoint = "localhost:4317"
)

var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.themiscli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.themisd")
)

var cdc = amino.NewCodec()

func init() {
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{}, secp256k1.PubKeyAminoName, nil)
	cdc.RegisterConcrete(secp256k1.PrivKeySecp256k1{}, secp256k1.PrivKeyAminoName, nil)

	Logger = logger.NewTMLogger(logger.NewSyncWriter(os.Stdout))
}

// Configuration represents themis config
type Configuration struct {
	EthRPCUrl        string `mapstructure:"eth_rpc_url"`        // RPC endpoint for main chain
	MetisRPCUrl      string `mapstructure:"metis_rpc_url"`      // RPC endpoint for metis chain
	TendermintRPCUrl string `mapstructure:"tendermint_rpc_url"` // tendemint node url

	EthRPCTimeout   time.Duration `mapstructure:"eth_rpc_timeout"`   // timeout for eth rpc
	MetisRPCTimeout time.Duration `mapstructure:"metis_rpc_timeout"` // timeout for metis rpc

	ThemisServerURL string `mapstructure:"themis_rest_server"` // themis server url
	MpcServerURL    string `mapstructure:"mpc_rpc_server"`     // mpc server url

	MainchainGasLimit    uint64 `mapstructure:"main_chain_gas_limit"`     // gas limit to mainchain transaction.
	MainchainMaxGasPrice int64  `mapstructure:"main_chain_max_gas_price"` // max gas price to mainchain transaction.

	// config related to bridge
	RootPollInterval      time.Duration `mapstructure:"root_poll_interval"`   // Poll interval for syncher service to sync for changes on main chain
	ThemisPollInterval    time.Duration `mapstructure:"themis_poll_interval"` // Poll interval for syncher service to sync for changes on themis chain
	MetisPollInterval     time.Duration `mapstructure:"metis_poll_interval"`  // Poll interval for syncher service to sync for changes on metis chain
	SpanPollInterval      time.Duration `mapstructure:"span_poll_interval"`
	ReSpanPollInterval    time.Duration `mapstructure:"respan_poll_interval"`
	ReSpanDelayTime       time.Duration `mapstructure:"respan_delay_time"`
	MpcPollInterval       time.Duration `mapstructure:"mpc_poll_interval"`
	EnableSH              bool          `mapstructure:"enable_self_heal"`         // Enable self healing
	SHStakeUpdateInterval time.Duration `mapstructure:"sh_stake_update_interval"` // Interval to self-heal StakeUpdate events if missing
	SHMaxDepthDuration    time.Duration `mapstructure:"sh_max_depth_duration"`    // Max duration that allows to suggest self-healing is not needed

	// Log related options
	LogsType       string `mapstructure:"logs_type"`        // if true, enable logging in json format
	LogsWriterFile string `mapstructure:"logs_writer_file"` // if given, Logs will be written to this file else os.Stdout

	// current chain - newSelectionAlgoHeight depends on this
	Chain string `mapstructure:"chain"`
}

var conf Configuration

// MainChainClient stores eth clie nt for Main chain Network
var mainChainClient *ethclient.Client
var mainRPCClient *rpc.Client

// MetisClient stores eth/rpc client for Metis Network
var metisClient *ethclient.Client
var metisRPCClient *rpc.Client

// private key object
var privObject secp256k1.PrivKeySecp256k1

var pubObject secp256k1.PubKeySecp256k1

// Logger stores global logger object
var Logger logger.Logger

// GenesisDoc contains the genesis file
var GenesisDoc tmTypes.GenesisDoc

var newSelectionAlgoHeight int64 = 0

var spanOverrideHeight int64 = 0

// InitThemisConfig initializes with viper config (from themis configuration)
func InitThemisConfig(homeDir string) {
	if strings.Compare(homeDir, "") == 0 {
		// get home dir from viper
		homeDir = viper.GetString(HomeFlag)
	}

	// get themis config filepath from viper/cobra flag
	themisConfigFileFromFlag := viper.GetString(WithThemisConfigFlag)

	// init themis with changed config files
	InitThemisConfigWith(homeDir, themisConfigFileFromFlag)

	// init mpc client
	mpcConnect()
}

// InitThemisConfigWith initializes passed themis/tendermint config files
func InitThemisConfigWith(homeDir string, themisConfigFileFromFLag string) {
	if strings.Compare(homeDir, "") == 0 {
		return
	}

	// read configuration from the standard configuration file
	configDir := filepath.Join(homeDir, "config")
	themisViper := viper.New()
	themisViper.SetEnvPrefix("THEMIS")
	themisViper.AutomaticEnv()

	if themisConfigFileFromFLag == "" {
		themisViper.SetConfigName("themis-config") // name of config file (without extension)
		themisViper.AddConfigPath(configDir)       // call multiple times to add many search paths
	} else {
		themisViper.SetConfigFile(themisConfigFileFromFLag) // set config file explicitly
	}

	// Handle errors reading the config file
	if err := themisViper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	// unmarshal configuration from the standard configuration file
	if err := themisViper.UnmarshalExact(&conf); err != nil {
		log.Fatalln("Unable to unmarshall config", "Error", err)
	}

	//  if there is a file with overrides submitted via flags => read it an merge it with the alreadey read standard configuration
	if themisConfigFileFromFLag != "" {
		themisViperFromFlag := viper.New()
		themisViperFromFlag.SetConfigFile(themisConfigFileFromFLag) // set flag config file explicitly

		err := themisViperFromFlag.ReadInConfig()
		if err != nil { // Handle errors reading the config file sybmitted as a flag
			log.Fatalln("Unable to read config file submitted via flag", "Error", err)
		}

		var confFromFlag Configuration
		// unmarshal configuration from the configuration file submited as a flag
		if err = themisViperFromFlag.UnmarshalExact(&confFromFlag); err != nil {
			log.Fatalln("Unable to unmarshall config file submitted via flag", "Error", err)
		}

		conf.Merge(&confFromFlag)
	}

	// update configuration data with submitted flags
	if err := conf.UpdateWithFlags(viper.GetViper(), Logger); err != nil {
		log.Fatalln("Unable to read flag values. Check log for details.", "Error", err)
	}

	// perform check for json logging
	if conf.LogsType == "json" {
		Logger = logger.NewTMJSONLogger(logger.NewSyncWriter(GetLogsWriter(conf.LogsWriterFile)))
	} else {
		// default fallback
		Logger = logger.NewTMLogger(logger.NewSyncWriter(GetLogsWriter(conf.LogsWriterFile)))
	}

	// perform checks for timeout
	if conf.EthRPCTimeout == 0 {
		// fallback to default
		// Logger.Debug("Invalid ETH RPC timeout provided, falling back to default value", "timeout", DefaultEthRPCTimeout)
		conf.EthRPCTimeout = DefaultEthRPCTimeout
	}

	if conf.MetisRPCTimeout == 0 {
		// fallback to default
		// Logger.Debug("Invalid Metis RPC timeout provided, falling back to default value", "timeout", DefaultMetisRPCTimeout)
		conf.MetisRPCTimeout = DefaultMetisRPCTimeout
	}

	if conf.SHStakeUpdateInterval == 0 {
		// fallback to default
		Logger.Debug("Invalid self-healing StakeUpdate interval provided, falling back to default value", "interval", DefaultSHStakeUpdateInterval)
		conf.SHStakeUpdateInterval = DefaultSHStakeUpdateInterval
	}

	if conf.SHMaxDepthDuration == 0 {
		// fallback to default
		Logger.Debug("Invalid self-healing max depth duration provided, falling back to default value", "duration", DefaultSHMaxDepthDuration)
		conf.SHMaxDepthDuration = DefaultSHMaxDepthDuration
	}

	// env over write
	conf.MergeFromEnv()

	var err error
	Logger.Info("dailing eth client", "url", conf.EthRPCUrl)
	if mainRPCClient, err = rpc.Dial(conf.EthRPCUrl); err != nil {
		log.Fatalln("Unable to dial via ethClient", "URL=", conf.EthRPCUrl, "chain=eth", "Error", err)
	}

	mainChainClient = ethclient.NewClient(mainRPCClient)

	Logger.Info("dailing metis client", "url", conf.MetisRPCUrl)
	if metisRPCClient, err = rpc.Dial(conf.MetisRPCUrl); err != nil {
		log.Fatal(err)
	}

	metisClient = ethclient.NewClient(metisRPCClient)
	// Loading genesis doc
	genDoc, err := tmTypes.GenesisDocFromFile(filepath.Join(configDir, "genesis.json"))
	if err != nil {
		log.Fatal(err)
	}

	GenesisDoc = *genDoc

	// load pv file, unmarshall and set to privObject
	err = file.PermCheck(file.Rootify("priv_validator_key.json", configDir), secretFilePerm)
	if err != nil {
		Logger.Error(err.Error())
	}

	privVal := privval.LoadFilePV(filepath.Join(configDir, "priv_validator_key.json"), filepath.Join(configDir, "priv_validator_key.json"))
	cdc.MustUnmarshalBinaryBare(privVal.Key.PrivKey.Bytes(), &privObject)
	cdc.MustUnmarshalBinaryBare(privObject.PubKey().Bytes(), &pubObject)

	switch conf.Chain {
	case MainChain:
	case TestChain:
	default:
		newSelectionAlgoHeight = 0
		spanOverrideHeight = 0
	}

}

// GetDefaultThemisConfig returns configration with default params
func GetDefaultThemisConfig() Configuration {
	defaultConf := Configuration{
		EthRPCUrl:        DefaultMainRPCUrl,
		MetisRPCUrl:      DefaultMetisRPCUrl,
		TendermintRPCUrl: DefaultTendermintNodeURL,

		EthRPCTimeout:   DefaultEthRPCTimeout,
		MetisRPCTimeout: DefaultMetisRPCTimeout,

		ThemisServerURL: DefaultThemisServerURL,
		MpcServerURL:    DefaultMpcServerURL,

		MainchainGasLimit:     DefaultMainchainGasLimit,
		MainchainMaxGasPrice:  DefaultMainchainMaxGasPrice,
		RootPollInterval:      DefaultRootSyncerPollInterval,
		ThemisPollInterval:    DefaultThemisSyncerPollInterval,
		MetisPollInterval:     DefaultMetisSyncerPollInterval,
		SpanPollInterval:      DefaultSpanPollInterval,
		ReSpanPollInterval:    DefaultReSpanPollInterval,
		ReSpanDelayTime:       DefaultReSpanDelayTime,
		MpcPollInterval:       DefaultMpcPollInterval,
		EnableSH:              DefaultEnableSH,
		SHStakeUpdateInterval: DefaultSHStakeUpdateInterval,
		SHMaxDepthDuration:    DefaultSHMaxDepthDuration,

		LogsType:       DefaultLogsType,
		Chain:          DefaultChain,
		LogsWriterFile: "", // default to stdout
	}

	// env over write
	defaultConf.MergeFromEnv()

	return defaultConf
}

func (c *Configuration) MergeFromEnv() {
	// env over write
	envEthRpcUrl := os.Getenv("ETH_RPC_URL")
	if envEthRpcUrl != "" {
		c.EthRPCUrl = envEthRpcUrl
	}

	envThemisRestUrl := os.Getenv("REST_SERVER")
	if envThemisRestUrl != "" {
		c.ThemisServerURL = envThemisRestUrl
	}

	envThemisTendermintUrl := os.Getenv("TENDERMINT_RPC_URL")
	if envThemisTendermintUrl != "" {
		c.TendermintRPCUrl = envThemisTendermintUrl
	}

	envMetisRpcUrl := os.Getenv("METIS_RPC_URL")
	if envMetisRpcUrl != "" {
		c.MetisRPCUrl = envMetisRpcUrl
	}

	envMpcRpcUrl := os.Getenv("MPC_RPC_URL")
	if envMpcRpcUrl != "" {
		c.MpcServerURL = envMpcRpcUrl
	}

	envSpanPollInterval := os.Getenv("SPAN_POLL_INTERNAL")
	if envSpanPollInterval != "" {
		var err error
		c.SpanPollInterval, err = time.ParseDuration(envSpanPollInterval)
		if err != nil {
			panic("invalid SPAN_POLL_INTERNAL" + err.Error())
		}
	}
	envReSpanPollInterval := os.Getenv("RESPAN_POLL_INTERNAL")
	if envReSpanPollInterval != "" {
		var err error
		c.ReSpanPollInterval, err = time.ParseDuration(envReSpanPollInterval)
		if err != nil {
			panic("invalid RESPAN_POLL_INTERNAL" + err.Error())
		}
	}

	envReSpanDelayTime := os.Getenv("RESPAN_DELAY_TIME")
	if envReSpanDelayTime != "" {
		var err error
		c.ReSpanDelayTime, err = time.ParseDuration(envReSpanDelayTime)
		if err != nil {
			panic("invalid RESPAN_DELAY_TIME" + err.Error())
		}
	}

	envMpcPollInterval := os.Getenv("MPC_POLL_INTERNAL")
	if envMpcPollInterval != "" {
		var err error
		c.MpcPollInterval, err = time.ParseDuration(envMpcPollInterval)
		if err != nil {
			panic("invalid MPC_POLL_INTERNAL" + err.Error())
		}
	}
}

// GetConfig returns cached configuration object
func GetConfig() Configuration {
	return conf
}

func GetGenesisDoc() tmTypes.GenesisDoc {
	return GenesisDoc
}

// TEST PURPOSE ONLY
// SetTestConfig sets test configuration
func SetTestConfig(_conf Configuration) {
	conf = _conf
}

//
// Get main/metis clients
//

// GetMainChainRPCClient returns main chain RPC client
func GetMainChainRPCClient() *rpc.Client {
	return mainRPCClient
}

// GetMainClient returns main chain's eth client
func GetMainClient() *ethclient.Client {
	return mainChainClient
}

// GetMetisClient returns metis's eth client
func GetMetisClient() *ethclient.Client {
	return metisClient
}

// GetMetisRPCClient returns metis's RPC client
func GetMetisRPCClient() *rpc.Client {
	return metisRPCClient
}

// GetPrivKey returns priv key object
func GetPrivKey() secp256k1.PrivKeySecp256k1 {
	return privObject
}

// GetECDSAPrivKey return ecdsa private key
func GetECDSAPrivKey() *ecdsa.PrivateKey {
	// get priv key
	pkObject := GetPrivKey()

	// create ecdsa private key
	ecdsaPrivateKey, _ := ethCrypto.ToECDSA(pkObject[:])

	return ecdsaPrivateKey
}

// GetPubKey returns pub key object
func GetPubKey() secp256k1.PubKeySecp256k1 {
	return pubObject
}

// GetAddress returns address object
func GetAddress() []byte {
	return GetPubKey().Address().Bytes()
}

// GetAddress returns address object
func GetAddressStr() string {
	return "0x" + GetPubKey().Address().String()
}

// GetValidChains returns all the valid chains
func GetValidChains() []string {
	return []string{"mainnet", "testnet", "local"}
}

// GetNewSelectionAlgoHeight returns newSelectionAlgoHeight
func GetNewSelectionAlgoHeight() int64 {
	return newSelectionAlgoHeight
}

// GetSpanOverrideHeight returns spanOverrideHeight
func GetSpanOverrideHeight() int64 {
	return spanOverrideHeight
}

// DecorateWithThemisFlags adds persistent flags for themis-config and bind flags with command
func DecorateWithThemisFlags(cmd *cobra.Command, v *viper.Viper, loggerInstance logger.Logger, caller string) {
	// add with-themis-config flag
	cmd.PersistentFlags().String(
		WithThemisConfigFlag,
		"",
		"Override of Themis config file (default <home>/config/themis-config.json)",
	)

	if err := v.BindPFlag(WithThemisConfigFlag, cmd.PersistentFlags().Lookup(WithThemisConfigFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, WithThemisConfigFlag), "Error", err)
	}

	// add MainRPCUrlFlag flag
	cmd.PersistentFlags().String(
		MainRPCUrlFlag,
		"",
		"Set RPC endpoint for ethereum chain",
	)

	if err := v.BindPFlag(MainRPCUrlFlag, cmd.PersistentFlags().Lookup(MainRPCUrlFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, MainRPCUrlFlag), "Error", err)
	}

	// add MetisRPCUrlFlag flag
	cmd.PersistentFlags().String(
		MetisRPCUrlFlag,
		"",
		"Set RPC endpoint for metis chain",
	)

	if err := v.BindPFlag(MetisRPCUrlFlag, cmd.PersistentFlags().Lookup(MetisRPCUrlFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, MetisRPCUrlFlag), "Error", err)
	}

	// add TendermintNodeURLFlag flag
	cmd.PersistentFlags().String(
		TendermintNodeURLFlag,
		"",
		"Set RPC endpoint for tendermint",
	)

	if err := v.BindPFlag(TendermintNodeURLFlag, cmd.PersistentFlags().Lookup(TendermintNodeURLFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, TendermintNodeURLFlag), "Error", err)
	}

	// add ThemisServerURLFlag flag
	cmd.PersistentFlags().String(
		ThemisServerURLFlag,
		"",
		"Set Themis REST server endpoint",
	)

	if err := v.BindPFlag(ThemisServerURLFlag, cmd.PersistentFlags().Lookup(ThemisServerURLFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, ThemisServerURLFlag), "Error", err)
	}

	// add MpcServerURLFlag flag
	cmd.PersistentFlags().String(
		MpcServerURLFlag,
		"",
		"Set mpc server endpoint",
	)

	if err := v.BindPFlag(MpcServerURLFlag, cmd.PersistentFlags().Lookup(MpcServerURLFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, MpcServerURLFlag), "Error", err)
	}

	// add RootPollIntervalFlag flag
	cmd.PersistentFlags().String(
		RootPollIntervalFlag,
		"",
		"Set root pull interval",
	)

	if err := v.BindPFlag(RootPollIntervalFlag, cmd.PersistentFlags().Lookup(RootPollIntervalFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, RootPollIntervalFlag), "Error", err)
	}

	// add ThemisPollIntervalFlag flag
	cmd.PersistentFlags().String(
		ThemisPollIntervalFlag,
		"",
		"Set themis pull interval",
	)

	if err := v.BindPFlag(ThemisPollIntervalFlag, cmd.PersistentFlags().Lookup(ThemisPollIntervalFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, ThemisPollIntervalFlag), "Error", err)
	}

	// add MetisPollIntervalFlag flag
	cmd.PersistentFlags().String(
		MetisPollIntervalFlag,
		"",
		"Set metis pull interval",
	)

	if err := v.BindPFlag(MetisPollIntervalFlag, cmd.PersistentFlags().Lookup(MetisPollIntervalFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, MetisPollIntervalFlag), "Error", err)
	}

	// add SpanPollIntervalFlag flag
	cmd.PersistentFlags().String(
		SpanPollIntervalFlag,
		"",
		"Set span pull interval",
	)

	if err := v.BindPFlag(SpanPollIntervalFlag, cmd.PersistentFlags().Lookup(SpanPollIntervalFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, SpanPollIntervalFlag), "Error", err)
	}

	// add ReSpanPollIntervalFlag flag
	cmd.PersistentFlags().String(
		ReSpanPollIntervalFlag,
		"",
		"Set re-span interval",
	)

	if err := v.BindPFlag(ReSpanPollIntervalFlag, cmd.PersistentFlags().Lookup(ReSpanPollIntervalFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, ReSpanPollIntervalFlag), "Error", err)
	}

	// add ReSpanDelayTimeFlag flag
	cmd.PersistentFlags().String(
		ReSpanDelayTimeFlag,
		"",
		"Set re-span delay time",
	)

	if err := v.BindPFlag(ReSpanDelayTimeFlag, cmd.PersistentFlags().Lookup(ReSpanDelayTimeFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, ReSpanDelayTimeFlag), "Error", err)
	}

	// add MpcPollIntervalFlag flag
	cmd.PersistentFlags().String(
		MpcPollIntervalFlag,
		"",
		"Set mpc pull interval",
	)

	if err := v.BindPFlag(MpcPollIntervalFlag, cmd.PersistentFlags().Lookup(MpcPollIntervalFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, MpcPollIntervalFlag), "Error", err)
	}

	// add MainchainGasLimitFlag flag
	cmd.PersistentFlags().Uint64(
		MainchainGasLimitFlag,
		0,
		"Set main chain gas limit",
	)

	if err := v.BindPFlag(MainchainGasLimitFlag, cmd.PersistentFlags().Lookup(MainchainGasLimitFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, MainchainGasLimitFlag), "Error", err)
	}

	// add MainchainMaxGasPriceFlag flag
	cmd.PersistentFlags().Int64(
		MainchainMaxGasPriceFlag,
		0,
		"Set main chain max gas limit",
	)

	if err := v.BindPFlag(MainchainMaxGasPriceFlag, cmd.PersistentFlags().Lookup(MainchainMaxGasPriceFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, MainchainMaxGasPriceFlag), "Error", err)
	}

	// add chain flag
	cmd.PersistentFlags().String(
		ChainFlag,
		"",
		fmt.Sprintf("Set one of the chains: [%s]", strings.Join(GetValidChains(), ",")),
	)

	if err := v.BindPFlag(ChainFlag, cmd.PersistentFlags().Lookup(ChainFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, ChainFlag), "Error", err)
	}

	// add logsWriterFile flag
	cmd.PersistentFlags().String(
		LogsWriterFileFlag,
		"",
		"Set logs writer file, Default is os.Stdout",
	)

	if err := v.BindPFlag(LogsWriterFileFlag, cmd.PersistentFlags().Lookup(LogsWriterFileFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, LogsWriterFileFlag), "Error", err)
	}
}

func (c *Configuration) UpdateWithFlags(v *viper.Viper, loggerInstance logger.Logger) error {
	const logErrMsg = "Unable to read flag."

	// get endpoint for ethereum chain from viper/cobra
	stringConfgValue := v.GetString(MainRPCUrlFlag)
	if stringConfgValue != "" {
		c.EthRPCUrl = stringConfgValue
	}

	// get endpoint for metis chain from viper/cobra
	stringConfgValue = v.GetString(MetisRPCUrlFlag)
	if stringConfgValue != "" {
		c.MetisRPCUrl = stringConfgValue
	}

	// get endpoint for tendermint from viper/cobra
	stringConfgValue = v.GetString(TendermintNodeURLFlag)
	if stringConfgValue != "" {
		c.TendermintRPCUrl = stringConfgValue
	}

	stringConfgValue = v.GetString(ThemisServerURLFlag)
	if stringConfgValue != "" {
		c.ThemisServerURL = stringConfgValue
	}

	// get mpc server endpoint from viper/cobra
	stringConfgValue = v.GetString(MpcServerURLFlag)
	if stringConfgValue != "" {
		c.MpcServerURL = stringConfgValue
	}

	// need this error for parsing Duration values
	var err error

	// get syncer pull interval from viper/cobra
	stringConfgValue = v.GetString(RootPollIntervalFlag)
	if stringConfgValue != "" {
		if c.RootPollInterval, err = time.ParseDuration(stringConfgValue); err != nil {
			loggerInstance.Error(logErrMsg, "Flag", RootPollIntervalFlag, "Error", err)
			return err
		}
	}

	// get syncer pull interval from viper/cobra
	stringConfgValue = v.GetString(ThemisPollIntervalFlag)
	if stringConfgValue != "" {
		if c.ThemisPollInterval, err = time.ParseDuration(stringConfgValue); err != nil {
			loggerInstance.Error(logErrMsg, "Flag", ThemisPollIntervalFlag, "Error", err)
			return err
		}
	}

	// get syncer pull interval from viper/cobra
	stringConfgValue = v.GetString(MetisPollIntervalFlag)
	if stringConfgValue != "" {
		if c.MetisPollInterval, err = time.ParseDuration(stringConfgValue); err != nil {
			loggerInstance.Error(logErrMsg, "Flag", MetisPollIntervalFlag, "Error", err)
			return err
		}
	}

	// get span poll interval from viper/cobra
	stringConfgValue = v.GetString(SpanPollIntervalFlag)
	if stringConfgValue != "" {
		if c.SpanPollInterval, err = time.ParseDuration(stringConfgValue); err != nil {
			loggerInstance.Error(logErrMsg, "Flag", SpanPollIntervalFlag, "Error", err)
			return err
		}
	}

	// get re-span poll interval from viper/cobra
	stringConfgValue = v.GetString(ReSpanPollIntervalFlag)
	if stringConfgValue != "" {
		if c.ReSpanPollInterval, err = time.ParseDuration(stringConfgValue); err != nil {
			loggerInstance.Error(logErrMsg, "Flag", ReSpanPollIntervalFlag, "Error", err)
			return err
		}
	}

	// get re-span delay time from viper/cobra
	stringConfgValue = v.GetString(ReSpanDelayTimeFlag)
	if stringConfgValue != "" {
		if c.ReSpanDelayTime, err = time.ParseDuration(stringConfgValue); err != nil {
			loggerInstance.Error(logErrMsg, "Flag", ReSpanDelayTimeFlag, "Error", err)
			return err
		}
	}

	// get mpc poll interval from viper/cobra
	stringConfgValue = v.GetString(MpcPollIntervalFlag)
	if stringConfgValue != "" {
		if c.MpcPollInterval, err = time.ParseDuration(stringConfgValue); err != nil {
			loggerInstance.Error(logErrMsg, "Flag", MpcPollIntervalFlag, "Error", err)
			return err
		}
	}

	// get mainchain gas limit from viper/cobra
	uint64ConfgValue := v.GetUint64(MainchainGasLimitFlag)
	if uint64ConfgValue != 0 {
		c.MainchainGasLimit = uint64ConfgValue
	}

	// get mainchain max gas price from viper/cobra. if it is greater then  zero => set it as configuration parameter
	int64ConfgValue := v.GetInt64(MainchainMaxGasPriceFlag)
	if int64ConfgValue > 0 {
		c.MainchainMaxGasPrice = int64ConfgValue
	}

	// get chain from viper/cobra flag
	stringConfgValue = v.GetString(ChainFlag)
	if stringConfgValue != "" {
		c.Chain = stringConfgValue
	}

	stringConfgValue = v.GetString(LogsWriterFileFlag)
	if stringConfgValue != "" {
		c.LogsWriterFile = stringConfgValue
	}

	return nil
}

func (c *Configuration) Merge(cc *Configuration) {
	if cc.EthRPCUrl != "" {
		c.EthRPCUrl = cc.EthRPCUrl
	}

	if cc.MetisRPCUrl != "" {
		c.MetisRPCUrl = cc.MetisRPCUrl
	}

	if cc.TendermintRPCUrl != "" {
		c.TendermintRPCUrl = cc.TendermintRPCUrl
	}

	if cc.ThemisServerURL != "" {
		c.ThemisServerURL = cc.ThemisServerURL
	}

	if cc.MpcServerURL != "" {
		c.MpcServerURL = cc.MpcServerURL
	}

	if cc.MainchainGasLimit != 0 {
		c.MainchainGasLimit = cc.MainchainGasLimit
	}

	if cc.MainchainMaxGasPrice != 0 {
		c.MainchainMaxGasPrice = cc.MainchainMaxGasPrice
	}

	if cc.RootPollInterval != 0 {
		c.RootPollInterval = cc.RootPollInterval
	}

	if cc.ThemisPollInterval != 0 {
		c.ThemisPollInterval = cc.ThemisPollInterval
	}

	if cc.MetisPollInterval != 0 {
		c.MetisPollInterval = cc.MetisPollInterval
	}

	if cc.SpanPollInterval != 0 {
		c.SpanPollInterval = cc.SpanPollInterval
	}

	if cc.ReSpanPollInterval != 0 {
		c.ReSpanPollInterval = cc.ReSpanPollInterval
	}

	if cc.ReSpanDelayTime != 0 {
		c.ReSpanDelayTime = cc.ReSpanDelayTime
	}

	if cc.MpcPollInterval != 0 {
		c.MpcPollInterval = cc.MpcPollInterval
	}

	if cc.Chain != "" {
		c.Chain = cc.Chain
	}

	if cc.LogsWriterFile != "" {
		c.LogsWriterFile = cc.LogsWriterFile
	}
}

// DecorateWithTendermintFlags creates tendermint flags for desired command and bind them to viper
func DecorateWithTendermintFlags(cmd *cobra.Command, v *viper.Viper, loggerInstance logger.Logger, message string) {
	// add seeds flag
	cmd.PersistentFlags().String(
		SeedsFlag,
		"",
		"Override seeds",
	)

	if err := v.BindPFlag(SeedsFlag, cmd.PersistentFlags().Lookup(SeedsFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", message, SeedsFlag), "Error", err)
	}
}

// UpdateTendermintConfig updates tenedermint config with flags and default values if needed
func UpdateTendermintConfig(tendermintConfig *cfg.Config, v *viper.Viper) {
	// update tendermintConfig.P2P.Seeds
	seedsFlagValue := v.GetString(SeedsFlag)
	if seedsFlagValue != "" {
		tendermintConfig.P2P.Seeds = seedsFlagValue
	}

	if tendermintConfig.P2P.Seeds == "" {
		switch conf.Chain {
		case MainChain:
			tendermintConfig.P2P.Seeds = DefaultMainnetSeeds
		case TestChain:
			tendermintConfig.P2P.Seeds = DefaultTestnetSeeds
		}
	}
}

func GetLogsWriter(logsWriterFile string) io.Writer {
	if logsWriterFile != "" {
		logWriter, err := os.OpenFile(logsWriterFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening log writer file: %v", err)
		}

		return logWriter
	} else {
		return os.Stdout
	}
}

type NotifyRespanStartResult struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result"`
}

func NotifyRespanStart(preSequencer, newSequencer string, startHeight uint64) error {
	payload := strings.NewReader(fmt.Sprintf(`{
    "jsonrpc": "2.0",
    "method": "rollupbridge_setPreRespan",
    "params": ["%v","%v",%v],
    "id": 0
}`, preSequencer, newSequencer, startHeight))

	client := &http.Client{}
	req, err := http.NewRequest("POST", conf.MetisRPCUrl, payload)
	if err != nil {
		Logger.Error("NotifyRespanStart http request failed", "err", err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		Logger.Error("NotifyRespanStart http request failed", "err", err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		Logger.Error("NotifyRespanStart http request failed", "err", err)
		return err
	}

	var result NotifyRespanStartResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		Logger.Error("NotifyRespanStart json unmarshal failed", "err", err)
		return err
	}
	Logger.Debug("NotifyRespanStart", "result", result)

	return nil
}

func CalcMetisSequencer(currentSequencer string, currentL2Height, currentBatch uint64, valSet *types.ValidatorSet) string {
	var allowedSequencers []string
	for _, val := range valSet.Validators {
		if val.EndBatch > 0 || !val.IsCurrentValidator(currentBatch) || strings.EqualFold(currentSequencer, val.Signer.EthAddress().Hex()) {
			continue
		}

		allowedSequencers = append(allowedSequencers, val.Signer.EthAddress().Hex())
	}
	sort.Strings(allowedSequencers)

	// make random sequencer
	rand.Seed(int64(currentL2Height))
	randomIndex := rand.Intn(len(allowedSequencers))
	newSequencer := allowedSequencers[randomIndex]

	Logger.Debug("CalcMetisSequencer", "allowedSequencersLen", len(allowedSequencers), "randomIndex", randomIndex, "newSequencer", newSequencer)
	return newSequencer
}

func CalcMetisSequencerWithSeed(currentSequencer string, currentL2Height, currentBatch uint64, valSet *types.ValidatorSet, seed common.Hash) string {
	var allowedSequencers []string
	for _, val := range valSet.Validators {
		if val.EndBatch > 0 || !val.IsCurrentValidator(currentBatch) || strings.EqualFold(currentSequencer, val.Signer.EthAddress().Hex()) {
			continue
		}

		allowedSequencers = append(allowedSequencers, val.Signer.EthAddress().Hex())
	}
	sort.Strings(allowedSequencers)

	// make random sequencer
	rand.Seed(int64(currentL2Height + seed.Big().Uint64()))
	randomIndex := rand.Intn(len(allowedSequencers))
	newSequencer := allowedSequencers[randomIndex]

	Logger.Debug("CalcMetisSequencer", "allowedSequencersLen", len(allowedSequencers), "randomIndex", randomIndex, "newSequencer", newSequencer)
	return newSequencer
}

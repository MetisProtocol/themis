package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/spf13/viper"
	tenderCommon "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/metis-seq/themis/bridge/setu/util"
	"github.com/metis-seq/themis/bridge/setu/util/sqlite"
)

const (
	metisEthServiceStr = "metis-eth-service"
)

type MetisEthService struct {
	// Base service
	tenderCommon.BaseService

	listenAddr string

	rpcServer       *rpc.Server
	ethereumService *EthereumService
}

type EthereumService struct {
	logger log.Logger
	// sql client
	sqlClient *sqlite.SqliteDB

	// cache
	cache *expirable.LRU[string, struct{}]
}

func NewMetisEthService(listenAddr string, cdc *codec.Codec) *MetisEthService {
	var logger = util.Logger().With("module", metisEthServiceStr)
	// creating metis eth object
	metisEthService := &MetisEthService{
		listenAddr: listenAddr,
	}
	metisEthService.BaseService = *tenderCommon.NewBaseService(logger, metisEthServiceStr, metisEthService)
	metisEthService.ethereumService = &EthereumService{
		logger:    logger,
		sqlClient: sqlite.GetBridgeSqlDBInstance(viper.GetString(util.BridgeSqliteDBFlag)),
	}

	cache := expirable.NewLRU[string, struct{}](10000, nil, 60*time.Second)
	metisEthService.ethereumService.cache = cache

	return metisEthService
}

// SendTransaction accept eth_sendTransaction RPC
func (es *EthereumService) SendTransaction(ctx context.Context, tx *types.Transaction) (common.Hash, error) {
	// mq producer
	err := es.saveTxToDb(tx)
	if err != nil {
		return common.Hash{}, err
	}
	es.logger.Debug("ethrpc SendTransaction received an tx", "hash", tx.Hash().Hex())
	return tx.Hash(), nil
}

// SendRawTransaction accept eth_sendRawTransaction RPC
func (es *EthereumService) SendRawTransaction(ctx context.Context, encodedTx hexutil.Bytes) (common.Hash, error) {
	tx := new(types.Transaction)
	if err := rlp.DecodeBytes(encodedTx, tx); err != nil {
		es.logger.Error("EthRPC | SendRawTransaction", "Error", err)
		return common.Hash{}, err
	}
	// mq producer
	err := es.saveTxToDb(tx)
	if err != nil {
		return common.Hash{}, err
	}
	es.logger.Debug("ethrpc SendRawTransaction received an tx", "hash", tx.Hash().Hex())
	return tx.Hash(), nil
}

func (es *EthereumService) saveTxToDb(tx *types.Transaction) error {
	if !util.RecoverSpanFinished {
		es.logger.Info("ethrpc SendRawTransaction received tx before recover span finished", "hash", tx.Hash().Hex())
		return errors.New("wait span recover")
	}

	// cache filter
	_, exist := es.cache.Get(tx.Hash().Hex())
	if exist {
		es.logger.Info("ethrpc SendRawTransaction received repeated tx", "hash", tx.Hash().Hex())
		return nil
	}
	es.cache.Add(tx.Hash().Hex(), struct{}{})

	// set health value
	HealthValue.RPC.TxHash = tx.Hash().Hex()
	HealthValue.RPC.Timestamp = time.Now().Unix()

	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return err
	}

	// send chan
	txData := hexutil.Encode(data)
	util.MetisTxCacheChan <- &sqlite.MetisTxCache{
		TxData: txData,
	}
	es.logger.Info("ethrpc SendRawTransaction saveTxToDb tx", "hash", tx.Hash().Hex())
	return nil
}

func (s *MetisEthService) OnStart() error {
	if err := s.BaseService.OnStart(); err != nil {
		s.Logger.Error("OnStart | OnStart", "Error", err)
	} // Always call the overridden method.

	// start eth rpc listeners
	s.rpcServer = rpc.NewServer()
	if err := s.rpcServer.RegisterName("eth", s.ethereumService); err != nil {
		s.Logger.Error("OnStart | Rpc RegisterName", "Error", err)
		return err
	}

	http.Handle("/", s.rpcServer)
	http.HandleFunc("/health", health)
	s.Logger.Info("starting Ethereum RPC server on :8646")
	err := http.ListenAndServe(s.listenAddr, nil)
	if err != nil {
		s.Logger.Error("OnStart | Starting Ethereum RPC", "Error", err)
		return err
	}
	s.Logger.Info("all listeners Started")
	return nil
}

// OnStop stops all necessary go routines
func (s *MetisEthService) OnStop() {
	s.BaseService.OnStop() // Always call the overridden method.
	// stop eth rpc listeners
	s.rpcServer.Stop()
	s.Logger.Info("all listeners stopped")
}

type healthResponse struct {
	IsCurrentSequencer int `json:"isCurrentSequencer"`
	Writer             struct {
		TxHash    string `json:"txHash"`
		Timestamp int64  `json:"timestamp"`
		Ret       string `json:"ret"`
	} `json:"writer"`
	L2 struct {
		BlockNumber uint64 `json:"blockNumber"`
		Timestamp   int64  `json:"timestamp"`
	} `json:"l2"`
	RPC struct {
		TxHash    string `json:"txHash"`
		Timestamp int64  `json:"timestamp"`
	} `json:"rpc"`
	Mpc struct {
		IsMpcProposer int    `json:"isMpcProposer"`
		SignID        string `json:"signId"`
		SignSuccess   int    `json:"signSuccess"`
		Timestamp     int64  `json:"timestamp"`
	} `json:"mpc"`
}

var HealthValue healthResponse

func health(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(HealthValue)
	w.Write(data)
}

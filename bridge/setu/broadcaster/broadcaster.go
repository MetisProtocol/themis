package broadcaster

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/cosmos/cosmos-sdk/client"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	metis "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"

	authTypes "github.com/metis-seq/themis/auth/types"
	"github.com/metis-seq/themis/bridge/setu/util"
	"github.com/metis-seq/themis/helper"

	"github.com/tendermint/tendermint/libs/log"

	hmTypes "github.com/metis-seq/themis/types"
)

// TxBroadcaster uses to broadcast transaction to each chain
type TxBroadcaster struct {
	logger log.Logger

	CliCtx cliContext.CLIContext

	themisMutex sync.Mutex
	metisMutex  sync.Mutex

	lastSeqNo uint64
	accNum    uint64
}

// NewTxBroadcaster creates new broadcaster
func NewTxBroadcaster(cdc *codec.Codec) *TxBroadcaster {
	cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	cliCtx.BroadcastMode = client.BroadcastSync
	cliCtx.TrustNode = true

	// current address
	address := hmTypes.BytesToThemisAddress(helper.GetAddress())

	account, err := util.GetAccount(cliCtx, address)
	if err != nil {
		panic("Error connecting to rest-server, please start server before bridge.")
	}

	return &TxBroadcaster{
		logger:    util.Logger().With("module", "txBroadcaster"),
		CliCtx:    cliCtx,
		lastSeqNo: account.GetSequence(),
		accNum:    account.GetAccountNumber(),
	}
}

// BroadcastToThemis broadcast to themis
func (tb *TxBroadcaster) BroadcastToThemis(msg sdk.Msg, event interface{}) error {
	tb.themisMutex.Lock()
	defer tb.themisMutex.Unlock()

	// tx encoder
	txEncoder := helper.GetTxEncoder(tb.CliCtx.Codec)
	// chain id
	chainID := helper.GetGenesisDoc().ChainID

	if tb.lastSeqNo == 0 {
		// get account number and sequence
		tb.resetSeqNo() //reset seqNo
	}

	txBldr := authTypes.NewTxBuilderFromCLI().
		WithTxEncoder(txEncoder).
		WithAccountNumber(tb.accNum).
		WithSequence(tb.lastSeqNo).
		WithChainID(chainID).
		// WithGas(helper.GetConfig().MainchainGasLimit).
		WithGas(helper.MainchainBuildGasLimit).
		WithGasAdjustment(float64(1.5))

	txResponse, err := helper.BuildAndBroadcastMsgs(tb.CliCtx, txBldr, []sdk.Msg{msg})
	if err != nil {
		tb.logger.Error("Error while broadcasting the themis transaction", "error", err)
		tb.resetSeqNo() //reset seqNo
		return err
	}

	txHash := txResponse.TxHash

	tb.logger.Info("Tx sent on themis", "txHash", txHash, "accSeq", tb.lastSeqNo, "accNum", tb.accNum)
	if txResponse.Code != 0 {
		tb.logger.Error("Tx failed on themis", "txHash", txHash, "accSeq", tb.lastSeqNo, "accNum", tb.accNum, "txResponseErr", txResponse.RawLog)
		tb.resetSeqNo() //reset seqNo
		return errors.New(txResponse.RawLog)
	}
	// increment account sequence
	tb.lastSeqNo += 1

	return nil
}

func (tb *TxBroadcaster) resetSeqNo() error {
	// current address
	address := hmTypes.BytesToThemisAddress(helper.GetAddress())
	// fetch from APIs
	account, errAcc := util.GetAccount(tb.CliCtx, address)
	if errAcc != nil {
		tb.logger.Error("Error fetching account from rest-api", "url", helper.GetThemisServerEndpoint(fmt.Sprintf(util.AccountDetailsURL, helper.GetAddress())))
		return errAcc
	}
	// update seqNo for safety
	tb.lastSeqNo = account.GetSequence()
	return nil
}

// BroadcastToMetis broadcast to metis
func (tb *TxBroadcaster) BroadcastToMetis(msg metis.CallMsg) error {
	tb.metisMutex.Lock()
	defer tb.metisMutex.Unlock()

	// get metis client
	metisClient := helper.GetMetisClient()

	// get auth
	auth, err := helper.GenerateAuthObj(metisClient, *msg.To, msg.Data)

	if err != nil {
		tb.logger.Error("Error generating auth object", "error", err)
		return err
	}

	// Create the transaction, sign it and schedule it for execution
	rawTx := types.NewTransaction(auth.Nonce.Uint64(), *msg.To, msg.Value, auth.GasLimit, auth.GasPrice, msg.Data)

	// signer
	signedTx, err := auth.Signer(auth.From, rawTx)
	if err != nil {
		tb.logger.Error("Error signing the transaction", "error", err)
		return err
	}
	tb.logger.Info("Sending transaction to metis", "txHash", signedTx.Hash())

	// create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), helper.GetConfig().MetisRPCTimeout)
	defer cancel()

	// broadcast transaction
	err = metisClient.SendTransaction(ctx, signedTx)
	if err != nil {
		tb.logger.Error("BroadcastToMetis", "error", err)

		if strings.Contains(err.Error(), "Timeout") || strings.Contains(err.Error(), "timeout") {
			tb.logger.Error("Error while broadcasting the transaction to metischain", "error", "timeout")
			return errors.New("timeout")
		}
	}

	return nil
}

// BroadcastToMetis broadcast to metis
func (tb *TxBroadcaster) BroadcastRawTxToMetis(signedTx *types.Transaction) error {
	tb.metisMutex.Lock()
	defer tb.metisMutex.Unlock()

	// get metis client
	metisClient := helper.GetMetisClient()

	// create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), helper.GetConfig().MetisRPCTimeout)
	defer cancel()

	// broadcast transaction
	err := metisClient.SendTransaction(ctx, signedTx)
	if err != nil {
		tb.logger.Error("Error broadcastToMetis while broadcasting rpcTx to metis", "txHash", signedTx.Hash().Hex(), "error", err)
		if strings.Contains(err.Error(), "dial tcp") ||
			strings.Contains(err.Error(), "sequencer") {
			return err
		}
	}
	return nil
}

// BroadcastToRootchain broadcast to rootchain
func (tb *TxBroadcaster) BroadcastToRootchain() {}

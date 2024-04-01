package processor

import (
	"context"
	"encoding/hex"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/metis-seq/themis/contracts/sequencerset"
	"github.com/metis-seq/themis/helper"
	"github.com/metis-seq/themis/helper/tss"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
func TestCommitSpan(t *testing.T) {
	genesisContractAddress := "0x5FC8d32690cc91D4c39d9d3abcBD16989F875707"

	metisRPCClient, err := rpc.Dial("http://127.0.0.1:8545")
	if err != nil {
		panic(err)
	}

	metisClient := ethclient.NewClient(metisRPCClient)
	metisContract, err := sequencerset.NewSequencerset(common.HexToAddress(genesisContractAddress), metisClient)
	if err != nil {
		panic(err)
	}

	chainID, err := metisClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	t.Logf("chainID: %v\n", chainID.Int64())

	mpcPri := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	mpcAddress := "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	privateKey, err := crypto.HexToECDSA(mpcPri)
	if err != nil {
		panic(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		panic(err)
	}

	nonce, err := metisClient.PendingNonceAt(context.Background(), common.HexToAddress(mpcAddress))
	if err != nil {
		panic(err)
	}
	t.Logf("nonce:%v", nonce)

	tx, err := metisContract.CommitEpoch(
		auth,
		big.NewInt(1),
		big.NewInt(600),
		big.NewInt(799),
		common.HexToAddress("0x3eb630c3c267395fee216b603a02061330d39642"),
	)
	if err != nil {
		panic(err)
	}
	t.Logf("tx:%v", tx.Hash().Hex())

	_, err = bind.WaitMined(context.Background(), metisClient, tx)
	if err != nil {
		panic(err)
	}
	t.Log("commitEpoch success")

	// query epoch
	epochId, err := metisContract.CurrentEpochNumber(nil)
	if err != nil {
		panic(err)
	}
	t.Logf("current epoch:%v", epochId)
}
*/

func TestCommitSpanWithOfflineSign(t *testing.T) {
	genesisContractAddress := "0x5FC8d32690cc91D4c39d9d3abcBD16989F875707"

	seqSetAbi, err := abi.JSON(strings.NewReader(sequencerset.SequencersetABI))
	if err != nil {
		panic(err)
	}

	metisRPCClient, err := rpc.Dial("http://127.0.0.1:8545")
	if err != nil {
		panic(err)
	}

	metisClient := ethclient.NewClient(metisRPCClient)
	metisContract, err := sequencerset.NewSequencerset(common.HexToAddress(genesisContractAddress), metisClient)
	if err != nil {
		panic(err)
	}
	epochId, err := metisContract.CurrentEpochNumber(nil)
	if err != nil {
		panic(err)
	}
	t.Logf("current epoch:%v", epochId)

	chainID, err := metisClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	t.Logf("chainID: %v\n", chainID.Int64())

	mpcPri := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	mpcAddress := "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	privateKey, err := crypto.HexToECDSA(mpcPri)
	if err != nil {
		panic(err)
	}
	// auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	// if err != nil {
	// 	panic(err)
	// }

	nonce, err := metisClient.PendingNonceAt(context.Background(), common.HexToAddress(mpcAddress))
	if err != nil {
		panic(err)
	}
	t.Logf("nonce:%v", nonce)

	gasPrice, err := metisClient.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}
	t.Logf("gasPrice:%v", gasPrice)

	input, err := seqSetAbi.Pack("commitEpoch",
		big.NewInt(1),
		big.NewInt(600),
		big.NewInt(799),
		common.HexToAddress("0x3eb630c3c267395fee216b603a02061330d39642"),
	)
	if err != nil {
		panic(err)
	}

	tx := types.NewTransaction(nonce, common.HexToAddress(genesisContractAddress), big.NewInt(0), 2000000, gasPrice, input)
	// calc sig hash
	txSigner := types.NewEIP155Signer(chainID)

	// generate sig msg
	signMsg := txSigner.Hash(tx).Bytes()

	// generate signature
	sig, err := crypto.Sign(signMsg[:], privateKey)
	if err != nil {
		panic(err)
	}

	// build signed tx
	signedTx, err := tx.WithSignature(txSigner, sig)
	if err != nil {
		panic(err)
	}
	t.Logf("tx:%v", signedTx.Hash().Hex())

	err = metisClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		panic(err)
	}

	_, err = bind.WaitMined(context.Background(), metisClient, signedTx)
	if err != nil {
		panic(err)
	}
	t.Log("commitEpoch success")

	// query epoch
	epochId, err = metisContract.CurrentEpochNumber(nil)
	if err != nil {
		panic(err)
	}
	t.Logf("current epoch:%v", epochId)

	epochStruct, err := metisContract.Epochs(nil, epochId)
	if err != nil {
		panic(err)
	}
	t.Logf("current epoch info:%v", epochStruct)
}

func TestParseSignedTx(t *testing.T) {
	txData := "f8e78001834c4b40945fc8d32690cc91d4c39d9d3abcbd16989f87570780b8844fb71bdd00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000258000000000000000000000000000000000000000000000000000000000000031f0000000000000000000000003eb630c3c267395fee216b603a02061330d396428208bfa026c8d20cb5fed9e307c0d11f901470f1f656cdc16a7ec4957f4e6286de51f9f8a02ed9b846d1a908a8f82c8c53483fb3188cda0da46c9d3d78f5791649fad21072"
	tx := new(types.Transaction)

	txBytes, _ := hex.DecodeString(txData)
	tx.UnmarshalBinary(txBytes)

	t.Logf("signed tx chainID:%v", tx.ChainId())
	t.Logf("signed tx data:%x", tx.Data())
	t.Logf("signed tx nonce:%v", tx.Nonce())
	t.Logf("signed tx to address :%v", tx.To().Hex())

	r, s, v := tx.RawSignatureValues()
	t.Logf("signed tx to sig :%v,%v,%v", r, s, v)
}

func TestDelayTime(t *testing.T) {
	delayTicker := time.NewTicker(10 * time.Second)
	chainBlockRefreshTicker := time.NewTicker(1 * time.Second)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
	DELAY:
		for {
			select {
			case <-delayTicker.C:
				wg.Done()
				break DELAY
			case <-chainBlockRefreshTicker.C:
				t.Log("refresh ticker")
			}
		}
	}()
	wg.Wait()

	t.Log("success")
}

func TestReCommitSpanWithOfflineSign(t *testing.T) {
	genesisContractAddress := "0xf83062f613527534a65ece52d93bcacfc4a2319b"

	seqSetAbi, err := abi.JSON(strings.NewReader(sequencerset.SequencersetABI))
	if err != nil {
		panic(err)
	}

	metisRPCClient, err := rpc.Dial("http://127.0.0.1:8545")
	if err != nil {
		panic(err)
	}

	metisClient := ethclient.NewClient(metisRPCClient)
	metisContract, err := sequencerset.NewSequencerset(common.HexToAddress(genesisContractAddress), metisClient)
	if err != nil {
		panic(err)
	}
	epochId, err := metisContract.CurrentEpochNumber(nil)
	if err != nil {
		panic(err)
	}
	t.Logf("current epoch:%v", epochId)

	chainID, err := metisClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	t.Logf("chainID: %v\n", chainID.Int64())

	// mpcPri := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	mpcAddress := "0x4835bd266b19887d56972474ec22fa769fd2a77b"
	// privateKey, err := crypto.HexToECDSA(mpcPri)
	// if err != nil {
	// 	panic(err)
	// }

	nonce, err := metisClient.PendingNonceAt(context.Background(), common.HexToAddress(mpcAddress))
	if err != nil {
		panic(err)
	}
	t.Logf("nonce:%v", nonce)

	gasPrice, err := metisClient.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}
	t.Logf("gasPrice:%v", gasPrice)

	input, err := seqSetAbi.Pack("recommitEpoch",
		big.NewInt(1),
		big.NewInt(2),
		big.NewInt(1100800),
		big.NewInt(1101399),
		common.HexToAddress("0x1267397fb5bf6f6dcc3d18d673616d512dbcd8f0"),
	)
	if err != nil {
		panic(err)
	}

	tx := types.NewTransaction(nonce, common.HexToAddress(genesisContractAddress), big.NewInt(0), 2000000, gasPrice, input)
	// calc sig hash
	txSigner := types.NewEIP155Signer(chainID)

	// generate sig msg
	signMsg := txSigner.Hash(tx).Bytes()
	t.Logf("signMsg:%x", signMsg)

	// mpc sign
	mpcClientConn, err := grpc.Dial(
		"3.213.188.165:8081",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	mpcClient := tss.NewTssServiceClient(mpcClientConn)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	request := &tss.KeySignRequest{
		KeyId:   "479f01f8-fca6-4ac5-b0e9-ff244821a43b",
		SignMsg: signMsg,
		SignId:  "69920fde-57e3-4a64-9265-20c109d3bcc5",
	}
	t.Logf("MpcSign request data:%v", request.String())

	keySignResp, err := mpcClient.KeySign(ctx, request)
	if err != nil {
		panic(err)
	}

	if keySignResp == nil {
		panic("nil sign result")
	}

	sig := helper.ConvertSignature(keySignResp.SignatureR, keySignResp.SignatureS, keySignResp.SignatureV)
	t.Logf("sig:%x", sig)

	// build signed tx
	signedTx, err := tx.WithSignature(txSigner, sig)
	if err != nil {
		panic(err)
	}
	t.Logf("tx:%v", signedTx.Hash().Hex())

	err = metisClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		panic(err)
	}

	_, err = bind.WaitMined(context.Background(), metisClient, signedTx)
	if err != nil {
		panic(err)
	}
	t.Log("commitEpoch success")

	// query epoch
	epochId, err = metisContract.CurrentEpochNumber(nil)
	if err != nil {
		panic(err)
	}
	t.Logf("current epoch:%v", epochId)

	epochStruct, err := metisContract.Epochs(nil, epochId)
	if err != nil {
		panic(err)
	}
	t.Logf("current epoch info:%v", epochStruct)
}

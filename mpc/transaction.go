package mpc

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	hmTypes "github.com/metis-seq/themis/types"
)

func GetCheckTxSignMsgData(ctx sdk.Context, k Keeper, msgType hmTypes.SignType, signData []byte) (types.Signer, []byte, []byte, error) {
	tx := new(types.Transaction)
	err := tx.UnmarshalBinary(signData)
	if err != nil {
		k.Logger(ctx).Error("CheckTxSignMsg tx unmarlshal failed", "err", err)
		return nil, nil, nil, err
	}

	chainParams := k.cm.GetChainParams(ctx)
	k.Logger(ctx).Info("get chain GetParams", "chainParams", chainParams)

	var chainID uint64
	switch msgType {
	case hmTypes.BatchSubmit, hmTypes.L1UpdateMpcAddress, hmTypes.BatchReward:
		chainID, _ = strconv.ParseUint(chainParams.ChainParams.MainChainID, 10, 64)
	case hmTypes.CommitEpochToMetis, hmTypes.ReCommitEpochToMetis, hmTypes.L2UpdateMpcAddress:
		chainID, _ = strconv.ParseUint(chainParams.ChainParams.MetisChainID, 10, 64)
	default:
		return nil, nil, nil, fmt.Errorf("CheckTxSignMsg unsupport sign type")
	}

	k.Logger(ctx).Info("GetCheckTxSignMsgData chain_id", "chain_id", chainID)
	txSigner := types.NewLondonSigner(big.NewInt(int64(chainID)))
	return txSigner, txSigner.Hash(tx).Bytes(), tx.Data(), nil
}

func DecodeTransactionInputData(contractABI *abi.ABI, data []byte, result interface{}) error {
	if len(data) < 4 {
		return errors.New("invalid input data")
	}

	methodSigData := data[:4]
	method, err := contractABI.MethodById(methodSigData)
	if err != nil {
		return err
	}

	inputsSigData := data[4:]
	inputs, err := method.Inputs.Unpack(inputsSigData)
	if err != nil {
		return err
	}

	return method.Inputs.Copy(result, inputs)
}

func CheckBatchRewardTxLogic(ctx sdk.Context, k Keeper, txData []byte, txMsg []byte) error {
	type BatchRewardInput struct {
		BatchId    *big.Int         `json:"_batchId"`
		StartEpoch *big.Int         `json:"_startEpoch"`
		EndEpoch   *big.Int         `json:"_endEpoch"`
		Seqs       []common.Address `json:"_seqs"`
		Blocks     []*big.Int       `json:"_blocks"`
	}

	var batchRewardInput BatchRewardInput
	err := DecodeTransactionInputData(&k.contractCaller.StakeManagerABI, txData, &batchRewardInput)
	if err != nil {
		k.Logger(ctx).Error("invalid mpc sign data",
			"batchId", batchRewardInput.BatchId, "signData", hex.EncodeToString(txData), "error", err)
		return err
	}

	address := k.cm.GetChainParams(ctx).ChainParams.ValidatorSetAddress.EthAddress()
	seqset, err := k.contractCaller.GetSequencerSetInstance(address)
	if err != nil {
		return err
	}

	k.Logger(ctx).Info("CheckBatchRewardTxLogic", "seqset", address, "batchId", batchRewardInput.BatchId,
		"txData", hex.EncodeToString(txData), "txMsg", hex.EncodeToString(txMsg))

	// Statistics span information in db
	sequencerCount := make(map[ethCommon.Address]uint64)

	for spanId := batchRewardInput.StartEpoch.Uint64(); spanId <= batchRewardInput.EndEpoch.Uint64(); spanId++ {
		info, err := seqset.Epochs(&bind.CallOpts{Context: ctx.Context()}, new(big.Int).SetUint64(spanId))
		if err != nil {
			k.Logger(ctx).Error("unable to get the epoch", "spanId", spanId, "error", err)
			return err
		}

		if info.Signer == (ethCommon.Address{}) || info.StartBlock.BitLen() == 0 || info.EndBlock.BitLen() == 0 {
			k.Logger(ctx).Error("invalid spanId, span not exits", "spanId", spanId)
			return fmt.Errorf("span not found")
		}

		finishedBlock := info.EndBlock.Uint64() - info.StartBlock.Uint64() + 1
		sequencerCount[info.Signer] += finishedBlock
	}

	if len(sequencerCount) != len(batchRewardInput.Seqs) {
		k.Logger(ctx).Error("invalid sequencer length", "batchId", batchRewardInput.BatchId,
			"sequencerLen", len(batchRewardInput.Seqs), "dbSequencerLen", len(sequencerCount))
		return errors.New("invalid sequencer length")
	}

	if len(sequencerCount) != len(batchRewardInput.Blocks) {
		k.Logger(ctx).Error("invalid finish block length", "batchId", batchRewardInput.BatchId,
			"finishedBlocksLen", len(batchRewardInput.Blocks), "dbFinishedBlocksLen", len(sequencerCount))
		return errors.New("invalid finish block length")
	}

	// check finished blocks
	for i := 0; i < len(batchRewardInput.Seqs); i++ {
		signer := batchRewardInput.Seqs[i]
		count := batchRewardInput.Blocks[i].Uint64()

		dbCount := sequencerCount[signer]
		if dbCount != count {
			k.Logger(ctx).Error("invalid batch reward block count", "batchId", batchRewardInput.BatchId,
				"signer", signer, "dbFinishedBlocks", dbCount, "txFinishedBlocks", count)
			return errors.New("invalid finish reward block count")
		}
	}
	return nil
}

func CheckBatchSubmitTxLogic(ctx sdk.Context, k Keeper, txData []byte) error {
	return nil
}

type CommitEpochToMetisInput struct {
	EpochId    *big.Int
	StartBlock *big.Int
	EndBlock   *big.Int
	Signer     common.Address
}

func CheckCommitEpochToMetisTxLogic(ctx sdk.Context, k Keeper, txData []byte) error {
	return nil
}

type ReCommitEpochToMetisInput struct {
	OldEpochId *big.Int
	NewEpochId *big.Int
	StartBlock *big.Int
	EndBlock   *big.Int
	NewSigner  common.Address
}

func CheckReCommitEpochToMetisTxLogic(ctx sdk.Context, k Keeper, txData []byte) error {
	return nil
}

type L1UpdateMpcAddressInput struct {
	NewMpcAddress common.Address
}

func CheckL1UpdateMpcAddressTxLogic(ctx sdk.Context, k Keeper, txData []byte) error {
	return nil
}

type L2UpdateMpcAddressInput struct {
	NewMpcAddress common.Address
}

func CheckL2UpdateMpcAddressTxLogic(ctx sdk.Context, k Keeper, txData []byte) error {
	return nil
}

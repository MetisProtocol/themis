package processor

import (
	"math/big"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/metis-seq/themis/bridge/setu/util"
	"github.com/metis-seq/themis/helper"
)

// getCurrentChildBlock gets the current child block
func getCurrentChildBlock(cc helper.ContractCaller) (uint64, error) {
	childBlock, err := cc.GetMetisChainBlock(nil)
	if err != nil {
		return 0, err
	}

	return childBlock.Number.Uint64(), nil
}

// getSequencerEpochLength gets the current sequencer
func getSequencerEpochLength(cc helper.ContractCaller, cliCtx context.CLIContext) (int64, error) {
	params, err := util.GetChainmanagerParams(cliCtx)
	if err != nil {
		return 0, err
	}
	sequencerSet, err := cc.GetSequencerSetInstance(params.ChainParams.ValidatorSetAddress.EthAddress())
	if err != nil {
		return 0, err
	}

	epochLength, err := sequencerSet.EpochLength(nil)
	if err != nil {
		return 0, err
	}

	return epochLength.Int64(), nil
}

// getCurrentChildSequencer gets the current sequencer
func getCurrentChildSequencer(cc helper.ContractCaller, cliCtx context.CLIContext, height uint64) (string, int64, error) {
	params, err := util.GetChainmanagerParams(cliCtx)
	if err != nil {
		return "", 0, err
	}
	sequencerSet, err := cc.GetSequencerSetInstance(params.ChainParams.ValidatorSetAddress.EthAddress())
	if err != nil {
		return "", 0, err
	}

	currentEpochNumber, err := sequencerSet.GetEpochByBlock(nil, big.NewInt(int64(height)))
	if err != nil {
		return "", 0, err
	}

	currentEpochInfo, err := sequencerSet.Epochs(nil, currentEpochNumber)
	if err != nil {
		return "", 0, err
	}

	return currentEpochInfo.Signer.Hex(), currentEpochNumber.Int64(), nil
}

func getChildLatestEpoch(cc helper.ContractCaller, cliCtx context.CLIContext) (int64, int64, int64, error) {
	params, err := util.GetChainmanagerParams(cliCtx)
	if err != nil {
		return 0, 0, 0, err
	}
	sequencerSet, err := cc.GetSequencerSetInstance(params.ChainParams.ValidatorSetAddress.EthAddress())
	if err != nil {
		return 0, 0, 0, err
	}

	currentEpochNumber, err := sequencerSet.CurrentEpochNumber(nil)
	if err != nil {
		return 0, 0, 0, err
	}

	currentEpochInfo, err := sequencerSet.Epochs(nil, currentEpochNumber)
	if err != nil {
		return 0, 0, 0, err
	}

	return currentEpochNumber.Int64(), currentEpochInfo.StartBlock.Int64(), currentEpochInfo.EndBlock.Int64(), nil
}

type epochInfo struct {
	ID         int64
	Signer     string
	StartBlock int64
	EndBlock   int64
}

func getChildEpochInfo(cc helper.ContractCaller, cliCtx context.CLIContext, epochNumber int64) (*epochInfo, error) {
	params, err := util.GetChainmanagerParams(cliCtx)
	if err != nil {
		return nil, err
	}
	sequencerSet, err := cc.GetSequencerSetInstance(params.ChainParams.ValidatorSetAddress.EthAddress())
	if err != nil {
		return nil, err
	}

	epoch, err := sequencerSet.Epochs(nil, big.NewInt(epochNumber))
	if err != nil {
		return nil, err
	}

	return &epochInfo{
		ID:         epoch.Number.Int64(),
		Signer:     epoch.Signer.Hex(),
		StartBlock: epoch.StartBlock.Int64(),
		EndBlock:   epoch.EndBlock.Int64(),
	}, nil
}

// getCurrentMainBatch gets the current epoch
func getCurrentMainBatch(cc helper.ContractCaller, cliCtx context.CLIContext) (uint64, error) {
	params, err := util.GetChainmanagerParams(cliCtx)
	if err != nil {
		return 0, err
	}
	stakeManager, err := cc.GetStakeManagerInstance(params.ChainParams.StakingManagerAddress.EthAddress())
	if err != nil {
		return 0, err
	}

	currentBatch, err := stakeManager.CurrentBatch(nil)
	if err != nil {
		return 0, err
	}

	return currentBatch.Uint64(), nil
}

package staking_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/metis-seq/themis/app"
	errs "github.com/metis-seq/themis/common"
	"github.com/metis-seq/themis/contracts/stakinginfo"
	"github.com/metis-seq/themis/helper"
	"github.com/metis-seq/themis/helper/mocks"
	"github.com/metis-seq/themis/staking"
	stakingSim "github.com/metis-seq/themis/staking/simulation"
	"github.com/metis-seq/themis/staking/types"

	hmTypes "github.com/metis-seq/themis/types"
	"github.com/metis-seq/themis/types/simulation"
)

type HandlerTestSuite struct {
	suite.Suite

	app    *app.ThemisApp
	ctx    sdk.Context
	cliCtx context.CLIContext

	handler        sdk.Handler
	contractCaller mocks.IContractCaller
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.app, suite.ctx, suite.cliCtx = createTestApp(false)
	suite.contractCaller = mocks.IContractCaller{}
	suite.handler = staking.NewHandler(suite.app.StakingKeeper, &suite.contractCaller)
}

func TestHandlerTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestHandleMsgValidatorJoin() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	txHash := hmTypes.HexToThemisHash("123")
	index := simulation.RandIntBetween(r1, 0, 100)
	logIndex := uint64(index)
	validatorId := uint64(1)
	amount, _ := big.NewInt(0).SetString("10000000000000000000", 10)

	privKey1 := secp256k1.GenPrivKey()
	pubkey := hmTypes.NewPubKey(privKey1.PubKey().Bytes())
	address := pubkey.Address()

	chainParams := app.ChainKeeper.GetParams(ctx)

	txreceipt := &ethTypes.Receipt{
		BlockNumber: big.NewInt(10),
	}

	msgValJoin := types.NewMsgValidatorJoin(
		hmTypes.BytesToThemisAddress(address.Bytes()),
		validatorId,
		uint64(1),
		sdk.NewInt(1000000000000000000),
		pubkey,
		txHash,
		logIndex,
		0,
		1,
	)

	stakinginfoStaked := &stakinginfo.StakinginfoLocked{
		Signer:          address,
		SequencerId:     new(big.Int).SetUint64(validatorId),
		ActivationBatch: big.NewInt(1),
		Amount:          amount,
		Total:           big.NewInt(10),
		SignerPubkey:    pubkey.Bytes()[1:],
	}

	suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txreceipt, nil)

	suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress.EthAddress(), txreceipt, msgValJoin.LogIndex).Return(stakinginfoStaked, nil)

	result := suite.handler(ctx, msgValJoin)
	require.True(t, result.IsOK(), "expected validator join to be ok, got %v", result)

	actualResult, ok := app.StakingKeeper.GetValidatorFromValID(ctx, hmTypes.ValidatorID(validatorId))
	require.Equal(t, false, ok, "Should not add validator")
	require.NotNil(t, actualResult, "got %v", actualResult)
}

func (suite *HandlerTestSuite) TestHandleMsgValidatorUpdate() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := suite.app.StakingKeeper
	// pass 0 as time alive to generate non de-activated validators
	staking.LoadValidatorSet(t, 4, keeper, ctx, false, 0)
	oldValSet := keeper.GetValidatorSet(ctx)

	// vals := oldValSet.(*Validators)
	oldSigner := oldValSet.Validators[0]
	newSigner := stakingSim.GenRandomVal(1, 0, 10, 10, false, 1)
	newSigner[0].ID = oldSigner.ID
	newSigner[0].VotingPower = oldSigner.VotingPower

	t.Log("To be Updated ===>", "Validator", newSigner[0].String())

	chainParams := app.ChainKeeper.GetParams(ctx)

	// gen msg
	msgTxHash := hmTypes.HexToThemisHash("123")
	msg := types.NewMsgSignerUpdate(newSigner[0].Signer, uint64(newSigner[0].ID), newSigner[0].PubKey, msgTxHash, 0, 0, 1)

	txreceipt := &ethTypes.Receipt{BlockNumber: big.NewInt(10)}
	suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txreceipt, nil)

	signerUpdateEvent := &stakinginfo.StakinginfoSignerChange{
		SequencerId:  new(big.Int).SetUint64(oldSigner.ID.Uint64()),
		OldSigner:    oldSigner.Signer.EthAddress(),
		NewSigner:    newSigner[0].Signer.EthAddress(),
		SignerPubkey: newSigner[0].PubKey.Bytes()[1:],
	}
	suite.contractCaller.On("DecodeSignerUpdateEvent", chainParams.ChainParams.StakingInfoAddress.EthAddress(), txreceipt, uint64(0)).Return(signerUpdateEvent, nil)

	result := suite.handler(ctx, msg)

	require.True(t, result.IsOK(), "expected validator update to be ok, got %v", result)

	newValidators := keeper.GetCurrentValidators(ctx)
	require.Equal(t, len(oldValSet.Validators), len(newValidators), "Number of current validators should be equal")

	setUpdates := helper.GetUpdatedValidators(&oldValSet, keeper.GetAllValidators(ctx), 5)

	err := oldValSet.UpdateWithChangeSet(setUpdates)
	require.NoError(t, err)

	_ = keeper.UpdateValidatorSetInStore(ctx, oldValSet)

	ValFrmID, ok := keeper.GetValidatorFromValID(ctx, oldSigner.ID)
	require.True(t, ok, "signer should be found, got %v", ok)
	require.NotEqual(t, oldSigner.Signer.Bytes(), newSigner[0].Signer.Bytes(), "Should not update state")
	require.Equal(t, ValFrmID.VotingPower, oldSigner.VotingPower, "VotingPower of new signer %v should be equal to old signer %v", ValFrmID.VotingPower, oldSigner.VotingPower)

	removedVal, err := keeper.GetValidatorInfo(ctx, oldSigner.Signer.Bytes())
	require.Empty(t, err)
	require.NotEqual(t, removedVal.VotingPower, int64(0), "should not update state")
}

func (suite *HandlerTestSuite) TestHandleMsgValidatorExit() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.StakingKeeper
	// pass 0 as time alive to generate non de-activated validators
	staking.LoadValidatorSet(t, 4, keeper, ctx, false, 0)
	validators := keeper.GetCurrentValidators(ctx)
	msgTxHash := hmTypes.HexToThemisHash("123")
	chainParams := app.ChainKeeper.GetParams(ctx)
	logIndex := uint64(0)

	txreceipt := &ethTypes.Receipt{
		BlockNumber: big.NewInt(10),
	}

	suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txreceipt, nil)

	amount, _ := big.NewInt(0).SetString("10000000000000000000", 10)
	stakinginfoUnstakeInit := &stakinginfo.StakinginfoUnlockInit{
		User:              validators[0].Signer.EthAddress(),
		SequencerId:       big.NewInt(0).SetUint64(validators[0].ID.Uint64()),
		DeactivationBatch: big.NewInt(10),
		Amount:            amount,
	}
	validators[0].EndBatch = 10

	suite.contractCaller.On("DecodeValidatorExitEvent", chainParams.ChainParams.StakingInfoAddress.EthAddress(), txreceipt, logIndex).Return(stakinginfoUnstakeInit, nil)

	msg := types.NewMsgValidatorExit(validators[0].Signer, uint64(validators[0].ID), validators[0].EndBatch, msgTxHash, 0, 0, 1)

	got := suite.handler(ctx, msg)

	require.True(t, got.IsOK(), "expected validator exit to be ok, got %v", got)

	updatedValInfo, err := keeper.GetValidatorInfo(ctx, validators[0].Signer.Bytes())
	// updatedValInfo.EndBatch = 10
	require.Empty(t, err, "Unable to get validator info from val address,ValAddr:%v Error:%v ", validators[0].Signer.String(), err)
	require.NotEqual(t, updatedValInfo.EndBatch, validators[0].EndBatch, "should not update deactivation batch")

	_, found := keeper.GetValidatorFromValID(ctx, validators[0].ID)
	require.True(t, found, "Validator should be present even after deactivation")

	got = suite.handler(ctx, msg)
	require.True(t, got.IsOK(), "should not fail, as state is not updated for validatorExit")
}

func (suite *HandlerTestSuite) TestHandleMsgStakeUpdate() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.StakingKeeper

	// pass 0 as time alive to generate non de-activated validators
	staking.LoadValidatorSet(t, 4, keeper, ctx, false, 0)
	oldValSet := keeper.GetValidatorSet(ctx)
	oldVal := oldValSet.Validators[0]

	t.Log("To be Updated ===>", "Validator", oldVal.String())

	chainParams := app.ChainKeeper.GetParams(ctx)

	msgTxHash := hmTypes.HexToThemisHash("123")
	msg := types.NewMsgStakeUpdate(oldVal.Signer, oldVal.ID.Uint64(), sdk.NewInt(2000000000000000000), msgTxHash, 0, 0, 1)

	txreceipt := &ethTypes.Receipt{BlockNumber: big.NewInt(10)}
	suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txreceipt, nil)

	stakinginfoStakeUpdate := &stakinginfo.StakinginfoLockUpdate{
		SequencerId: new(big.Int).SetUint64(oldVal.ID.Uint64()),
		NewAmount:   new(big.Int).SetInt64(2000000000000000000),
	}

	suite.contractCaller.On("DecodeValidatorStakeUpdateEvent", chainParams.ChainParams.StakingInfoAddress.EthAddress(), txreceipt, uint64(0)).Return(stakinginfoStakeUpdate, nil)

	got := suite.handler(ctx, msg)
	require.True(t, got.IsOK(), "expected validator stake update to be ok, got %v", got)

	updatedVal, err := keeper.GetValidatorInfo(ctx, oldVal.Signer.Bytes())
	require.Empty(t, err, "unable to fetch validator info %v-", err)
	require.NotEqual(t, stakinginfoStakeUpdate.NewAmount.Int64(), updatedVal.VotingPower, "Validator VotingPower should not be updated to %v", stakinginfoStakeUpdate.NewAmount.Uint64())
}

func (suite *HandlerTestSuite) TestExitedValidatorJoiningAgain() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	accounts := simulation.RandomAccounts(r1, 1)
	pubKey := hmTypes.NewPubKey(accounts[0].PubKey.Bytes())
	signerAddress := hmTypes.HexToThemisAddress(pubKey.Address().Hex())

	txHash := hmTypes.HexToThemisHash("123")
	index := simulation.RandIntBetween(r1, 0, 100)
	logIndex := uint64(index)
	amount, _ := big.NewInt(0).SetString("10000000000000000000", 10)

	// update ACK
	app.StakingKeeper.UpdateL1BatchWithValue(ctx, 20)

	validatorId := hmTypes.NewValidatorID(uint64(1))
	validator := hmTypes.NewValidator(
		validatorId,
		10,
		15,
		1,
		int64(0), // power
		pubKey,
		signerAddress,
	)

	err := app.StakingKeeper.AddValidator(ctx, *validator)
	if err != nil {
		t.Error("Error while adding validator to store", err)
	}

	isCurrentValidator := validator.IsCurrentValidator(14)
	require.False(t, isCurrentValidator)

	totalValidators := app.StakingKeeper.GetAllValidators(ctx)
	require.Equal(t, totalValidators[0].Signer, signerAddress)

	chainParams := app.ChainKeeper.GetParams(ctx)

	txreceipt := &ethTypes.Receipt{
		BlockNumber: big.NewInt(10),
	}
	msgValJoin := types.NewMsgValidatorJoin(
		signerAddress,
		validatorId.Uint64(),
		uint64(1),
		sdk.NewInt(100000),
		pubKey,
		txHash,
		logIndex,
		0,
		1,
	)

	stakinginfoStaked := &stakinginfo.StakinginfoLocked{
		Signer:          signerAddress.EthAddress(),
		SequencerId:     new(big.Int).SetUint64(validatorId.Uint64()),
		ActivationBatch: big.NewInt(1),
		Amount:          amount,
		Total:           big.NewInt(10),
		SignerPubkey:    pubKey.Bytes()[1:],
	}

	suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txreceipt, nil)

	suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress.EthAddress(), txreceipt, msgValJoin.LogIndex).Return(stakinginfoStaked, nil)

	result := suite.handler(ctx, msgValJoin)
	require.True(t, !result.IsOK(), errs.CodeToDefaultMsg(result.Code))
}

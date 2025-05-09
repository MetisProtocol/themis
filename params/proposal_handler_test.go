package params_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/metis-seq/themis/params"
	"github.com/metis-seq/themis/params/subspace"
	"github.com/metis-seq/themis/params/types"
	paramTypes "github.com/metis-seq/themis/params/types"
)

type testInput struct {
	ctx    sdk.Context
	cdc    *codec.Codec
	keeper params.Keeper
}

var (
	_ subspace.ParamSet = (*testParams)(nil)

	keyMaxValidators = "MaxValidators"
	keySlashingRate  = "SlashingRate"
	testSubspace     = "TestSubspace"
)

type testParamsSlashingRate struct {
	DoubleSign uint16 `json:"double_sign,omitempty" yaml:"double_sign,omitempty"`
	Downtime   uint16 `json:"downtime,omitempty" yaml:"downtime,omitempty"`
}

type testParams struct {
	MaxValidators uint16                 `json:"max_validators" yaml:"max_validators"` // maximum number of validators (max uint16 = 65535)
	SlashingRate  testParamsSlashingRate `json:"slashing_rate" yaml:"slashing_rate"`
}

func (tp *testParams) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{Key: []byte(keyMaxValidators), Value: &tp.MaxValidators},
		{Key: []byte(keySlashingRate), Value: &tp.SlashingRate},
	}
}

type invalidParamProposal struct{}

func (invalidParamProposal) GetTitle() string         { return "" }
func (invalidParamProposal) GetDescription() string   { return "" }
func (invalidParamProposal) ProposalRoute() string    { return "" }
func (invalidParamProposal) ProposalType() string     { return "" }
func (invalidParamProposal) ValidateBasic() sdk.Error { return nil }
func (invalidParamProposal) String() string           { return "" }

func testProposal(changes ...paramTypes.ParamChange) paramTypes.ParameterChangeProposal {
	return paramTypes.NewParameterChangeProposal(
		"Test",
		"description",
		changes,
	)
}

func newTestInput(t *testing.T) testInput {
	t.Helper()

	cdc := codec.New()
	types.RegisterCodec(cdc)

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)

	keyParams := sdk.NewKVStoreKey("params")
	tKeyParams := sdk.NewTransientStoreKey("transient_params")

	cms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)

	err := cms.LoadLatestVersion()
	require.Nil(t, err)

	keeper := params.NewKeeper(cdc, keyParams, tKeyParams, paramTypes.DefaultCodespace)
	ctx := sdk.NewContext(cms, abci.Header{}, false, log.NewNopLogger())

	return testInput{ctx, cdc, keeper}
}

func TestProposalHandlerPassed(t *testing.T) {
	input := newTestInput(t)
	ss := input.keeper.Subspace(testSubspace).WithKeyTable(
		subspace.NewKeyTable().RegisterParamSet(&testParams{}),
	)

	tp := testProposal(paramTypes.NewParamChange(testSubspace, keyMaxValidators, "1"))
	hdlr := params.NewParamChangeProposalHandler(input.keeper)
	require.NoError(t, hdlr(input.ctx, tp))

	var param uint16
	ss.Get(input.ctx, []byte(keyMaxValidators), &param)
	require.Equal(t, param, uint16(1))
}

func TestProposalHandlerFailed(t *testing.T) {
	input := newTestInput(t)
	ss := input.keeper.Subspace(testSubspace).WithKeyTable(
		subspace.NewKeyTable().RegisterParamSet(&testParams{}),
	)

	tp := testProposal(paramTypes.NewParamChange(testSubspace, keyMaxValidators, "invalidType"))
	hdlr := params.NewParamChangeProposalHandler(input.keeper)
	require.Error(t, hdlr(input.ctx, tp))

	require.False(t, ss.Has(input.ctx, []byte(keyMaxValidators)))

	require.Error(t, hdlr(input.ctx, invalidParamProposal{}))
}

func TestProposalHandlerSubspaceFailed(t *testing.T) {
	input := newTestInput(t)

	// without subspace
	tp := testProposal(paramTypes.NewParamChange(testSubspace, keySlashingRate, `{"downtime": 7}`))
	hdlr := params.NewParamChangeProposalHandler(input.keeper)
	require.Error(t, hdlr(input.ctx, tp))
}

func TestProposalHandlerUpdateOmitempty(t *testing.T) {
	input := newTestInput(t)
	ss := input.keeper.Subspace(testSubspace).WithKeyTable(
		subspace.NewKeyTable().RegisterParamSet(&testParams{}),
	)

	hdlr := params.NewParamChangeProposalHandler(input.keeper)
	var param testParamsSlashingRate

	tp := testProposal(paramTypes.NewParamChange(testSubspace, keySlashingRate, `{"downtime": 7}`))
	require.NoError(t, hdlr(input.ctx, tp))

	ss.Get(input.ctx, []byte(keySlashingRate), &param)
	require.Equal(t, testParamsSlashingRate{0, 7}, param)

	tp = testProposal(paramTypes.NewParamChange(testSubspace, keySlashingRate, `{"double_sign": 10}`))
	require.NoError(t, hdlr(input.ctx, tp))

	ss.Get(input.ctx, []byte(keySlashingRate), &param)
	require.Equal(t, testParamsSlashingRate{10, 7}, param)
}

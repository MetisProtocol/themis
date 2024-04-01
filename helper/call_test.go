package helper

import (
	"os"
	"testing"

	authTypes "github.com/metis-seq/themis/auth/types"
	"github.com/metis-seq/themis/contracts/erc20"

	"github.com/metis-seq/themis/contracts/stakemanager"
	"github.com/metis-seq/themis/contracts/stakinginfo"
	"github.com/metis-seq/themis/types"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	testTendermintNode = "tcp://localhost:26657"
)

// FetchSigners fetches the signers' list
func FetchSigners(voteBytes []byte, sigInput []byte) ([]string, error) {
	const sigLength = 65

	signersList := make([]string, len(sigInput))

	// Calculate total stake Power of all Signers.
	for i := 0; i < len(sigInput); i += sigLength {
		signature := sigInput[i : i+sigLength]

		pKey, err := authTypes.RecoverPubkey(voteBytes, signature)
		if err != nil {
			return nil, err
		}

		signersList[i] = types.NewPubKey(pKey).Address().String()
	}

	return signersList, nil
}

// TestPopulateABIs tests that package level ABIs cache works as expected
// by not invoking json methods after contracts ABIs' init
func TestPopulateABIs(t *testing.T) {
	t.Parallel()

	viper.Set(TendermintNodeFlag, testTendermintNode)
	viper.Set("log_level", "info")
	InitThemisConfig(os.ExpandEnv("$HOME/.themisd"))

	t.Log("ABIs map should be empty and all ABIs not found")
	assert.True(t, len(ContractsABIsMap) == 0)
	_, found := ContractsABIsMap[stakinginfo.StakinginfoABI]
	assert.False(t, found)
	_, found = ContractsABIsMap[stakemanager.StakemanagerABI]
	assert.False(t, found)
	_, found = ContractsABIsMap[erc20.Erc20ABI]
	assert.False(t, found)

	t.Log("Should create a new contract caller and populate its ABIs by decoding json")

	contractCallerObjFirst, err := NewContractCaller()
	if err != nil {
		t.Error("Error creating contract caller")
	}

	assert.Equalf(t, ContractsABIsMap[stakinginfo.StakinginfoABI], &contractCallerObjFirst.StakingInfoABI,
		"values for %s not equals", stakinginfo.StakinginfoABI)
	assert.Equalf(t, ContractsABIsMap[stakemanager.StakemanagerABI], &contractCallerObjFirst.StakeManagerABI,
		"values for %s not equals", stakemanager.StakemanagerABI)
	assert.Equalf(t, ContractsABIsMap[erc20.Erc20ABI], &contractCallerObjFirst.MetisTokenABI,
		"values for %s not equals", erc20.Erc20ABI)

	t.Log("ABIs map should not be empty and all ABIs found")
	assert.True(t, len(ContractsABIsMap) == 8)
	_, found = ContractsABIsMap[stakinginfo.StakinginfoABI]
	assert.True(t, found)
	_, found = ContractsABIsMap[stakemanager.StakemanagerABI]
	assert.True(t, found)
	_, found = ContractsABIsMap[erc20.Erc20ABI]
	assert.True(t, found)

	t.Log("Should create a new contract caller and populate its ABIs by using cached map")

	contractCallerObjSecond, err := NewContractCaller()
	if err != nil {
		t.Log("Error creating contract caller")
	}

	assert.Equalf(t, ContractsABIsMap[stakinginfo.StakinginfoABI], &contractCallerObjSecond.StakingInfoABI,
		"values for %s not equals", stakinginfo.StakinginfoABI)
	assert.Equalf(t, ContractsABIsMap[stakemanager.StakemanagerABI], &contractCallerObjSecond.StakeManagerABI,
		"values for %s not equals", stakemanager.StakemanagerABI)
	assert.Equalf(t, ContractsABIsMap[erc20.Erc20ABI], &contractCallerObjSecond.MetisTokenABI,
		"values for %s not equals", erc20.Erc20ABI)
}

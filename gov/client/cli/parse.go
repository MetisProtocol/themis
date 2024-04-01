package cli

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"

	govutils "github.com/metis-seq/themis/gov/client/utils"
)

func parseSubmitProposalFlags() (*proposal, error) {
	proposal := &proposal{}
	proposalFile := viper.GetString(FlagProposal)

	if proposalFile == "" {
		proposal.Title = viper.GetString(FlagTitle)
		proposal.Description = viper.GetString(FlagDescription)
		proposal.Type = govutils.NormalizeProposalType(viper.GetString(flagProposalType))
		proposal.Deposit = viper.GetString(FlagDeposit)
		return proposal, nil
	}

	for _, flag := range ProposalFlags {
		if viper.GetString(flag) != "" {
			return nil, fmt.Errorf("--%s flag provided alongside --proposal, which is a noop", flag)
		}
	}

	contents, err := os.ReadFile(proposalFile)
	if err != nil {
		return nil, err
	}

	err = jsoniter.ConfigFastest.Unmarshal(contents, proposal)
	if err != nil {
		return nil, err
	}

	return proposal, nil
}

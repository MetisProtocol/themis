package utils

import "github.com/metis-seq/themis/gov/types"

// NormalizeVoteOption - normalize user specified vote option
func NormalizeVoteOption(option string) string {
	switch option {
	case "Yes", "yes":
		return types.OptionYes.String()

	case "Abstain", "abstain":
		return types.OptionAbstain.String()

	case "No", "no":
		return types.OptionNo.String()

	case "NoWithVeto", "no_with_veto":
		return types.OptionNoWithVeto.String()

	default:
		return ""
	}
}

// NormalizeProposalType - normalize user specified proposal type
func NormalizeProposalType(proposalType string) string {
	switch proposalType {
	default:
		return ""
	}
}

// NormalizeProposalStatus - normalize user specified proposal status
func NormalizeProposalStatus(status string) string {
	switch status {
	case "DepositPeriod", "deposit_period":
		return "DepositPeriod"
	case "VotingPeriod", "voting_period":
		return "VotingPeriod"
	case "Passed", "passed":
		return "Passed"
	case "Rejected", "rejected":
		return "Rejected"
	}
	return ""
}

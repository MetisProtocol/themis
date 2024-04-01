package main

import (
	"errors"
	"math/big"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ethereum/go-ethereum/common"

	chainmanagerTypes "github.com/metis-seq/themis/chainmanager/types"
	"github.com/metis-seq/themis/helper"
	stakingcli "github.com/metis-seq/themis/staking/client/cli"
)

var chainParamsEndpoint = "/chainmanager/params"

// StakeCmd stakes for a validator
func StakeCmd(cliCtx cliContext.CLIContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stake",
		Short: "Stake metis tokens for your account",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			helper.InitThemisConfig("")

			validatorStr := viper.GetString(stakingcli.FlagValidatorAddress)
			stakeAmountStr := viper.GetString(stakingcli.FlagAmount)
			feeAmountStr := viper.GetString(stakingcli.FlagFeeAmount)
			acceptDelegation := viper.GetBool(stakingcli.FlagAcceptDelegation)

			// validator str
			if validatorStr == "" {
				return errors.New("Validator address is required")
			}

			// stake amount
			stakeAmount, ok := big.NewInt(0).SetString(stakeAmountStr, 10)
			if !ok {
				return errors.New("Invalid stake amount")
			}

			// fee amount
			feeAmount, ok := big.NewInt(0).SetString(feeAmountStr, 10)
			if !ok {
				return errors.New("Invalid fee amount")
			}

			// contract caller
			contractCaller, err := helper.NewContractCaller()
			if err != nil {
				return err
			}

			params, err := GetChainmanagerParams(cliCtx)
			if err != nil {
				return err
			}

			stakingManagerAddress := params.ChainParams.StakingManagerAddress.EthAddress()
			stakeManagerInstance, err := contractCaller.GetStakeManagerInstance(stakingManagerAddress)
			if err != nil {
				return err
			}

			return contractCaller.StakeFor(
				common.HexToAddress(validatorStr),
				stakeAmount,
				feeAmount,
				acceptDelegation,
				stakingManagerAddress,
				stakeManagerInstance,
			)
		},
	}

	cmd.Flags().String(stakingcli.FlagValidatorAddress, "", "--validator=<validator address here>")
	cmd.Flags().String(stakingcli.FlagAmount, "10000000000000000000", "--staked-amount=<stake amount>, if left blank it will be assigned as 10 metis tokens")
	cmd.Flags().String(stakingcli.FlagFeeAmount, "5000000000000000000", "--fee-amount=<themis fee amount>, if left blank will be assigned as 5 metis tokens")
	cmd.Flags().Bool(stakingcli.FlagAcceptDelegation, true, "--accept-delegation=<accept delegation>, if left blank will be assigned as true")

	return cmd
}

// ApproveCmd approves tokens for a validator
func ApproveCmd(cliCtx cliContext.CLIContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve",
		Short: "Approve the tokens to stake",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			helper.InitThemisConfig("")

			stakeAmountStr := viper.GetString(stakingcli.FlagAmount)
			feeAmountStr := viper.GetString(stakingcli.FlagFeeAmount)

			// stake amount
			stakeAmount, ok := big.NewInt(0).SetString(stakeAmountStr, 10)
			if !ok {
				return errors.New("Invalid stake amount")
			}

			// fee amount
			feeAmount, ok := big.NewInt(0).SetString(feeAmountStr, 10)
			if !ok {
				return errors.New("Invalid fee amount")
			}

			contractCaller, err := helper.NewContractCaller()
			if err != nil {
				return err
			}

			params, err := GetChainmanagerParams(cliCtx)
			if err != nil {
				return err
			}

			stakingManagerAddress := params.ChainParams.StakingManagerAddress.EthAddress()
			metisTokenAddress := params.ChainParams.MetisTokenAddress.EthAddress()

			metisTokenInstance, err := contractCaller.GetMetisTokenInstance(metisTokenAddress)
			if err != nil {
				return err
			}

			return contractCaller.ApproveTokens(stakeAmount.Add(stakeAmount, feeAmount), stakingManagerAddress, metisTokenAddress, metisTokenInstance)
		},
	}

	cmd.Flags().String(stakingcli.FlagAmount, "10000000000000000000", "--staked-amount=<stake amount>, if left blank will be assigned as 10 metis tokens")
	cmd.Flags().String(stakingcli.FlagFeeAmount, "5000000000000000000", "--fee-amount=<themis fee amount>, if left blank will be assigned as 5 metis tokens")

	return cmd
}

// GetChainmanagerParams return configManager params
func GetChainmanagerParams(cliCtx cliContext.CLIContext) (*chainmanagerTypes.Params, error) {
	response, err := helper.FetchFromAPI(
		cliCtx,
		helper.GetThemisServerEndpoint(chainParamsEndpoint),
	)

	if err != nil {
		return nil, err
	}

	var params chainmanagerTypes.Params
	if err := jsoniter.ConfigFastest.Unmarshal(response.Result, &params); err != nil {
		return nil, err
	}

	return &params, nil
}

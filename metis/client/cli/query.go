package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ethereum/go-ethereum/common"

	hmClient "github.com/metis-seq/themis/client"
	"github.com/metis-seq/themis/helper"
	"github.com/metis-seq/themis/metis/types"
	hmTypes "github.com/metis-seq/themis/types"
	"github.com/metis-seq/themis/version"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	// Group supply queries under a subcommand
	queryCmds := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the metis module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       hmClient.ValidateCmd,
	}

	// clerk query command
	queryCmds.AddCommand(
		client.GetCommands(
			GetSpan(cdc),
			GetLatestSpan(cdc),
			GetQueryParams(cdc),
			GetSpanList(cdc),
			GetNextSpanSeed(cdc),
			GetPreparedProposeSpan(cdc),
		)...,
	)

	return queryCmds
}

// GetSpan get state record
func GetSpan(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "span",
		Short: "show span",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			spanIDStr := viper.GetString(FlagSpanId)
			if spanIDStr == "" {
				return fmt.Errorf("span id cannot be empty")
			}

			spanID, err := strconv.ParseUint(spanIDStr, 10, 64)
			if err != nil {
				return err
			}

			// get query params
			queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQuerySpanParams(spanID))
			if err != nil {
				return err
			}

			// fetch span
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySpan), queryParams)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("Span not found")
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().Uint64(FlagSpanId, 0, "--id=<span ID here>")

	if err := cmd.MarkFlagRequired(FlagSpanId); err != nil {
		cliLogger.Error("GetSpan | MarkFlagRequired | FlagSpanId", "Error", err)
	}

	return cmd
}

// GetLatestSpan get state record
func GetLatestSpan(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "latest-span",
		Short: "show latest span",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// fetch latest span
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLatestSpan), nil)

			// fetch span
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("Latest span not found")
			}

			fmt.Println(string(res))
			return nil
		},
	}

	return cmd
}

// GetQueryParams implements the params query command.
func GetQueryParams(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "show the current metis parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as metis parameters.

Example:
$ %s query metis params
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParams)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			err = jsoniter.ConfigFastest.Unmarshal(bz, &params)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(params)
		},
	}
}

// GetSpan get state record
func GetSpanList(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "spanlist",
		Short: "show span list",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			pageStr := viper.GetString(FlagPage)
			if pageStr == "" {
				return fmt.Errorf("page can't be empty")
			}

			limitStr := viper.GetString(FlagLimit)
			if limitStr == "" {
				return fmt.Errorf("limit can't be empty")
			}

			page, err := strconv.ParseUint(pageStr, 10, 64)
			if err != nil {
				return err
			}

			limit, err := strconv.ParseUint(limitStr, 10, 64)
			if err != nil {
				return err
			}

			// get query params
			queryParams, err := cliCtx.Codec.MarshalJSON(hmTypes.NewQueryPaginationParams(page, limit))
			if err != nil {
				return err
			}

			// query span list
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySpanList), queryParams)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("Span list not found")
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().Uint64(FlagPage, 0, "--page=<page number here>")
	cmd.Flags().Uint64(FlagLimit, 0, "--id=<limit here>")

	if err := cmd.MarkFlagRequired(FlagPage); err != nil {
		cliLogger.Error("GetSpanList | MarkFlagRequired | FlagPage", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagLimit); err != nil {
		cliLogger.Error("GetSpanList | MarkFlagRequired | FlagLimit", "Error", err)
	}

	return cmd
}

// GetNextSpanSeed implements the next span seed.
func GetNextSpanSeed(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "next-span-seed",
		Args:  cobra.NoArgs,
		Short: "show the next span seed",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryNextSpanSeed), nil)
			if err != nil {

				fmt.Println("Error while fetching the span seed")
				return err
			}

			if len(res) == 0 {
				fmt.Println("No span seed found")
				return nil
			}

			fmt.Println(string(res))
			return nil

		},
	}
}

// PostSendProposeSpanTx send propose span transaction
func GetPreparedProposeSpan(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "propose-span",
		Short: "send propose span tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			metisChainID := viper.GetString(FlagMetisChainId)
			if metisChainID == "" {
				return fmt.Errorf("MetisChainID cannot be empty")
			}

			// get proposer
			proposer := hmTypes.HexToThemisAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(cliCtx)
			}

			// start block

			startBlockStr := viper.GetString(FlagStartBlock)
			if startBlockStr == "" {
				return fmt.Errorf("Start block cannot be empty")
			}

			startBlock, err := strconv.ParseUint(startBlockStr, 10, 64)
			if err != nil {
				return err
			}

			// span

			spanIDStr := viper.GetString(FlagSpanId)
			if spanIDStr == "" {
				return fmt.Errorf("Span Id cannot be empty")
			}

			spanID, err := strconv.ParseUint(spanIDStr, 10, 64)
			if err != nil {
				return err
			}

			//
			// Query data
			//

			// fetch duration
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryParams, types.ParamSpan), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				return errors.New("span duration not found")
			}

			var spanDuration uint64
			if err := jsoniter.Unmarshal(res, &spanDuration); err != nil {
				return err
			}

			res, _, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryNextSpanSeed), nil)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("next span seed not found")
			}

			var seed common.Hash
			if err := jsoniter.Unmarshal(res, &seed); err != nil {
				return err
			}

			l2EpochIDStr := viper.GetString(FlagL2EpochId)
			if l2EpochIDStr == "" {
				return fmt.Errorf("l2 Epoch Id cannot be empty")
			}

			l2EpochID, err := strconv.ParseUint(l2EpochIDStr, 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgProposeSpan(
				spanID,
				proposer,
				l2EpochID,
				startBlock,
				startBlock+spanDuration-1,
				metisChainID,
				seed,
				false,
				proposer,
				nil,
			)

			result, err := jsoniter.Marshal(&msg)
			if err != nil {
				return err
			}

			fmt.Println(string(result))
			return nil
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagSpanId, "", "--span-id=<span-id>")
	cmd.Flags().String(FlagL2EpochId, "", "--l2-epoch-id=<l2-epoch-id>")
	cmd.Flags().String(FlagMetisChainId, "", "--metis-chain-id=<metis-chain-id>")
	cmd.Flags().String(FlagStartBlock, "", "--start-block=<start-block-number>")

	if err := cmd.MarkFlagRequired(FlagMetisChainId); err != nil {
		cliLogger.Error("GetPreparedProposeSpan | MarkFlagRequired | FlagMetisChainId", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagStartBlock); err != nil {
		cliLogger.Error("GetPreparedProposeSpan | MarkFlagRequired | FlagStartBlock", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagSpanId); err != nil {
		cliLogger.Error("GetPreparedProposeSpan | MarkFlagRequired | FlagSpanId", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagL2EpochId); err != nil {
		cliLogger.Error("GetPreparedProposeSpan | MarkFlagRequired | FlagL2EpochId", "Error", err)
	}

	return cmd
}

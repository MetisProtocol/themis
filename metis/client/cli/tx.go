package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ethereum/go-ethereum/common"

	hmClient "github.com/metis-seq/themis/client"
	"github.com/metis-seq/themis/helper"
	"github.com/metis-seq/themis/metis/types"
	hmTypes "github.com/metis-seq/themis/types"
)

var cliLogger = helper.Logger.With("module", "metis/client/cli")

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Metis transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       hmClient.ValidateCmd,
	}

	txCmd.AddCommand(
		client.PostCommands(
			PostSendProposeSpanTx(cdc),
		)...,
	)

	return txCmd
}

// PostSendProposeSpanTx send propose span transaction
func PostSendProposeSpanTx(cdc *codec.Codec) *cobra.Command {
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
			if err := jsoniter.ConfigFastest.Unmarshal(res, &spanDuration); err != nil {
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
			if err := jsoniter.ConfigFastest.Unmarshal(res, &seed); err != nil {
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

			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagSpanId, "", "--span-id=<span-id>")
	cmd.Flags().String(FlagL2EpochId, "", "--l2-epoch-id=<l2-epoch-id>")
	cmd.Flags().String(FlagMetisChainId, "", "--metis-chain-id=<metis-chain-id>")
	cmd.Flags().String(FlagStartBlock, "", "--start-block=<start-block-number>")

	if err := cmd.MarkFlagRequired(FlagMetisChainId); err != nil {
		cliLogger.Error("PostSendProposeSpanTx | MarkFlagRequired | FlagMetisChainId", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagL2EpochId); err != nil {
		cliLogger.Error("PostSendProposeSpanTx | MarkFlagRequired | FlagL2EpochId", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagSpanId); err != nil {
		cliLogger.Error("PostSendProposeSpanTx | MarkFlagRequired | FlagSpanId", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagStartBlock); err != nil {
		cliLogger.Error("PostSendProposeSpanTx | MarkFlagRequired | FlagStartBlock", "Error", err)
	}

	return cmd
}

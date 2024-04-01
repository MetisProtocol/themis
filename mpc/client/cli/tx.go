package cli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	hmClient "github.com/metis-seq/themis/client"
	"github.com/metis-seq/themis/helper"
	"github.com/metis-seq/themis/mpc/types"
	hmTypes "github.com/metis-seq/themis/types"
)

var cliLogger = helper.Logger.With("module", "mpc/client/cli")

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Mpc transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       hmClient.ValidateCmd,
	}

	txCmd.AddCommand(
		client.PostCommands(
			PostSendProposeMpcCreateTx(cdc),
			PostSendProposeMpcSignTx(cdc),
			PostUpdateMpcSignTx(cdc),
		)...,
	)

	return txCmd
}

// PostSendProposeMpcCreateTx send propose mpc transaction
func PostSendProposeMpcCreateTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "propose-mpc-create",
		Short: "send propose mpc-create tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// get proposer
			proposer := hmTypes.HexToThemisAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(cliCtx)
			}

			// mpc
			mpcIdStr := viper.GetString(FlagMpcId)
			if mpcIdStr == "" {
				return fmt.Errorf("mpc id cannot be empty")
			}

			threshold := viper.GetUint64(FlagMpcThreshold)
			if threshold <= 0 {
				return fmt.Errorf("invalid mpc threshold")
			}

			parties := viper.GetString(FlagMpcParties)
			if parties == "" {
				return fmt.Errorf("invalid mpc parties")
			}

			var mpcParty []hmTypes.PartyID
			err := json.Unmarshal([]byte(parties), &mpcParty)
			if err != nil {
				return fmt.Errorf("invalid mpc participants err:%v", err)
			}

			mpcAddress := viper.GetString(FlagMpcAddress)
			if mpcAddress == "" {
				return fmt.Errorf("invalid mpc address")
			}

			mpcPubkey := viper.GetString(FlagMpcPubkey)
			if mpcPubkey == "" {
				return fmt.Errorf("invalid mpc pubkey")
			}
			mpbPubkeyBytes, err := hex.DecodeString(mpcPubkey)
			if err != nil {
				return fmt.Errorf("invalid  mpcPubkey err:%v", err)
			}

			mpcType := viper.GetUint64(FlagMpcType)

			msg := types.NewMsgProposeMpcCreate(
				mpcIdStr,
				threshold,
				mpcParty,
				proposer,
				hmTypes.HexToThemisAddress(mpcAddress),
				mpbPubkeyBytes,
				hmTypes.MpcType(mpcType),
			)

			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagMpcId, "", "--mpc-id=<mpc-id>")
	cmd.Flags().Uint64(FlagMpcThreshold, 0, "--threshold=<threshold>")
	cmd.Flags().String(FlagMpcParties, "", "--parties=<parties>")
	cmd.Flags().String(FlagMpcAddress, "", "--mpc-address=<mpc-address>")

	if err := cmd.MarkFlagRequired(FlagMpcId); err != nil {
		cliLogger.Error("PostSendProposeMpcTx | MarkFlagRequired | FlagMpcId", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagMpcThreshold); err != nil {
		cliLogger.Error("PostSendProposeMpcTx | MarkFlagRequired | FlagMpcThreshold", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagMpcParties); err != nil {
		cliLogger.Error("PostSendProposeMpcTx | MarkFlagRequired | FlagMpcParties", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagMpcAddress); err != nil {
		cliLogger.Error("PostSendProposeMpcTx | MarkFlagRequired | FlagMpcAddress", "Error", err)
	}

	return cmd
}

// PostSendProposeMpcSignTx send propose mpc transaction
func PostSendProposeMpcSignTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "propose-mpc-sign",
		Short: "send propose mpc-sign tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// get proposer
			proposer := hmTypes.HexToThemisAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(cliCtx)
			}

			// sign id
			signIdStr := viper.GetString(FlagMpcSignId)
			if signIdStr == "" {
				return fmt.Errorf("mpc sign id cannot be empty")
			}

			// mpc id
			keyIdStr := viper.GetString(FlagMpcId)
			if keyIdStr == "" {
				return fmt.Errorf("mpc id cannot be empty")
			}

			// sign type
			signType := viper.GetUint64(FlagMpcSignType)
			switch signType {
			case uint64(hmTypes.BatchSubmit), uint64(hmTypes.BatchReward), uint64(hmTypes.CommitEpochToMetis), uint64(hmTypes.ReCommitEpochToMetis):
			default:
				return fmt.Errorf("unsupport sign type")
			}

			// sign data
			signDataStr := viper.GetString(FlagMpcSignData)
			if signDataStr == "" {
				return fmt.Errorf("invalid mpc sign data")
			}
			signData, err := hex.DecodeString(strings.TrimPrefix(signDataStr, "0x"))
			if err != nil {
				return fmt.Errorf("invalid mpc sign data err:%v", err)
			}

			// sign msg
			signMsgStr := viper.GetString(FlagMpcSignMsg)
			if signMsgStr == "" {
				return fmt.Errorf("invalid mpc sign msg")
			}
			signMsg, err := hex.DecodeString(strings.TrimPrefix(signMsgStr, "0x"))
			if err != nil {
				return fmt.Errorf("invalid mpc sign msg err:%v", err)
			}

			msg := types.NewMsgProposeMpcSign(
				signIdStr,
				keyIdStr,
				hmTypes.SignType(signType),
				signData,
				signMsg,
				proposer,
			)

			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagMpcSignId, "", "--sign-id=<sign-id>")
	cmd.Flags().String(FlagMpcId, "", "--mpc-id=<mpc-id>")
	cmd.Flags().Uint64(FlagMpcSignType, 0, "--sign-type=<sign-type>")
	cmd.Flags().String(FlagMpcSignData, "", "--sign-data=<sign-data>")
	cmd.Flags().String(FlagMpcSignMsg, "", "--sign-msg=<sign-msg>")

	if err := cmd.MarkFlagRequired(FlagMpcId); err != nil {
		cliLogger.Error("PostSendProposeMpcSignTx | MarkFlagRequired | FlagMpcId", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagMpcSignId); err != nil {
		cliLogger.Error("PostSendProposeMpcSignTx | MarkFlagRequired | FlagMpcSignId", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagMpcSignType); err != nil {
		cliLogger.Error("PostSendProposeMpcSignTx | MarkFlagRequired | FlagMpcSignType", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagMpcSignData); err != nil {
		cliLogger.Error("PostSendProposeMpcSignTx | MarkFlagRequired | FlagMpcSignData", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagMpcSignMsg); err != nil {
		cliLogger.Error("PostSendProposeMpcSignTx | MarkFlagRequired | FlagMpcSignMsg", "Error", err)
	}

	return cmd
}

// PostUpdateMpcSignTx send mpc sign transaction
func PostUpdateMpcSignTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mpc-sign",
		Short: "send mpc-sign tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// get proposer
			proposer := hmTypes.HexToThemisAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(cliCtx)
			}

			// sign id
			signIdStr := viper.GetString(FlagMpcSignId)
			if signIdStr == "" {
				return fmt.Errorf("mpc sign id cannot be empty")
			}

			// signature
			signature := viper.GetString(FlagMpcSignature)
			if signIdStr == "" {
				return fmt.Errorf("mpc signature cannot be empty")
			}

			msg := types.NewMsgMpcSign(
				signIdStr,
				signature,
				proposer,
			)

			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagMpcSignId, "", "--sign-id=<sign-id>")
	cmd.Flags().String(FlagMpcSignature, "", "--signature=<signature>")

	if err := cmd.MarkFlagRequired(FlagMpcSignId); err != nil {
		cliLogger.Error("PostSendProposeMpcSignTx | MarkFlagRequired | FlagMpcSignId", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagMpcSignature); err != nil {
		cliLogger.Error("PostSendProposeMpcSignTx | MarkFlagRequired | FlagMpcSignature", "Error", err)
	}

	return cmd
}

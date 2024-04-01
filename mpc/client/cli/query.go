package cli

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	hmClient "github.com/metis-seq/themis/client"
	"github.com/metis-seq/themis/mpc/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	// Group supply queries under a subcommand
	queryCmds := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the mpc module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       hmClient.ValidateCmd,
	}

	// clerk query command
	queryCmds.AddCommand(
		client.GetCommands(
			GetMpc(cdc),
			GetLatestMpc(cdc),
			GetMpcList(cdc),
			GetMpcSet(cdc),
			GetMpcSign(cdc),
		)...,
	)

	return queryCmds
}

// GetMpc get state record
func GetMpc(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mpc",
		Short: "show mpc info",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			mpcID := viper.GetString(FlagMpcId)
			if mpcID == "" {
				return fmt.Errorf("mpc id cannot be empty")
			}

			// get query params
			queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryMpcParams(mpcID))
			if err != nil {
				return err
			}

			// fetch mpc
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMpc), queryParams)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("mpc not found")
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().String(FlagMpcId, "", "--mpc-id=<mpc ID here>")

	if err := cmd.MarkFlagRequired(FlagMpcId); err != nil {
		cliLogger.Error("GetMpc | MarkFlagRequired | FlagMpcId", "Error", err)
	}

	return cmd
}

// GetLatestMpc get state record
func GetLatestMpc(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "latest-mpc",
		Short: "show latest mpc",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			mpcType := viper.GetInt(FlagMpcType)
			// get query params
			queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryLatestMpcParams(mpcType))
			if err != nil {
				return err
			}

			// fetch latest mpc
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLatestMpc), queryParams)

			// fetch mpc
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("latest mpc not found")
			}

			fmt.Println(string(res))
			return nil
		},
	}

	return cmd
}

// GetMpc get state record
func GetMpcList(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "show mpc list",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// query mpc list
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMpcList), nil)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("mpc list not found")
			}

			fmt.Println(string(res))
			return nil
		},
	}

	return cmd
}

// GetMpcSet get state record
func GetMpcSet(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mpc-set",
		Short: "show mpc parties set",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// query mpc list
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMpcSet), nil)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("mpc list not found")
			}

			fmt.Println(string(res))
			return nil
		},
	}

	return cmd
}

// GetMpcSign get state record
func GetMpcSign(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mpc-sign",
		Short: "show mpc sign info",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			signID := viper.GetString(FlagMpcSignId)
			if signID == "" {
				return fmt.Errorf("mpc id cannot be empty")
			}
			fmt.Println("sign-id:", signID)

			// get query params
			queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryMpcSignParams(signID))
			if err != nil {
				return err
			}
			fmt.Println("query params sign-id:", queryParams)

			// query mpc sign info
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMpcSign), queryParams)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("mpc sign not found")
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().String(FlagMpcSignId, "", "--sign-id=<mpc sign ID here>")

	if err := cmd.MarkFlagRequired(FlagMpcSignId); err != nil {
		cliLogger.Error("GetMpcSign | MarkFlagRequired | FlagMpcSignId", "Error", err)
	}

	return cmd
}

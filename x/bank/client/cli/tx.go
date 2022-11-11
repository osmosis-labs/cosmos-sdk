package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

// NewTxCmd returns a root CLI command handler for all x/bank transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Bank transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(NewSendTxCmd())
	txCmd.AddCommand(NewMultiSendTxCmd())

	return txCmd
}

// NewSendTxCmd returns a CLI command handler for creating a MsgSend transaction.
func NewSendTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "send [from_key_or_address] [to_address] [amount]",
		Short: `Send funds from one account to another. Note, the'--from' flag is
ignored as it is implied from [from_key_or_address].`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[0])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {

				return err
			}

			coins, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {

				return err
			}

			msg := &types.MsgSend{
				FromAddress: clientCtx.GetFromAddress().String(),
				ToAddress:   args[1],
				Amount:      coins,
			}
			if err := msg.ValidateBasic(); err != nil {

				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewMultiSendTxCmd returns a CLI command handler for creating a MsgSend transaction.
func NewMultiSendTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "multisend [from_key_or_address] [to_address_csv] [amount_csv]",
		Short: `Send funds from one account to many other accounts. Note, the'--from' flag is
ignored as it is implied from [from_key_or_address].`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[0])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				fmt.Printf("TEST1\n")
				return err
			}

			// coinsCombined, err := sdk.ParseCoinsNormalized(args[2])
			// if err != nil {
			// 	fmt.Printf("TEST2\n")
			// 	return err
			// }
			// fmt.Printf("%v TEST3\n", coinsCombined)

			toAddresses := args[1]

			toAddressesArray := strings.Split(toAddresses, ",")

			outputs := []types.Output{}

			coinsString := args[2]

			coinsStringArray := strings.Split(coinsString, ",")
			var coinsCombined sdk.Coins

			for i, coinString := range coinsStringArray {
				coins, err := sdk.ParseCoinsNormalized(coinString)
				if err != nil {
					return err
				}
				for _, coin := range coins {
					coinsCombined.Add(coin)
				}
				outputs = append(outputs, types.Output{
					Address: toAddressesArray[i],
					Coins:   coins,
				})
			}

			input := types.Input{
				Address: clientCtx.GetFromAddress().String(),
				Coins:   coinsCombined,
			}

			msg := &types.MsgMultiSend{
				Inputs:  []types.Input{input},
				Outputs: outputs,
			}
			if err := msg.ValidateBasic(); err != nil {
				fmt.Printf("TEST4\n")
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

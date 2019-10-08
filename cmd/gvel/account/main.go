package account

import "github.com/spf13/cobra"

func NewAccountCmd() *cobra.Command {
	cmd := cobra.Command{
		Use: "account <arg>",
		Short: "Use account command for managing the account interacting with Velo",
	}

	cmd.AddCommand(newCreateCmd())

	return &cmd
}

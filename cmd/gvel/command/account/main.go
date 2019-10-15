package account

import (
	"github.com/spf13/cobra"
)

func NewAccountCmd(rootCmd *cobra.Command) {
	accountMainCmd := cobra.Command{
		Use:   "account <arg>",
		Short: "Use account command for managing the account interacting with Velo",
	}

	accountMainCmd.AddCommand(newCreateCmd())

	rootCmd.AddCommand(&accountMainCmd)
}

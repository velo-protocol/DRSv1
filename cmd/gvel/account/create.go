package account

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	cmd := cobra.Command{
		Use: "create",
		Short: "Create an account and store on your disk",
		Run: createAccountRunner,
	}

	return &cmd
}

func createAccountRunner(cmd *cobra.Command, args []string) {
	fmt.Println("TBI")
}

package account

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
)

func (accountCommand *CommandHandler) List(cmd *cobra.Command, args []string) {
	accounts, err := accountCommand.Logic.ListAccount()
	if err != nil {
		console.ExitWithError(console.ExitError, err)
	}

	var data [][]string
	headers := []string{"Index", "Address"}
	for index, account := range *accounts {
		data = append(data, []string{
			fmt.Sprintf("%d", index),
			fmt.Sprintf("%s", account.Address),
		})
	}

	console.WriteTable(headers, data)
}
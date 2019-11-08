package account

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
)

func (accountCommand *CommandHandler) List(cmd *cobra.Command, args []string) {
	console.StartLoading("Getting account details")
	accounts, err := accountCommand.Logic.ListAccount()
	console.StopLoading()
	if err != nil {
		console.ExitWithError(console.ExitError, err)
	}

	var data [][]string
	headers := []string{"Index", "Address", "Default"}
	for index, account := range *accounts {
		data = append(data, []string{
			fmt.Sprintf("%d", index+1),
			fmt.Sprintf("%s", account.Address),
			fmt.Sprintf("%v", account.IsDefault),
		})
	}

	console.WriteTable(headers, data)
}

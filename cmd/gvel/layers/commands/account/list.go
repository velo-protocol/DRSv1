package account

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/console"
)

func (accountCommand *CommandHandler) List(cmd *cobra.Command, args []string) {
	accounts, err := accountCommand.Logic.ListAccount()
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

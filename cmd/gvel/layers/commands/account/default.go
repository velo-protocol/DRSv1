package account

import (
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/validation"
)

func (accountCommand *CommandHandler) Default(cmd *cobra.Command, args []string) {
	address := accountCommand.Prompt.RequestString("Please input public address to be made default", validation.ValidateStellarAddress)

	console.StartLoading("Setting the default account")
	output, err := accountCommand.Logic.SetDefaultAccount(&entity.SetDefaultAccountInput{
		Account: address,
	})
	console.StopLoading()
	if err != nil {
		console.ExitWithError(console.ExitError, err)
	}

	console.Logger.Infof("%s is now set as the default account for signing transaction.", output.Account)
}

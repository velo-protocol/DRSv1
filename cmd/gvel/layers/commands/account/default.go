package account

import (
	"github.com/spf13/cobra"
	"github.com/velo-protocol/DRSv1/cmd/gvel/entity"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/console"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/validation"
)

func (accountCommand *CommandHandler) Default(cmd *cobra.Command, args []string) {
	address := accountCommand.Prompt.RequestString("Please input public address to be made default", validation.ValidateStellarAddress)

	output, err := accountCommand.Logic.SetDefaultAccount(&entity.SetDefaultAccountInput{
		Account: address,
	})
	if err != nil {
		console.ExitWithError(console.ExitError, err)
	}

	console.Logger.Infof("%s is now set as the default account for signing transaction.", output.Account)
}

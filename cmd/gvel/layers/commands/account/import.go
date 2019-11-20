package account

import (
	"github.com/spf13/cobra"
	"github.com/velo-protocol/DRSv1/cmd/gvel/entity"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/console"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/validation"
)

func (accountCommand *CommandHandler) ImportAccount(cmd *cobra.Command, args []string) {
	setAsDefault, err := cmd.Flags().GetBool("default")
	if err != nil {
		console.ExitWithError(console.ExitInvalidInput, err)
	}

	seedKey := accountCommand.Prompt.RequestHiddenString("Please enter seed key of the address", validation.ValidateSeedKey)
	passphrase := accountCommand.Prompt.RequestPassphrase()

	console.StartLoading("Importing account")
	output, err := accountCommand.Logic.ImportAccount(&entity.ImportAccountInput{
		Passphrase:   passphrase,
		SeedKey:      seedKey,
		SetAsDefault: setAsDefault,
	})
	console.StopLoading()
	if err != nil {
		console.ExitWithError(console.ExitError, err)
	}

	if output.IsDefault {
		console.Logger.Printf("Account with address %s is now the default account. Please keep your passphrase safe.", output.ImportedKeyPair.Address())
	} else {
		console.Logger.Printf("Add account with address %s to gvel. Please keep your passphrase safe.", output.ImportedKeyPair.Address())
	}
}

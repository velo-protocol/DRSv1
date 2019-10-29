package account

import (
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
)

func (accountCommand *CommandHandler) Create(cmd *cobra.Command, args []string) {
	setAsDefault, err := cmd.Flags().GetBool("default")
	if err != nil {
		console.ExitWithError(console.ExitInvalidInput, err)
	}

	passphrase := accountCommand.Prompt.RequestPassphrase()

	output, err := accountCommand.Logic.CreateAccount(&entity.CreateAccountInput{
		Passphrase:          passphrase,
		SetAsDefaultAccount: setAsDefault,
	})
	if err != nil {
		console.ExitWithError(console.ExitError, err)
	}

	console.Logger.Printf("A new account is created with address %s Please remember to keep your passphrase safe. You will not be able to recover this passphrase.", output.GeneratedKeyPair.Address())
}

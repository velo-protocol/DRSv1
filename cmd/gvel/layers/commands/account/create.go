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

	console.Logger.Println("generating a new stellar account")
	output, err := accountCommand.Logic.CreateAccount(&entity.CreateAccountInput{
		Passphrase:          passphrase,
		SetAsDefaultAccount: setAsDefault,
	})
	if err != nil {
		console.ExitWithError(console.ExitError, err)
	}

	console.Logger.Printf("%s has been created\n", output.GeneratedKeyPair.Address())
}

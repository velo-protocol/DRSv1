package account

import (
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
)

func (accountCommand *CommandHandler) Create(cmd *cobra.Command, args []string) {
	passphrase := accountCommand.Prompt.RequestPassphrase()

	console.Logger.Println("generating a new stellar account")
	kp, err := accountCommand.Logic.CreateAccount(passphrase)
	if err != nil {
		console.ExitWithError(console.ExitError, err)
	}

	console.Logger.Printf("%s has been created\n", kp.Address())
}

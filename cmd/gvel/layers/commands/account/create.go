package account

import (
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
	errManager "gitlab.com/velo-labs/cen/cmd/gvel/utils/error_manager"
	"log"
)

func (accountCommand *CommandHandler) Create(cmd *cobra.Command, args []string) {
	passphrase := console.RequestPassphrase()

	log.Println("generating a new stellar account")

	kp, err := accountCommand.Logic.CreateAccount(passphrase)
	if err != nil {
		errManager.ExitWithError(errManager.ExitError, err)
	}

	log.Printf("%s has been created\n", kp.Address())
}

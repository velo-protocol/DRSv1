package account

import (
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/constant"
	errManager "gitlab.com/velo-labs/cen/cmd/gvel/error_manager"
	"gitlab.com/velo-labs/cen/cmd/gvel/prompt"
	"gitlab.com/velo-labs/cen/cmd/gvel/util"
	"log"
)

func newCreateCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "create",
		Short: "Create an account and store on your disk",
		Run:   createAccountRunner,
	}

	return &cmd
}

func createAccountRunner(cmd *cobra.Command, args []string) {
	lo, err := util.InitLogic(constant.DefaultGevelAccountDbPath, constant.URL_FRIENDBOT)
	if err != nil {
		panic(err)
	}

	passphrase := prompt.RequestPassphrase()

	log.Println("generating a new stellar account")

	kp, err := lo.CreateAccount(passphrase)
	if err != nil {
		errManager.ExitWithError(errManager.ExitError, err)
	}

	log.Printf("%s has been created\n", kp.Address())
}

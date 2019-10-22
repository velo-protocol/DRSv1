package account

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/config"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/logic"
	errManager "gitlab.com/velo-labs/cen/cmd/gvel/utils/error_manager"
)

type CommandHandler struct {
	Logic logic.Logic
}

func NewCommandHandler(logic logic.Logic) *CommandHandler {
	return &CommandHandler{
		Logic: logic,
	}
}

func (accountCommand *CommandHandler) Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "account <arg>",
		Short: "Use account command for managing the account interacting with Velo",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if !config.Exists() {
				errManager.ExitWithError(errManager.ExitError, errors.New("config file not found, please run `gvel init`"))
			}
		},
	}

	command.AddCommand(
		&cobra.Command{
			Use:   "create",
			Short: "Create an account and store on your disk",
			Run:   accountCommand.Create,
		},
		&cobra.Command{
			Use:   "list",
			Short: "Print all accounts that were created",
			Run:   accountCommand.List,
		},
	)

	return command
}

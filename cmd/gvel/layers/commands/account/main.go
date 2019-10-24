package account

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/logic"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/config"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
)

type CommandHandler struct {
	Logic  logic.Logic
	Prompt console.Prompt
}

func NewCommandHandler(logic logic.Logic, prompt console.Prompt) *CommandHandler {
	return &CommandHandler{
		Logic:  logic,
		Prompt: prompt,
	}
}

func (accountCommand *CommandHandler) Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "account <arg>",
		Short: "Use account command for managing the account interacting with Velo",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if !config.Exists() {
				console.ExitWithError(console.ExitError, errors.New("config file not found, please run `gvel init`"))
			}
		},
	}

	command.AddCommand(
		accountCommand.GetCreateCommand(),
		accountCommand.GetListCommand(),
	)

	return command
}

func (accountCommand *CommandHandler) GetCreateCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "create",
		Short: "Create an account and store on your disk",
		Run:   accountCommand.Create,
	}

	command.Flags().BoolP("default", "d", false, "set as default account")
	return command
}

func (accountCommand *CommandHandler) GetListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Print all accounts that were created",
		Run:   accountCommand.List,
	}
}

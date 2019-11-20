package account

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/logic"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/config"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/console"
)

type CommandHandler struct {
	Logic     logic.Logic
	Prompt    console.Prompt
	AppConfig config.Configuration
}

func NewCommandHandler(logic logic.Logic, prompt console.Prompt, config config.Configuration) *CommandHandler {
	return &CommandHandler{
		Logic:     logic,
		Prompt:    prompt,
		AppConfig: config,
	}
}

func (accountCommand *CommandHandler) Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "account <arg>",
		Short: "Use account command for managing the account interacting with Velo",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if !accountCommand.AppConfig.Exists() {
				console.ExitWithError(console.ExitError, errors.New("config file not found, please run `gvel init`"))
			}
		},
	}

	command.AddCommand(
		accountCommand.GetCreateCommand(),
		accountCommand.GetDefaultCommand(),
		accountCommand.GetImportCommand(),
		accountCommand.GetListCommand(),
		accountCommand.GetExportCommand(),
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

func (accountCommand *CommandHandler) GetDefaultCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "default",
		Short: "Set default account to be used as signer",
		Run:   accountCommand.Default,
	}
}

func (accountCommand *CommandHandler) GetImportCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "import",
		Short: "Import your account",
		Run:   accountCommand.ImportAccount,
	}

	command.Flags().BoolP("default", "d", false, "set as default account")
	return command
}

func (accountCommand *CommandHandler) GetExportCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "export",
		Short: "Export your account",
		Run:   accountCommand.ExportAccount,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if accountCommand.AppConfig.GetDefaultAccount() == "" {
				console.ExitWithError(console.ExitError, errors.New("default account not found in config file, please run `gvel account create` or `gvel account import`"))
			}
		},
	}

	return command
}

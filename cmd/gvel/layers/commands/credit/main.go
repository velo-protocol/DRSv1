package credit

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

func (creditCommand *CommandHandler) Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "credit <arg>",
		Short: "Use credit command for managing the stable credit on Velo",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if !creditCommand.AppConfig.Exists() {
				console.ExitWithError(console.ExitError, errors.New("config file not found, please run `gvel init`"))
			}

			if creditCommand.AppConfig.GetDefaultAccount() == "" {
				console.ExitWithError(console.ExitError, errors.New("default account not found in config file, please run `gvel account create` or `gvel account import`"))
			}
		},
	}

	command.AddCommand(
		creditCommand.GetExchangeRateCommand(),
		creditCommand.GetSetupCommand(),
		creditCommand.GetMintCommand(),
		creditCommand.GetRedeemCommand(),
	)

	return command
}

func (creditCommand *CommandHandler) GetSetupCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "setup",
		Short: "Setup a stable credit on Velo",
		Run:   creditCommand.Setup,
	}

	return command
}

func (creditCommand *CommandHandler) GetMintCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "mint",
		Short: "Mint a stable credit on Velo",
		Run:   creditCommand.Mint,
	}

	return command
}

func (creditCommand *CommandHandler) GetExchangeRateCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "get-exchange",
		Short: "Get exchange rate of a stable credit on Velo",
		Run:   creditCommand.GetExchangeRate,
	}

	return command
}

func (creditCommand *CommandHandler) GetRedeemCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "redeem",
		Short: "Redeemed of a stable credit on Velo",
		Run:   creditCommand.Redeem,
	}

	return command
}

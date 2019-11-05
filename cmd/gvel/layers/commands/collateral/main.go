package collateral

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/logic"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/config"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
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

func (collateralCommand *CommandHandler) Command() *cobra.Command {

	command := &cobra.Command{
		Use:   "collateral <arg>",
		Short: "Use collateral command for managing the collateral asset on Velo",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if !collateralCommand.AppConfig.Exists() {
				console.ExitWithError(console.ExitError, errors.New("config file not found, please run `gvel init`"))
			}

			if collateralCommand.AppConfig.GetDefaultAccount() == "" {
				console.ExitWithError(console.ExitError, errors.New("default account not found in config file, please run `gvel account create` or `gvel account import`"))
			}
		},
	}

	command.AddCommand(
		collateralCommand.GetHealthCheckCommand(),
		collateralCommand.GetRebalanceReserveCommand(),
	)

	return command
}

func (collateralCommand *CommandHandler) GetHealthCheckCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "health-check",
		Short: "Get collateral health check of Velo",
		Run:   collateralCommand.GetHealthCheck,
	}
	return command
}

func (collateralCommand *CommandHandler) GetRebalanceReserveCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "rebalance",
		Short: "Rebalance the Collateral and Reserve pool of Velo",
		Run:   collateralCommand.RebalanceReserve,
	}

	return command
}

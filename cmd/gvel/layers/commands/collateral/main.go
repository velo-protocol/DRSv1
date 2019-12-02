package collateral

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/velo-protocol/DRSv1/cmd/gvel/constants"
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

func (collateralCommand *CommandHandler) Command() *cobra.Command {

	command := &cobra.Command{
		Use:   fmt.Sprintf("%s %s", constants.CmdCollateral, "<arg>"),
		Short: "Use collateral command for managing the collateral asset on Velo",
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
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
		Use:   constants.CmdCollateralHealthCheck,
		Short: "Get collateral health check of Velo",
		Run:   collateralCommand.GetHealthCheck,
	}
	return command
}

func (collateralCommand *CommandHandler) GetRebalanceReserveCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   constants.CmdCollateralRebalance,
		Short: "Rebalance the Collateral and Reserve pool of Velo",
		Run:   collateralCommand.RebalanceReserve,
	}

	return command
}

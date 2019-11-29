package commands

import (
	"github.com/spf13/cobra"
	"github.com/velo-protocol/DRSv1/cmd/gvel/constants"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/commands/account"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/commands/collateral"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/commands/credit"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/commands/initialize"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/logic"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/config"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/console"
)

type GvelHandler struct {
	Logic             logic.Logic
	RootCommand       *cobra.Command
	InitCommand       *cobra.Command
	AccountCommand    *cobra.Command
	CreditCommand     *cobra.Command
	CollateralCommand *cobra.Command
	Prompt            console.Prompt
	AppConfig         config.Configuration
}

func NewGvelHandler(logic logic.Logic, config config.Configuration) *GvelHandler {
	return &GvelHandler{
		Logic:       logic,
		RootCommand: NewGvelRootCommand(),
		Prompt:      console.NewPrompt(),
		AppConfig:   config,
	}
}

func NewGvelRootCommand() *cobra.Command {
	return &cobra.Command{
		Use: constants.CmdRootGvel,
	}
}

func (gvelHandler *GvelHandler) Init() {
	// init InitCommand
	if gvelHandler.InitCommand == nil {
		gvelHandler.InitCommand = initialize.
			NewCommandHandler(gvelHandler.Logic, gvelHandler.Prompt, gvelHandler.AppConfig).
			Command()
	}

	// init AccountCommand
	if gvelHandler.AccountCommand == nil {
		gvelHandler.AccountCommand = account.
			NewCommandHandler(gvelHandler.Logic, gvelHandler.Prompt, gvelHandler.AppConfig).
			Command()
	}

	// init CreditCommand
	if gvelHandler.CreditCommand == nil {
		gvelHandler.CreditCommand = credit.
			NewCommandHandler(gvelHandler.Logic, gvelHandler.Prompt, gvelHandler.AppConfig).
			Command()
	}

	// init CollateralCommand
	if gvelHandler.CollateralCommand == nil {
		gvelHandler.CollateralCommand = collateral.
			NewCommandHandler(gvelHandler.Logic, gvelHandler.Prompt, gvelHandler.AppConfig).
			Command()
	}

	// Add commands to root
	gvelHandler.RootCommand.AddCommand(
		gvelHandler.InitCommand,
		gvelHandler.AccountCommand,
		gvelHandler.CreditCommand,
		gvelHandler.CollateralCommand,
	)
}

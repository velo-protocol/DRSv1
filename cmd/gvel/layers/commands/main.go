package commands

import (
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/commands/account"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/commands/initialize"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/logic"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
	"log"
)

type GvelHandler struct {
	Logic          logic.Logic
	RootCommand    *cobra.Command
	InitCommand    *cobra.Command
	AccountCommand *cobra.Command
	Prompt         console.Prompt
}

func NewGvelHandler(logic logic.Logic) *GvelHandler {
	return &GvelHandler{
		Logic:       logic,
		RootCommand: NewGvelRootCommand(),
		Prompt:      console.NewPrompt(),
	}
}

func NewGvelRootCommand() *cobra.Command {
	return &cobra.Command{
		Use: "gvel",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.Println("IMPORTANT NOTICE: Heavily WIP, expect anything.")
		},
	}
}

func (gvelHandler *GvelHandler) Init() {
	// init InitCommand
	if gvelHandler.InitCommand == nil {
		gvelHandler.InitCommand = initialize.
			NewCommandHandler(gvelHandler.Logic, gvelHandler.Prompt).
			Command()
	}

	// init AccountCommand
	if gvelHandler.AccountCommand == nil {
		gvelHandler.AccountCommand = account.
			NewCommandHandler(gvelHandler.Logic, gvelHandler.Prompt).
			Command()
	}

	// Add commands to root
	gvelHandler.RootCommand.AddCommand(
		gvelHandler.InitCommand,
		gvelHandler.AccountCommand,
	)
}

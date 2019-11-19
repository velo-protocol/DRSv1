package initialize

import (
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

func (initCommand *CommandHandler) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Use init command for initializing all configurations",
		Run:   initCommand.Init,
	}
}

package initialize

import (
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/logic"
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

func (initCommand *CommandHandler) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Use init command for initializing all configurations",
		Run:   initCommand.Init,
	}
}

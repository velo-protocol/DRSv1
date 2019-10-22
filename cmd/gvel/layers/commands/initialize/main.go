package initialize

import (
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/logic"
)

type CommandHandler struct {
	Logic logic.Logic
}

func NewCommandHandler(logic logic.Logic) *CommandHandler {
	return &CommandHandler{
		Logic: logic,
	}
}

func (initCommand *CommandHandler) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Use init command for initializing all configurations",
		Run:   initCommand.Init,
	}
}

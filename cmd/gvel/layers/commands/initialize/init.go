package initialize

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/velo-protocol/DRSv1/cmd/gvel/constants"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/console"
)

func (initCommand *CommandHandler) Init(cmd *cobra.Command, args []string) {
	if initCommand.AppConfig.Exists() {
		console.ExitWithError(console.ExitError, errors.Errorf("gvel has already been initialized, configuration can be found at %s", constants.DefaultConfigFilePath))
	}

	err := initCommand.Logic.Init(constants.DefaultConfigFilePath)
	if err != nil {
		console.ExitWithError(console.ExitError, err)
	}

	console.Logger.Printf("gvel has been initialized\n")
	console.Logger.Printf("using config file at: %s\n", constants.DefaultConfigFilePath)
}

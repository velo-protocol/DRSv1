package initialize

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/config"
	"gitlab.com/velo-labs/cen/cmd/gvel/constants"
	errManager "gitlab.com/velo-labs/cen/cmd/gvel/utils/error_manager"
	"log"
)

func (initCommand *CommandHandler) Init(cmd *cobra.Command, args []string) {
	if config.Exists() {
		errManager.ExitWithError(errManager.ExitError, errors.Errorf("gvel has already been initialized, configuration can be found at %s", constants.DefaultConfigFilePath))
	}

	err := initCommand.Logic.Init(constants.DefaultConfigFilePath)
	if err != nil {
		errManager.ExitWithError(errManager.ExitError, err)
	}

	log.Printf("gvel has been initialized\n")
	log.Printf("using config file at: %s\n", constants.DefaultConfigFilePath)
}

package init

import (
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/constant"
	errManager "gitlab.com/velo-labs/cen/cmd/gvel/error_manager"
	"gitlab.com/velo-labs/cen/cmd/gvel/util"
	"log"
)

func NewInitCmd(rootCmd *cobra.Command) {
	cmd := cobra.Command{
		Use:   "init",
		Short: "Use init command for initializing all configurations",
		Run:   initRunner,
	}

	rootCmd.AddCommand(&cmd)
}

func initRunner(cmd *cobra.Command, args []string) {
	logic := util.InitLogicWithoutDB()

	err := logic.Init(constant.DefaultConfigFilePath)
	if err != nil {
		errManager.ExitWithError(errManager.ExitError, err)
	}

	log.Printf("gvel had been initialized\n")
	log.Printf("using config file at: %s\n", constant.DefaultConfigFilePath)
}

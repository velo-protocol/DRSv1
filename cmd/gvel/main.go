package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/command/account"
	_init "gitlab.com/velo-labs/cen/cmd/gvel/command/init"
	"log"
	"os"
)

func main() {
	var rootCommand = &cobra.Command{
		Use: "gvel",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.Println("IMPORTANT NOTICE: Heavily WIP, expect anything.")
		},
	}

	account.NewAccountCmd(rootCommand)
	_init.NewInitCmd(rootCommand)

	err := rootCommand.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/account"
	_init "gitlab.com/velo-labs/cen/cmd/gvel/init"
	"os"
)

func main() {
	var rootCommand = &cobra.Command{
		Use: "gvel",
	}

	rootCommand.AddCommand(account.NewAccountCmd())
	rootCommand.AddCommand(_init.NewInitCmd())

	err := rootCommand.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

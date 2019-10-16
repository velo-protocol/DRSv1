package account

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/constant"
	errManager "gitlab.com/velo-labs/cen/cmd/gvel/error_manager"
	"gitlab.com/velo-labs/cen/cmd/gvel/prompt"
	"gitlab.com/velo-labs/cen/cmd/gvel/util"
)

func newListCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list",
		Short: "Print all accounts that were created",
		Run:   listAccountRunner,
	}

	return &cmd
}

func listAccountRunner(cmd *cobra.Command, args []string) {
	lo, err := util.InitLogic(constant.DefaultGevelAccountDbPath, constant.UrlFriendbot)
	if err != nil {
		panic(err)
	}

	accounts, err := lo.ListAccount()
	if err != nil {
		errManager.ExitWithError(errManager.ExitError, err)
	}

	var data [][]string
	headers := []string{"Index", "Address"}
	for index, account := range *accounts {
		data = append(data, []string{
			fmt.Sprintf("%d", index),
			fmt.Sprintf("%s", account.Address),
		})
	}

	prompt.TableWriter(headers, data)
}

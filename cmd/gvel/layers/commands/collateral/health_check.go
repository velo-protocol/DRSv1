package collateral

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/console"
)

func (collateralCommand *CommandHandler) GetHealthCheck(cmd *cobra.Command, args []string) {

	console.StartLoading("Getting collateral health check")
	output, err := collateralCommand.Logic.GetCollateralHealthCheck()
	console.StopLoading()

	if err != nil {
		console.ExitWithError(console.ExitError, err)
	}

	var data [][]string
	headers := []string{"Asset", "Collateral pool", "Required collateral"}
	data = append(data, []string{
		fmt.Sprintf("%s", output.Asset),
		fmt.Sprintf("%s", output.PoolAmount),
		fmt.Sprintf("%s", output.RequiredAmount),
	})

	console.WriteTable(headers, data)
}

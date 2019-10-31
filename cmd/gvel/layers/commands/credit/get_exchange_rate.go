package credit

import (
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
)

func (creditCommand *CommandHandler) GetExchangeRate(cmd *cobra.Command, args []string) {

	getExchangeRateInput := &entity.GetExchangeRateInput{
		AssetCode: creditCommand.Prompt.RequestString("Please input asset code of the stable credit", nil),
		Issuer:    creditCommand.Prompt.RequestString("Please input issuing account of the stable credit", nil),
	}

	console.StartLoading("Getting exchange rate")
	output, err := creditCommand.Logic.GetExchangeRate(getExchangeRateInput)
	console.StopLoading()

	if err != nil {
		console.ExitWithError(console.ExitError, err)
	}

	console.Logger.Infof("You will get %s %s for 1 %s.", output.RedeemablePricePerUnit, output.RedeemableCollateral, output.AssetCode)
}

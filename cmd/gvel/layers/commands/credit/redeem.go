package credit

import (
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
)

func (creditCommand *CommandHandler) Redeem(cmd *cobra.Command, args []string) {

	redeemCreditInput := &entity.RedeemCreditInput{
		AssetCode:   creditCommand.Prompt.RequestString("Please input asset code of the stable credit to be redeemed", nil),
		AssetIssuer: creditCommand.Prompt.RequestString("Please input issuing account of the stable credit", nil),
		Amount:      creditCommand.Prompt.RequestString("Please input the amount of stable credit", nil),
		Passphrase:  creditCommand.Prompt.RequestHiddenString("Please input passphrase", nil),
	}

	console.StartLoading("Redeeming %s", redeemCreditInput.AssetCode)
	output, err := creditCommand.Logic.RedeemCredit(redeemCreditInput)
	console.StopLoading()

	if err != nil {
		console.ExitWithError(console.ExitError, err)
	}

	console.Logger.Infof("Redeemed successfully.")
	console.Logger.Infof("🔗 Stellar Transaction Hash: %s", output.TxResult.Hash)
}

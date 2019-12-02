package credit

import (
	"github.com/spf13/cobra"
	"github.com/velo-protocol/DRSv1/cmd/gvel/entity"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/console"
)

func (creditCommand *CommandHandler) Redeem(cmd *cobra.Command, args []string) {

	redeemCreditInput := &entity.RedeemCreditInput{
		AssetCodeToBeRedeemed:   creditCommand.Prompt.RequestString("Please input asset code of the stable credit to be redeemed", nil),
		AssetIssuerToBeRedeemed: creditCommand.Prompt.RequestString("Please input issuing account of the stable credit", nil),
		AmountToBeRedeemed:      creditCommand.Prompt.RequestString("Please input the amount of stable credit", nil),
		Passphrase:              creditCommand.Prompt.RequestHiddenString("🔑 Please input passphrase", nil),
	}

	console.StartLoading("Redeeming %s", redeemCreditInput.AssetCodeToBeRedeemed)
	output, err := creditCommand.Logic.RedeemCredit(redeemCreditInput)
	console.StopLoading()

	if err != nil {
		console.ExitWithError(console.ExitError, err)
	}

	console.Logger.Infof(
		"Redeemed successfully. You got %s %s.",
		output.CollateralAmount,
		output.CollateralCode,
	)
	console.Logger.Infof("🔗 Stellar Transaction Hash: %s", output.TxResult.Hash)
}

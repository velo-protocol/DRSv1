package credit

import (
	"github.com/spf13/cobra"
	"github.com/velo-protocol/DRSv1/cmd/gvel/entity"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/console"
)

func (creditCommand *CommandHandler) Mint(cmd *cobra.Command, args []string) {

	mintCreditInput := &entity.MintCreditInput{
		AssetCodeToBeMinted: creditCommand.Prompt.RequestString("Please input asset code of credit to be minted", nil),
		CollateralAssetCode: creditCommand.Prompt.RequestString("Please input asset code of collateral", nil),
		CollateralAmount:    creditCommand.Prompt.RequestString("Please input amount of collateral", nil),
		Passphrase:          creditCommand.Prompt.RequestHiddenString("🔑 Please input passphrase", nil),
	}

	console.StartLoading("Minting stable credit")
	output, err := creditCommand.Logic.MintCredit(mintCreditInput)
	console.StopLoading()

	if err != nil {
		console.ExitWithError(console.ExitError, err)
	}

	console.Logger.Infof(
		"%s %s minted successfully. The stable credit is in %s",
		output.CollateralAmount,
		output.AssetCodeToBeMinted,
		output.AssetDistributorToBeIssued,
	)
	console.Logger.Infof("🔗 Stellar Transaction Hash: %s", output.TxResult.Hash)
}

package credit

import (
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
)

func (creditCommand *CommandHandler) Mint(cmd *cobra.Command, args []string) {

	mintCreditInput := &entity.MintCreditInput{
		AssetToBeMinted:     creditCommand.Prompt.RequestString("Please input asset code of credit to be minted", nil),
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

	console.Logger.Infof("%s minted successfully.", output.AssetToBeMinted)
	console.Logger.Infof("🔗 Stellar Transaction Hash: %s", output.TxResult.Hash)
}

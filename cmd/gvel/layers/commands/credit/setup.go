package credit

import (
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
)

func (creditCommand *CommandHandler) Setup(cmd *cobra.Command, args []string) {

	setupCreditInput := &entity.SetupCreditInput{
		AssetCode:      creditCommand.Prompt.RequestString("Please input asset code", nil),
		PeggedValue:    creditCommand.Prompt.RequestString("Please input pegged value", nil),
		PeggedCurrency: creditCommand.Prompt.RequestString("Please input pegged currency", nil),
		Passphrase:     creditCommand.Prompt.RequestHiddenString("Please input passphrase", nil),
	}

	console.StartLoading("Setting up stable credit")
	output, err := creditCommand.Logic.SetupCredit(setupCreditInput)
	console.StopLoading()

	if err != nil {
		console.ExitWithError(console.ExitError, err)
	}

	console.Logger.Infof("Stable credit %s set up for account %s successfully.", output.AssetCode, output.SourceAddress)
	console.Logger.Infof("🔗 Stellar Transaction Hash: %s", output.TxResult.Hash)
}

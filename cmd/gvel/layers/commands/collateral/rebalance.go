package collateral

import (
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
)

func (collateralCommand *CommandHandler) RebalanceReserve(cmd *cobra.Command, args []string) {

	rebalanceInput := &entity.RebalanceInput{
		Passphrase: collateralCommand.Prompt.RequestHiddenString("🔑 Please input passphrase", nil),
	}

	console.StartLoading("Rebalancing the Collateral and Reserve pool")
	output, err := collateralCommand.Logic.RebalanceReserve(rebalanceInput)
	console.StopLoading()

	if err != nil {
		console.ExitWithError(console.ExitError, err)
	}

	console.Logger.Infof("Rebalancing completed.")
	console.Logger.Infof("🔗 Stellar Transaction Hash: %s", output.TxResult.Hash)
}

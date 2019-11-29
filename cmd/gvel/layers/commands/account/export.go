package account

import (
	"github.com/spf13/cobra"
	"github.com/velo-protocol/DRSv1/cmd/gvel/entity"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/console"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/validation"
)

func (accountCommand *CommandHandler) ExportAccount(cmd *cobra.Command, args []string) {
	publicKey := accountCommand.Prompt.RequestString("Please input the public key you want to export", validation.ValidateStellarAddress)
	passphrase := accountCommand.Prompt.RequestHiddenString("🔑 Please input the passphrase of the account", nil)

	output, err := accountCommand.Logic.ExportAccount(&entity.ExportAccountInput{
		PublicKey:  publicKey,
		Passphrase: passphrase,
	})
	if err != nil {
		console.ExitWithError(console.ExitError, err)
	}

	console.Logger.Printf("Your public key is: %s", output.ExportedKeyPair.Address())
	console.Logger.Printf("Your seed key is: %s", output.ExportedKeyPair.Seed())
}

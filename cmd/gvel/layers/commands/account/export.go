package account

import (
	"github.com/spf13/cobra"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/validation"
)

func (accountCommand *CommandHandler) ExportAccount(cmd *cobra.Command, args []string) {
	publicKey := accountCommand.Prompt.RequestString("Please input the public key you want to export", validation.ValidateStellarAddress)
	passphrase := accountCommand.Prompt.RequestHiddenString("🔑 Please input the passphrase of the account", nil)

	console.StartLoading("Decrypting seed key")
	output, err := accountCommand.Logic.ExportAccount(&entity.ExportAccountInput{
		PublicKey:  publicKey,
		Passphrase: passphrase,
	})

	console.StopLoading()
	if err != nil {
		console.ExitWithError(console.ExitError, err)
	}

	console.Logger.Printf("Your public key is: %s", output.ExportedKeyPair.Address())
	console.Logger.Printf("Your seed key is: %s", output.ExportedKeyPair.Seed())
}

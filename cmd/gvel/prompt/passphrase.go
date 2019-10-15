package prompt

import (
	"github.com/bgentry/speakeasy"
	errManager "gitlab.com/velo-labs/cen/cmd/gvel/error_manager"
)

func RequestPassphrase() string {
	passphrase, err := speakeasy.Ask("please enter passphrase: ")
	if err != nil {
		errManager.ExitWithError(errManager.ExitBadArgs, err)
	}

	confirm, err := speakeasy.Ask("please repeat a passphrase to confirm: ")
	if err != nil {
		errManager.ExitWithError(errManager.ExitBadArgs, err)
	}

	if passphrase != confirm {
		errManager.ExitWithError(errManager.ExitBadArgs, err)
	}

	return passphrase
}

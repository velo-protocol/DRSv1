package console

import (
	"github.com/bgentry/speakeasy"
	"github.com/pkg/errors"
)

func RequestPassphrase() string {
	passphrase, err := speakeasy.Ask("please enter passphrase: ")
	if err != nil {
		ExitWithError(ExitBadArgs, err)
	}

	confirm, err := speakeasy.Ask("please repeat a passphrase to confirm: ")
	if err != nil {
		ExitWithError(ExitBadArgs, err)
	}

	if passphrase != confirm {
		ExitWithError(ExitBadArgs, errors.New("passphrase does not match"))
	}

	return passphrase
}

type prompt struct{}

func NewPrompt() Prompt {
	return &prompt{}
}

func (prompt *prompt) RequestPassphrase() string {
	return RequestPassphrase()
}

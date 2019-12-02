package console

import (
	"fmt"
	"github.com/bgentry/speakeasy"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
)

func RequestPassphrase() string {
	passphrase, err := speakeasy.Ask("🔑 please enter passphrase: ")
	if err != nil {
		ExitWithError(ExitBadArgs, err)
	}

	confirm, err := speakeasy.Ask("🔑 please repeat a passphrase to confirm: ")
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
	passphrase := prompt.RequestHiddenString("🔑 Please enter passphrase", nil)

	_, err := (&promptui.Prompt{
		Label: "🔑 Please repeat a passphrase to confirm ",
		Mask:  '*',
		Validate: func(s string) error {
			if s != passphrase {
				return errors.New("passphrase does not match")
			}
			return nil
		},
	}).Run()
	if err != nil {
		ExitWithError(ExitBadArgs, err)
	}

	return passphrase
}

func (prompt *prompt) RequestHiddenString(label string, validate promptui.ValidateFunc) string {
	userInput, err := (&promptui.Prompt{
		Label:    fmt.Sprintf("%s ", label),
		Mask:     '*',
		Validate: validate,
	}).Run()

	if err != nil {
		ExitWithError(ExitBadArgs, err)
	}
	return userInput
}

func (prompt *prompt) RequestString(label string, validate promptui.ValidateFunc) string {
	userInput, err := (&promptui.Prompt{
		Label:    fmt.Sprintf("%s ", label),
		Validate: validate,
	}).Run()

	if err != nil {
		ExitWithError(ExitBadArgs, err)
	}

	return userInput
}

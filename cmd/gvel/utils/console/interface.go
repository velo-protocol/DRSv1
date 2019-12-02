package console

import "github.com/manifoldco/promptui"

type Prompt interface {
	RequestPassphrase() string
	RequestString(label string, validate promptui.ValidateFunc) string
	RequestHiddenString(label string, validate promptui.ValidateFunc) string
}

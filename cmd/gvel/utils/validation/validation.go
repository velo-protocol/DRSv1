package validation

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/strkey"
)

func ValidateStellarAddress(input string) error {
	if !strkey.IsValidEd25519PublicKey(input) {
		return errors.New("invalid account format")
	}
	return nil
}

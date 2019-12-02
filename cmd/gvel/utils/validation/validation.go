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

func ValidateSeedKey(input string) error {
	if !strkey.IsValidEd25519SecretSeed(input) {
		return errors.New("invalid seed key format")
	}
	return nil
}

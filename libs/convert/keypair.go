package vconvert

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/strkey"
)

func PublicKeyToKeyPair(publicKey string) (*keypair.FromAddress, error) {
	kp, err := keypair.Parse(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get keyPair from publicKey")
	}

	kpFromAddress, ok := kp.(*keypair.FromAddress)
	if !ok {
		return nil, errors.Wrap(err, "unable to cast KP to keypair.FromAddress")
	}
	return kpFromAddress, nil
}

func SecretKeyToKeyPair(secretKey string) (*keypair.Full, error) {
	seedKey, err := stringToByte32(secretKey)
	if err != nil {
		return nil, errors.Wrap(err, "incorrect signature format")
	}

	kp, err := keypair.FromRawSeed(seedKey)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get keyPair from secretKey")
	}

	return kp, nil
}

func stringToByte32(s string) ([32]byte, error) {
	rawSeed, err := strkey.Decode(strkey.VersionByteSeed, s)
	if err != nil {
		return [32]byte{}, err
	}

	var rawSeed32 [32]byte
	copy(rawSeed32[:], rawSeed)
	return rawSeed32, nil
}

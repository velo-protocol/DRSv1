package vconvert

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/strkey"
)

func SecretKeyToKeyPair(secretKey string) (*keypair.Full, error) {
	seedKey, err := stringToByte32(secretKey)
	if err != nil {
		return nil, errors.New("unable to convert secretKey key to byte")
	}

	kp, err := keypair.FromRawSeed(seedKey)
	if err != nil {
		return nil, errors.New("unable to get keyPair from secretKey key")
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

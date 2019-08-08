package utils

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/strkey"
)

func KpFromSeedString(seed string) (*keypair.Full, error) {
	seedKeyb, err := StringToByte32(seed)
	if err != nil {
		return nil, errors.New("unable to convert seed key to byte")
	}

	kp, err := keypair.FromRawSeed(*seedKeyb)
	if err != nil {
		return nil, errors.New("unable to get keypair from seed key")
	}

	return kp, nil
}

func StringToByte32(s string) (*[32]byte, error) {
	rawSeed, err := strkey.Decode(strkey.VersionByteSeed, s)
	if err != nil {
		return nil, err
	}

	rawSeed32 := [32]byte{}
	for i := 0; i < 32; i++ {
		rawSeed32[i] = rawSeed[i]
	}

	return &rawSeed32, nil
}

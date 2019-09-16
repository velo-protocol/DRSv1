package vcrypto

import (
	"encoding/base64"
	"github.com/stellar/go/keypair"
)

type Interface interface {
	Sign(secretSeed string, message []byte) (string, error)
	Verify(publicKey string, message, signature []byte) error
}

type Crypto struct{}

func New() *Crypto {
	return &Crypto{}
}

func (c *Crypto) Sign(secretSeed string, message []byte) (string, error) {
	kp, err := keypair.Parse(secretSeed)
	if err != nil {
		return "", err
	}

	signature, err := kp.Sign(message)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func (c *Crypto) Verify(publicKey string, message, signature []byte) error {
	kp, err := keypair.Parse(publicKey)
	if err != nil {
		return err
	}

	err = kp.Verify(message, signature)
	if err != nil {
		return err
	}

	return nil
}

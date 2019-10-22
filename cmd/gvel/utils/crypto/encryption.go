package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"github.com/pkg/errors"
	"io"
)

func Encrypt(rawMessage []byte, passphrase string) ([]byte, []byte, error) {
	hasher := sha256.New()
	hasher.Write([]byte(passphrase))

	c, err := aes.NewCipher(hasher.Sum(nil))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to build cipher from passphrase")
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to build gcm")
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to rand nonce")
	}

	encryptedSeed := gcm.Seal(nonce, nonce, rawMessage, nil)

	return encryptedSeed, nonce, nil
}

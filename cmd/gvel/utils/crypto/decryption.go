package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"github.com/pkg/errors"
)

func Decrypt(cipherText []byte, passphrase string) ([]byte, error) {
	hasher := sha256.New()
	hasher.Write([]byte(passphrase))

	c, err := aes.NewCipher(hasher.Sum(nil))
	if err != nil {
		return nil, errors.Wrap(err, "failed to build cipher from passphrase")
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build GCM")
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, errors.New("nonce size is larger than cipher text")
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decipher and authenticate")
	}

	return plaintext, nil
}

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"github.com/pkg/errors"
)

func Decrypt(ciphertext []byte, nonce []byte, passphrase string) ([]byte, error) {
	c, err := aes.NewCipher([]byte(passphrase))
	if err != nil {
		return nil, errors.Wrap(err, "failed to build cipher from passphrase")
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err, "failed to build cipher from")
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
	}

	return plaintext, nil
}

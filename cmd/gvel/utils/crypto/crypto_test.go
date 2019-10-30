package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecrypt(t *testing.T) {
	message := "Hello Velo!"
	passphrase := "password"

	encryptedMessage, _, err := Encrypt([]byte(message), passphrase)
	assert.NoError(t, err)

	plaintext, err := Decrypt(encryptedMessage, passphrase)
	assert.NoError(t, err)
	assert.Equal(t, message, string(plaintext))
}

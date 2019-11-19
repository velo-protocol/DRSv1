package vconvert

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPublicKeyToKeyPair(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		kp, err := PublicKeyToKeyPair("GA2NBNMRJAYAGMD37P3I2TMIKMR7ASRIU7EGDPWPURNVKHBIJOSNB6VE")
		assert.NoError(t, err)
		assert.Equal(t, "GA2NBNMRJAYAGMD37P3I2TMIKMR7ASRIU7EGDPWPURNVKHBIJOSNB6VE", kp.Address())
	})
	t.Run("error, bad public key", func(t *testing.T) {
		_, err := PublicKeyToKeyPair("GA2NBNMRJAYAGMD37P3I2TMIKMR7ASRIU7EGDPWPURNVKHBIJOSNB6VA")
		assert.Error(t, err)
	})
}

func TestSecretKeyToKeyPair(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		kp, err := SecretKeyToKeyPair("SCYHQSJG5JH5SXKQEWTRY3FUFYHQAVH7B6ZQEVLW22A5YYMYPBXLQBTH")
		assert.NoError(t, err)
		assert.Equal(t, "GA2NBNMRJAYAGMD37P3I2TMIKMR7ASRIU7EGDPWPURNVKHBIJOSNB6VE", kp.Address())
	})
	t.Run("success", func(t *testing.T) {
		_, err := SecretKeyToKeyPair("SCYHQSJG5JH5SXKQEWTRY3FUFYHQAVH7B6ZQEVLW22A5YYMYPBXLQBTA")
		assert.Error(t, err)
	})
}

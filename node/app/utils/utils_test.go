package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeBase64(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		result, err := DecodeBase64("SGVsbG8gVmVsbyE=")
		assert.NoError(t, err)
		assert.Equal(t, "Hello Velo!", result)
	})
	t.Run("error, bad string input", func(t *testing.T) {
		_, err := DecodeBase64("AAAAA")
		assert.Error(t, err)
	})
}

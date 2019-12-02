package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateSeedKey(t *testing.T) {
	err := ValidateSeedKey("SCYHQSJG5JH5SXKQEWTRY3FUFYHQAVH7B6ZQEVLW22A5YYMYPBXLQBTH")
	assert.NoError(t, err)

	err = ValidateSeedKey("GA2NBNMRJAYAGMD37P3I2TMIKMR7ASRIU7EGDPWPURNVKHBIJOSNB6VE")
	assert.Error(t, err)
}

func TestValidateStellarAddress(t *testing.T) {
	err := ValidateStellarAddress("GA2NBNMRJAYAGMD37P3I2TMIKMR7ASRIU7EGDPWPURNVKHBIJOSNB6VE")
	assert.NoError(t, err)

	err = ValidateStellarAddress("SCYHQSJG5JH5SXKQEWTRY3FUFYHQAVH7B6ZQEVLW22A5YYMYPBXLQBTH")
	assert.Error(t, err)
}

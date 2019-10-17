package entities_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"testing"
)

func TestGetExchangeRateInput_Validate(t *testing.T) {

	t.Run("Happy", func(t *testing.T) {
		entity := entities.GetExchangeRateInput{
			AssetCode: "vTHB",
			Issuer:    "GCQDOOHRLBZW2A6COOMMWI5RAKGEZPBXSGZ6L6WA7M7GK3ZMHODDRAS3",
		}
		err := entity.Validate()
		assert.NoError(t, err)
	})

	t.Run("Error, asset code is missing", func(t *testing.T) {
		entity := entities.GetExchangeRateInput{
			Issuer: "GCQDOOHRLBZW2A6COOMMWI5RAKGEZPBXSGZ6L6WA7M7GK3ZMHODDRAS3",
		}
		err := entity.Validate()
		assert.Error(t, err)
		assert.Equal(t, fmt.Sprintf("%s %s", constants.AssetCode, constants.ErrMustNotBeBlank), err.Error())
	})

	t.Run("Error, issuer is missing", func(t *testing.T) {
		entity := entities.GetExchangeRateInput{
			AssetCode: "vTHB",
		}
		err := entity.Validate()
		assert.Error(t, err)
		assert.Equal(t, fmt.Sprintf("%s %s", constants.Issuer, constants.ErrMustNotBeBlank), err.Error())
	})

	t.Run("Error, invalid asset code", func(t *testing.T) {
		entity := entities.GetExchangeRateInput{
			AssetCode: "_AssetC0de",
			Issuer:    "GCQDOOHRLBZW2A6COOMMWI5RAKGEZPBXSGZ6L6WA7M7GK3ZMHODDRAS3",
		}
		err := entity.Validate()
		assert.Error(t, err)
		assert.Equal(t, fmt.Sprintf(constants.ErrInvalidFormat, constants.AssetCode), err.Error())
	})

	t.Run("Error, invalid issuer", func(t *testing.T) {
		entity := entities.GetExchangeRateInput{
			AssetCode: "vTHB",
			Issuer:    "WRONG_ISSUER",
		}
		err := entity.Validate()
		assert.Error(t, err)
		assert.Equal(t, fmt.Sprintf(constants.ErrInvalidFormat, constants.Issuer), err.Error())
	})
}

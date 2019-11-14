package subusecases_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/environments"
	"testing"
)

func TestUseCase_GetIssuerAccount(t *testing.T) {
	var (
		peggedValue = decimal.NewFromFloat(15000000)

		vThbIssuerAddress = "GAN6D232HXTF4OHL7J36SAJD3M22H26B2O4QFVRO32OEM523KTMB6Q72"
		vThbAsset         = "vTHB"

		trustedPartnerAddress = publicKey2
		peggedCurrency        = "THB"

		getMockGetIssuerAccountInput = &entities.GetIssuerAccountInput{
			IssuerAddress: vThbIssuerAddress,
		}
	)

	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(vThbIssuerAddress).
			Return(&horizon.Account{
				AccountID: vThbIssuerAddress,
				Sequence:  "1",
				Signers: []horizon.Signer{{
					Key:    env.DrsPublicKey,
					Weight: 1,
				}, {
					Key:    trustedPartnerAddress,
					Weight: 1,
				}, {
					Key:    vThbIssuerAddress,
					Weight: 0,
				}},
				Data: map[string]string{
					"peggedValue":    base64.StdEncoding.EncodeToString([]byte(peggedValue.String())),
					"peggedCurrency": base64.StdEncoding.EncodeToString([]byte(peggedCurrency)),
					"assetCode":      base64.StdEncoding.EncodeToString([]byte(vThbAsset)),
				},
			}, nil)

		output, err := helper.subUseCase.GetIssuerAccount(context.Background(), getMockGetIssuerAccountInput)

		assert.NoError(t, err)
		assert.NotEmpty(t, output)
		assert.Equal(t, vThbIssuerAddress, output.Account.AccountID)
		assert.Equal(t, "1.5", output.PeggedValue.String())
		assert.Equal(t, peggedCurrency, output.PeggedCurrency)
		assert.Equal(t, vThbAsset, output.AssetCode)
	})

	t.Run("Error - fail to get issuer account", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(vThbIssuerAddress).
			Return(nil, errors.New("some error has occurred"))

		_, err := helper.subUseCase.GetIssuerAccount(context.Background(), getMockGetIssuerAccountInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrGetIssuerAccount)
	})

	t.Run("Error - signer count must be 3", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(vThbIssuerAddress).
			Return(&horizon.Account{
				AccountID: vThbIssuerAddress,
				Sequence:  "1",
				Signers:   []horizon.Signer{},
			}, nil)

		_, err := helper.subUseCase.GetIssuerAccount(context.Background(), getMockGetIssuerAccountInput)
		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf(constants.ErrInvalidIssuerAccount, "signer count must be 3"))
	})

	t.Run("Error - pegged value cannot parse base64", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(vThbIssuerAddress).
			Return(&horizon.Account{
				AccountID: vThbIssuerAddress,
				Sequence:  "1",
				Signers: []horizon.Signer{{
					Key:    env.DrsPublicKey,
					Weight: 1,
				}, {
					Key:    trustedPartnerAddress,
					Weight: 1,
				}, {
					Key:    vThbIssuerAddress,
					Weight: 0,
				}},
				Data: map[string]string{
					"peggedValue":    "BAD_VALUE",
					"peggedCurrency": "BAD_VALUE",
				},
			}, nil)

		_, err := helper.subUseCase.GetIssuerAccount(context.Background(), getMockGetIssuerAccountInput)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf(constants.ErrInvalidIssuerAccount, "invalid pegged value format"))
	})

	t.Run("Error - invalid pegged value format", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(vThbIssuerAddress).
			Return(&horizon.Account{
				AccountID: vThbIssuerAddress,
				Sequence:  "1",
				Signers: []horizon.Signer{{
					Key:    env.DrsPublicKey,
					Weight: 1,
				}, {
					Key:    trustedPartnerAddress,
					Weight: 1,
				}, {
					Key:    vThbIssuerAddress,
					Weight: 0,
				}},
				Data: map[string]string{
					"peggedValue":    base64.StdEncoding.EncodeToString([]byte("BAD_VALUE")),
					"peggedCurrency": "BAD_VALUE",
				},
			}, nil)

		_, err := helper.subUseCase.GetIssuerAccount(context.Background(), getMockGetIssuerAccountInput)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf(constants.ErrInvalidIssuerAccount, "invalid pegged value format"))
	})

	t.Run("Error - invalid pegged currency format", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(vThbIssuerAddress).
			Return(&horizon.Account{
				AccountID: vThbIssuerAddress,
				Sequence:  "1",
				Signers: []horizon.Signer{{
					Key:    env.DrsPublicKey,
					Weight: 1,
				}, {
					Key:    trustedPartnerAddress,
					Weight: 1,
				}, {
					Key:    vThbIssuerAddress,
					Weight: 0,
				}},
				Data: map[string]string{
					"peggedValue":    base64.StdEncoding.EncodeToString([]byte(peggedValue.String())),
					"peggedCurrency": "BAD_VALUE",
				},
			}, nil)

		_, err := helper.subUseCase.GetIssuerAccount(context.Background(), getMockGetIssuerAccountInput)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf(constants.ErrInvalidIssuerAccount, "invalid pegged currency format"))
	})

	t.Run("Error - invalid asset code format", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(vThbIssuerAddress).
			Return(&horizon.Account{
				AccountID: vThbIssuerAddress,
				Sequence:  "1",
				Signers: []horizon.Signer{{
					Key:    env.DrsPublicKey,
					Weight: 1,
				}, {
					Key:    trustedPartnerAddress,
					Weight: 1,
				}, {
					Key:    vThbIssuerAddress,
					Weight: 0,
				}},
				Data: map[string]string{
					"peggedValue":    base64.StdEncoding.EncodeToString([]byte(peggedValue.String())),
					"peggedCurrency": base64.StdEncoding.EncodeToString([]byte(peggedCurrency)),
					"assetCode":      "BAD_VALUE",
				},
			}, nil)

		_, err := helper.subUseCase.GetIssuerAccount(context.Background(), getMockGetIssuerAccountInput)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf(constants.ErrInvalidIssuerAccount, "invalid asset code format"))
	})

}

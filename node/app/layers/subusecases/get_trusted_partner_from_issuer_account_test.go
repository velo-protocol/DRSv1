package subusecases_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	env "gitlab.com/velo-labs/cen/node/app/environments"
	"testing"
)

func TestUseCase_GetTrustedPartnerFromIssuerAccount(t *testing.T) {
	var (
		vThbIssuerAddress = "GAN6D232HXTF4OHL7J36SAJD3M22H26B2O4QFVRO32OEM523KTMB6Q72"

		trustedPartnerAddress = publicKey2
	)

	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(trustedPartnerAddress).
			Return(&horizon.Account{
				AccountID: trustedPartnerAddress,
				Sequence:  "1",
			}, nil)

		output, err := helper.subUseCase.GetTrustedPartnerFromIssuerAccount(context.Background(), &entities.GetTrustedPartnerFromIssuerAccountInput{
			IssuerAccount: &horizon.Account{
				AccountID: vThbIssuerAddress,
				Sequence:  "1",
				Signers: []horizon.Signer{
					{
						Key:    env.DrsPublicKey,
						Weight: 1,
					}, {
						Key:    trustedPartnerAddress,
						Weight: 1,
					}, {
						Key:    vThbIssuerAddress,
						Weight: 0,
					},
				},
			},
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, output)
	})

	t.Run("Error - no drs account as a signer in issuer account", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		_, err := helper.subUseCase.GetTrustedPartnerFromIssuerAccount(context.Background(), &entities.GetTrustedPartnerFromIssuerAccountInput{
			IssuerAccount: &horizon.Account{
				AccountID: vThbIssuerAddress,
				Sequence:  "1",
				Signers: []horizon.Signer{
					{
						Key:    "",
						Weight: 1,
					}, {
						Key:    trustedPartnerAddress,
						Weight: 1,
					}, {
						Key:    vThbIssuerAddress,
						Weight: 0,
					},
				},
			},
		})

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf(constants.ErrInvalidIssuerAccount, "no drs as signer"))
	})

	t.Run("Error - fail to get trusted partner account", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(trustedPartnerAddress).
			Return(nil, errors.New(constants.ErrGetTrustedPartnerAccountDetail))

		_, err := helper.subUseCase.GetTrustedPartnerFromIssuerAccount(context.Background(), &entities.GetTrustedPartnerFromIssuerAccountInput{
			IssuerAccount: &horizon.Account{
				AccountID: vThbIssuerAddress,
				Sequence:  "1",
				Signers: []horizon.Signer{
					{
						Key:    env.DrsPublicKey,
						Weight: 1,
					}, {
						Key:    trustedPartnerAddress,
						Weight: 1,
					}, {
						Key:    vThbIssuerAddress,
						Weight: 0,
					},
				},
			},
		})

		assert.Error(t, err)
		assert.EqualError(t, err, constants.ErrGetTrustedPartnerAccountDetail)
	})

}

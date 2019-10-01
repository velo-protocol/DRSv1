package usecases_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/assert"
	vconvert "gitlab.com/velo-labs/cen/libs/convert"
	vtxnbuild "gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	nerrors "gitlab.com/velo-labs/cen/node/app/errors"
	"testing"
)

func TestUseCase_SetupCredit(t *testing.T) {
	var (
		kp1, _ = vconvert.SecretKeyToKeyPair(secretKey1)
		kp2, _ = vconvert.SecretKeyToKeyPair(secretKey2)

		trustedPartnerListAddress = publicKey2
		trustedPartnerMetaAddress = publicKey3
		trustedPartnerMetaEncoded = base64.StdEncoding.EncodeToString([]byte(publicKey3))

		getMockVeloTx = func() *vtxnbuild.VeloTx {
			return &vtxnbuild.VeloTx{
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: publicKey1,
				},
				VeloOp: &vtxnbuild.SetupCredit{
					PeggedValue:    "1.00",
					PeggedCurrency: "THB",
					AssetCode:      "vTHB",
				},
			}
		}
	)

	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		// get tx sender account
		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)

		// get drs account data
		helper.mockStellarRepo.EXPECT().GetDrsAccountData().
			Return(
				&entities.DrsAccountData{
					TrustedPartnerListAddress: trustedPartnerListAddress,
				},
				nil)

		// validate trusted partner role
		helper.mockStellarRepo.EXPECT().GetAccountData(trustedPartnerListAddress).
			Return(map[string]string{publicKey1: trustedPartnerMetaEncoded}, nil)

		// get trusted partner meta
		helper.mockStellarRepo.EXPECT().GetAccountData(trustedPartnerMetaAddress).
			Return(map[string]string{"vSGD_GAN6D232HXTF4OHL7J36SAJD3M22H26B2O4QFVRO32OEM523KTMB6Q72": "R0RXQUZZM1pRSlZEQ0tOVVVOTFZHNTVOVkZCRFpWVlBZRFNGWlIzRURQTEtJWkwzNDRKWkxUNlU="}, nil)

		signedStellarTxXdr, err := helper.useCase.SetupCredit(context.Background(), veloTx)
		assert.NoError(t, err)
		assert.NotNil(t, signedStellarTxXdr)
	})

	t.Run("Error - velo op validation fail", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.SetupCredit{
				PeggedCurrency: "",
			},
		}

		signedStellarTxXdr, err := helper.useCase.SetupCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, signedStellarTxXdr)
		assert.IsType(t, nerrors.ErrInvalidArgument{}, err)
	})

	t.Run("Error - VeloTx missing signer", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign()

		signedStellarTxXdr, err := helper.useCase.SetupCredit(context.Background(), veloTx)

		assert.Nil(t, signedStellarTxXdr)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), constants.ErrSignatureNotFound)
		assert.IsType(t, nerrors.ErrUnAuthenticated{}, err)
	})

	t.Run("Error - VeloTx wrong signer", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp2)

		signedStellarTxXdr, err := helper.useCase.SetupCredit(context.Background(), veloTx)

		assert.Nil(t, signedStellarTxXdr)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), constants.ErrSignatureNotMatchSourceAccount)
		assert.IsType(t, nerrors.ErrUnAuthenticated{}, err)
	})

	t.Run("Error - can't get whitelist", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(nil, errors.New("stellar return error"))

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		signedStellarTxXdr, err := helper.useCase.SetupCredit(context.Background(), veloTx)
		assert.Nil(t, signedStellarTxXdr)
		assert.Contains(t, err.Error(), "fail to get tx sender account")
		assert.IsType(t, nerrors.ErrNotFound{}, err)
	})

	t.Run("Error - can't get DRS account data", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		// get tx sender account
		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)

		// get drs account data
		helper.mockStellarRepo.EXPECT().GetDrsAccountData().
			Return(nil, errors.New("stellar return error"))

		signedStellarTxXdr, err := helper.useCase.SetupCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, signedStellarTxXdr)
		assert.Contains(t, err.Error(), "fail to get data of drs account")
		assert.IsType(t, nerrors.ErrInternal{}, err)
	})

	t.Run("Error - can't get trusted partner role", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		// get tx sender account
		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)

		// get drs account data
		helper.mockStellarRepo.EXPECT().GetDrsAccountData().
			Return(
				&entities.DrsAccountData{
					TrustedPartnerListAddress: trustedPartnerListAddress,
				},
				nil)

		// validate trusted partner role
		helper.mockStellarRepo.EXPECT().GetAccountData(trustedPartnerListAddress).
			Return(nil, errors.New("stellar return error"))

		signedStellarTxXdr, err := helper.useCase.SetupCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, signedStellarTxXdr)
		assert.Contains(t, err.Error(), "fail to get data of trusted partner list account")
		assert.IsType(t, nerrors.ErrInternal{}, err)
	})

	t.Run("Error - don't have trusted partner permission denied", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		// get tx sender account
		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)

		// get drs account data
		helper.mockStellarRepo.EXPECT().GetDrsAccountData().
			Return(
				&entities.DrsAccountData{
					TrustedPartnerListAddress: trustedPartnerListAddress,
				},
				nil)

		// validate trusted partner role
		helper.mockStellarRepo.EXPECT().GetAccountData(trustedPartnerListAddress).
			Return(map[string]string{}, nil)

		signedStellarTxXdr, err := helper.useCase.SetupCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, signedStellarTxXdr)
		assert.Equal(t, err.Error(), fmt.Sprintf(constants.ErrFormatSignerNotHavePermission, constants.VeloOpSetupCredit))
		assert.IsType(t, nerrors.ErrPermissionDenied{}, err)
	})

	t.Run("Error - can't decode trusted partner data", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		// get tx sender account
		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)

		// get drs account data
		helper.mockStellarRepo.EXPECT().GetDrsAccountData().
			Return(
				&entities.DrsAccountData{
					TrustedPartnerListAddress: trustedPartnerListAddress,
				},
				nil)

		// validate trusted partner role
		helper.mockStellarRepo.EXPECT().GetAccountData(trustedPartnerListAddress).
			Return(map[string]string{publicKey1: "BAD_ENCODED_VALUE"}, nil)

		signedStellarTxXdr, err := helper.useCase.SetupCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, signedStellarTxXdr)
		assert.Contains(t, err.Error(), "fail to decode data")
		assert.IsType(t, nerrors.ErrInternal{}, err)
	})

	t.Run("Error - can't get trusted partner meta", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		// get tx sender account
		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)

		// get drs account data
		helper.mockStellarRepo.EXPECT().GetDrsAccountData().
			Return(
				&entities.DrsAccountData{
					TrustedPartnerListAddress: trustedPartnerListAddress,
				},
				nil)

		// validate trusted partner role
		helper.mockStellarRepo.EXPECT().GetAccountData(trustedPartnerListAddress).
			Return(map[string]string{publicKey1: trustedPartnerMetaEncoded}, nil)

		// get trusted partner meta
		helper.mockStellarRepo.EXPECT().GetAccountData(trustedPartnerMetaAddress).
			Return(nil, errors.New("stellar return error"))

		signedStellarTxXdr, err := helper.useCase.SetupCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, signedStellarTxXdr)
		assert.IsType(t, nerrors.ErrInternal{}, err)
	})

	t.Run("Error - the issuing and distribution account for asset code to specified already", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		// get tx sender account
		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)

		// get drs account data
		helper.mockStellarRepo.EXPECT().GetDrsAccountData().
			Return(
				&entities.DrsAccountData{
					TrustedPartnerListAddress: trustedPartnerListAddress,
				},
				nil)

		// validate trusted partner role
		helper.mockStellarRepo.EXPECT().GetAccountData(trustedPartnerListAddress).
			Return(map[string]string{publicKey1: trustedPartnerMetaEncoded}, nil)

		// get trusted partner meta
		helper.mockStellarRepo.EXPECT().GetAccountData(trustedPartnerMetaAddress).
			Return(map[string]string{
				"vSGD_GAN6D232HXTF4OHL7J36SAJD3M22H26B2O4QFVRO32OEM523KTMB6Q72": "R0RXQUZZM1pRSlZEQ0tOVVVOTFZHNTVOVkZCRFpWVlBZRFNGWlIzRURQTEtJWkwzNDRKWkxUNlU=",
				"vTHB_GAHLHUVDHRJ3U3CUOYQRW2TVNRIC6QC6R2MWVCMKYSVYESO5CQMA6PYM": "R0NTWExLS0tFRzdDWE9WVEVTRVI2SDRYNkk0V1lIWkFCVkpNRkxFUU42MlVCTFNMVlhPUEFUSFY="}, nil)

		signedStellarTxXdr, err := helper.useCase.SetupCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, signedStellarTxXdr)
		assert.IsType(t, nerrors.ErrInternal{}, err)
		assert.Equal(t, "the issuing and distribution account for asset code to specified already", err.Error())
	})

}

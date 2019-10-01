package usecases_test

import (
	"context"
	"encoding/base64"
	"github.com/pkg/errors"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/assert"
	vconvert "gitlab.com/velo-labs/cen/libs/convert"
	vtxnbuild "gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
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
			Return(map[string]string{"SGD_GAN6D232HXTF4OHL7J36SAJD3M22H26B2O4QFVRO32OEM523KTMB6Q72": "R0RXQUZZM1pRSlZEQ0tOVVVOTFZHNTVOVkZCRFpWVlBZRFNGWlIzRURQTEtJWkwzNDRKWkxUNlU="}, nil)

		signedStellarTxXdr, err := helper.useCase.SetupCredit(context.Background(), veloTx)
		assert.NoError(t, err)
		assert.NotNil(t, signedStellarTxXdr)
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
	})

	t.Run("Error - can't query on whitelist on stellar", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(nil, errors.New("stellar error"))

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		signedStellarTxXdr, err := helper.useCase.SetupCredit(context.Background(), veloTx)
		assert.Nil(t, signedStellarTxXdr)
		assert.Contains(t, err.Error(), "fail to get tx sender account")
	})

	//t.Run("Error - this user has no permission", func(t *testing.T) {
	//	helper := initTest(t)
	//	defer mockCtrl.Finish()
	//
	//	testHelper.MockWhiteListRepo.EXPECT().
	//		FindOneWhitelist(entities.WhiteListFilter{
	//			StellarPublicKey: &publicKey1,
	//			RoleCode:         pointer.ToString(string(vxdr.RoleTrustedPartner)),
	//		}).
	//		Return(nil, nil)
	//
	//	veloTx := getMockVeloTx()
	//	_ = veloTx.Build()
	//	_ = veloTx.Sign(kp1)
	//
	//	_, err := useCase.SetupCredit(context.Background(), veloTx)
	//	assert.Contains(t, err.Error(), fmt.Sprintf(constants.ErrFormatSignerNotHavePermission, constants.VeloOpSetupCredit))
	//	assert.IsType(t, nerrors.ErrPermissionDenied{}, err)
	//})
	//
	//t.Run("Error - fail to load trusted partner account", func(t *testing.T) {
	//	helper := initTest(t)
	//	defer mockCtrl.Finish()
	//
	//	testHelper.MockWhiteListRepo.EXPECT().
	//		FindOneWhitelist(entities.WhiteListFilter{
	//			StellarPublicKey: &publicKey1,
	//			RoleCode:         pointer.ToString(string(vxdr.RoleTrustedPartner)),
	//		}).
	//		Return(&entities.WhiteList{StellarPublicKey: publicKey1}, nil)
	//
	//	helper.mockStellarRepo.EXPECT().
	//		GetAccount(publicKey1).
	//		Return(nil, errors.New("some error has occurred"))
	//
	//	veloTx := getMockVeloTx()
	//	_ = veloTx.Build()
	//	_ = veloTx.Sign(kp1)
	//
	//	_, err := useCase.SetupCredit(context.Background(), veloTx)
	//	assert.IsType(t, nerrors.ErrNotFound{}, err)
	//})
	//
	//t.Run("Error - fail to build tx, bad tp account format", func(t *testing.T) {
	//	helper := initTest(t)
	//	defer mockCtrl.Finish()
	//
	//	testHelper.MockWhiteListRepo.EXPECT().
	//		FindOneWhitelist(entities.WhiteListFilter{
	//			StellarPublicKey: &publicKey1,
	//			RoleCode:         pointer.ToString(string(vxdr.RoleTrustedPartner)),
	//		}).
	//		Return(&entities.WhiteList{StellarPublicKey: publicKey1}, nil)
	//
	//	helper.mockStellarRepo.EXPECT().
	//		GetAccount(publicKey1).
	//		Return(&horizon.Account{
	//			AccountID: "GBAD_ACCOUNT",
	//		}, nil)
	//
	//	veloTx := getMockVeloTx()
	//	_ = veloTx.Build()
	//	_ = veloTx.Sign(kp1)
	//
	//	_, err := useCase.SetupCredit(context.Background(), veloTx)
	//	assert.IsType(t, nerrors.ErrInternal{}, err)
	//})

}

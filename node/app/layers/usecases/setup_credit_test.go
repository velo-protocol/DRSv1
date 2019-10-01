package usecases_test

import (
	"context"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/assert"
	vconvert "gitlab.com/velo-labs/cen/libs/convert"
	vtxnbuild "gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"testing"
)

func TestUseCase_SetupCredit(t *testing.T) {
	var (
		kp1, _ = vconvert.SecretKeyToKeyPair(secretKey1)
		//kp2, _ = vconvert.SecretKeyToKeyPair(secretKey2)

		trustedPartnerListAddress = "GATOCTTV6EUBWHEXHRDK4GT63GIZJWLM6NT4YPOVKURSUTDU3PN6V6PK"

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
		useCase, testHelper, mockCtrl := initUseCaseTest(t)
		defer mockCtrl.Finish()

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		// get tx sender account
		testHelper.MockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)

		// get drs account data
		testHelper.MockStellarRepo.EXPECT().GetDrsAccountData().
			Return(
				&entities.DrsAccountData{
					TrustedPartnerListAddress: trustedPartnerListAddress,
				},
				nil)

		// validate trusted partner role
		testHelper.MockStellarRepo.EXPECT().GetAccountData(trustedPartnerListAddress).
			Return(map[string]string{publicKey1: "R0FPQVVZSEg1SkxPSDJGVUVEWlpMSlQyRkdCN1NPWlE2TkdVRUJGUE9WR1JBR0s3VFJGTDJVVFI="}, nil)

		testHelper.MockStellarRepo.EXPECT().GetAccountData("GAOAUYHH5JLOH2FUEDZZLJT2FGB7SOZQ6NGUEBFPOVGRAGK7TRFL2UTR").
			Return(map[string]string{"SGD_GAN6D232HXTF4OHL7J36SAJD3M22H26B2O4QFVRO32OEM523KTMB6Q72": "R0RXQUZZM1pRSlZEQ0tOVVVOTFZHNTVOVkZCRFpWVlBZRFNGWlIzRURQTEtJWkwzNDRKWkxUNlU="}, nil)

		signedStellarTxXdr, err := useCase.SetupCredit(context.Background(), veloTx)
		assert.NoError(t, err)
		assert.NotNil(t, signedStellarTxXdr)
	})
	//
	//t.Run("Error - VeloTx missing signer", func(t *testing.T) {
	//	useCase, _, mockCtrl := initUseCaseTest(t)
	//	defer mockCtrl.Finish()
	//
	//	veloTx := getMockVeloTx()
	//	_ = veloTx.Build()
	//	_ = veloTx.Sign()
	//
	//	signedStellarTxXdr, err := useCase.SetupCredit(context.Background(), veloTx)
	//
	//	assert.Nil(t, signedStellarTxXdr)
	//	assert.NotNil(t, err)
	//	assert.Contains(t, err.Error(), constants.ErrSignatureNotFound)
	//})
	//
	//t.Run("Error - VeloTx wrong signer", func(t *testing.T) {
	//	useCase, _, mockCtrl := initUseCaseTest(t)
	//	defer mockCtrl.Finish()
	//
	//	veloTx := getMockVeloTx()
	//	_ = veloTx.Build()
	//	_ = veloTx.Sign(kp2)
	//
	//	signedStellarTxXdr, err := useCase.SetupCredit(context.Background(), veloTx)
	//
	//	assert.Nil(t, signedStellarTxXdr)
	//	assert.NotNil(t, err)
	//	assert.Contains(t, err.Error(), constants.ErrSignatureNotMatchSourceAccount)
	//})
	//
	//t.Run("Error - can't query on whitelist table", func(t *testing.T) {
	//	useCase, testHelper, mockCtrl := initUseCaseTest(t)
	//	defer mockCtrl.Finish()
	//
	//	testHelper.MockWhiteListRepo.EXPECT().
	//		FindOneWhitelist(entities.WhiteListFilter{
	//			StellarPublicKey: &publicKey1,
	//			RoleCode:         pointer.ToString(string(vxdr.RoleTrustedPartner)),
	//		}).
	//		Return(nil, errors.New(constants.ErrToGetDataFromDatabase))
	//
	//	veloTx := getMockVeloTx()
	//	_ = veloTx.Build()
	//	_ = veloTx.Sign(kp1)
	//
	//	_, err := useCase.SetupCredit(context.Background(), veloTx)
	//	assert.Contains(t, err.Error(), constants.ErrToGetDataFromDatabase)
	//})
	//
	//t.Run("Error - this user has no permission", func(t *testing.T) {
	//	useCase, testHelper, mockCtrl := initUseCaseTest(t)
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
	//	useCase, testHelper, mockCtrl := initUseCaseTest(t)
	//	defer mockCtrl.Finish()
	//
	//	testHelper.MockWhiteListRepo.EXPECT().
	//		FindOneWhitelist(entities.WhiteListFilter{
	//			StellarPublicKey: &publicKey1,
	//			RoleCode:         pointer.ToString(string(vxdr.RoleTrustedPartner)),
	//		}).
	//		Return(&entities.WhiteList{StellarPublicKey: publicKey1}, nil)
	//
	//	testHelper.MockStellarRepo.EXPECT().
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
	//	useCase, testHelper, mockCtrl := initUseCaseTest(t)
	//	defer mockCtrl.Finish()
	//
	//	testHelper.MockWhiteListRepo.EXPECT().
	//		FindOneWhitelist(entities.WhiteListFilter{
	//			StellarPublicKey: &publicKey1,
	//			RoleCode:         pointer.ToString(string(vxdr.RoleTrustedPartner)),
	//		}).
	//		Return(&entities.WhiteList{StellarPublicKey: publicKey1}, nil)
	//
	//	testHelper.MockStellarRepo.EXPECT().
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

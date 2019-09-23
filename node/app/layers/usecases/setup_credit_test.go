package usecases_test

import (
	"context"
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/libs/convert"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/errors"
	"testing"
)

func TestUseCase_SetupCredit(t *testing.T) {
	var (
		kp1, _ = vconvert.SecretKeyToKeyPair(secretKey1)
		kp2, _ = vconvert.SecretKeyToKeyPair(secretKey2)

		getMockVeloTx = func() *vtxnbuild.VeloTx {
			return &vtxnbuild.VeloTx{
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: publicKey1,
				},
				VeloOp: &vtxnbuild.SetupCredit{
					PeggedValue:    "1.00",
					PeggedCurrency: "THB",
					AssetName:      "vTHB",
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
		veloTxEnvelope := veloTx.TxEnvelope()

		testHelper.MockWhiteListRepo.EXPECT().
			FindOneWhitelist(entities.WhiteListFilter{
				StellarPublicKey: &publicKey1,
				RoleCode:         pointer.ToString(string(vxdr.RoleTrustedPartner)),
			}).
			Return(&entities.WhiteList{
				StellarPublicKey: publicKey1,
				RoleCode:         string(vxdr.RoleTrustedPartner),
			}, nil)

		testHelper.MockStellarRepo.EXPECT().
			LoadAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)

		signedStellarTxXdr, err := useCase.SetupCredit(context.Background(), veloTxEnvelope)
		assert.NoError(t, err)
		assert.NotNil(t, signedStellarTxXdr)
	})

	t.Run("Error - VeloTx missing signer", func(t *testing.T) {
		useCase, _, mockCtrl := initUseCaseTest(t)
		defer mockCtrl.Finish()

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign()
		veloTxEnvelope := veloTx.TxEnvelope()

		signedStellarTxXdr, err := useCase.SetupCredit(context.Background(), veloTxEnvelope)

		assert.Nil(t, signedStellarTxXdr)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), constants.ErrSignatureNotFound)
	})

	t.Run("Error - VeloTx wrong signer", func(t *testing.T) {
		useCase, _, mockCtrl := initUseCaseTest(t)
		defer mockCtrl.Finish()

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp2)
		veloTxEnvelope := veloTx.TxEnvelope()

		signedStellarTxXdr, err := useCase.SetupCredit(context.Background(), veloTxEnvelope)

		assert.Nil(t, signedStellarTxXdr)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), constants.ErrSignatureNotMatchSourceAccount)
	})

	t.Run("Error - can't query on whitelist table", func(t *testing.T) {
		useCase, testHelper, mockCtrl := initUseCaseTest(t)
		defer mockCtrl.Finish()

		testHelper.MockWhiteListRepo.EXPECT().
			FindOneWhitelist(entities.WhiteListFilter{
				StellarPublicKey: &publicKey1,
				RoleCode:         pointer.ToString(string(vxdr.RoleTrustedPartner)),
			}).
			Return(nil, errors.New(constants.ErrToGetDataFromDatabase))

		veloTxB64, _ := getMockVeloTx().BuildSignEncode(kp1)
		veloTx, _ := vtxnbuild.TransactionFromXDR(veloTxB64)
		envelope := veloTx.TxEnvelope()

		_, err := useCase.SetupCredit(context.Background(), envelope)
		assert.Contains(t, err.Error(), constants.ErrToGetDataFromDatabase)
	})

	t.Run("Error - this user has no permission", func(t *testing.T) {
		useCase, testHelper, mockCtrl := initUseCaseTest(t)
		defer mockCtrl.Finish()

		testHelper.MockWhiteListRepo.EXPECT().
			FindOneWhitelist(entities.WhiteListFilter{
				StellarPublicKey: &publicKey1,
				RoleCode:         pointer.ToString(string(vxdr.RoleTrustedPartner)),
			}).
			Return(nil, nil)

		veloTxB64, _ := getMockVeloTx().BuildSignEncode(kp1)
		veloTx, _ := vtxnbuild.TransactionFromXDR(veloTxB64)
		envelope := veloTx.TxEnvelope()

		_, err := useCase.SetupCredit(context.Background(), envelope)
		assert.Contains(t, err.Error(), fmt.Sprintf(constants.ErrFormatSignerNotHavePermission, "setup credit"))
		assert.IsType(t, nerrors.ErrPermissionDenied{}, err)
	})

	t.Run("Error - fail to load trusted partner account", func(t *testing.T) {
		useCase, testHelper, mockCtrl := initUseCaseTest(t)
		defer mockCtrl.Finish()

		testHelper.MockWhiteListRepo.EXPECT().
			FindOneWhitelist(entities.WhiteListFilter{
				StellarPublicKey: &publicKey1,
				RoleCode:         pointer.ToString(string(vxdr.RoleTrustedPartner)),
			}).
			Return(&entities.WhiteList{StellarPublicKey: publicKey1}, nil)

		testHelper.MockStellarRepo.EXPECT().
			LoadAccount(publicKey1).
			Return(nil, errors.New("some error has occured"))

		veloTxB64, _ := getMockVeloTx().BuildSignEncode(kp1)
		veloTx, _ := vtxnbuild.TransactionFromXDR(veloTxB64)
		envelope := veloTx.TxEnvelope()

		_, err := useCase.SetupCredit(context.Background(), envelope)
		assert.IsType(t, nerrors.ErrNotFound{}, err)
	})

	t.Run("Error - fail to build tx (bad tp account format)", func(t *testing.T) {
		useCase, testHelper, mockCtrl := initUseCaseTest(t)
		defer mockCtrl.Finish()

		testHelper.MockWhiteListRepo.EXPECT().
			FindOneWhitelist(entities.WhiteListFilter{
				StellarPublicKey: &publicKey1,
				RoleCode:         pointer.ToString(string(vxdr.RoleTrustedPartner)),
			}).
			Return(&entities.WhiteList{StellarPublicKey: publicKey1}, nil)

		testHelper.MockStellarRepo.EXPECT().
			LoadAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: "GBAD_ACCOUNT",
			}, nil)

		veloTxB64, _ := getMockVeloTx().BuildSignEncode(kp1)
		veloTx, _ := vtxnbuild.TransactionFromXDR(veloTxB64)
		envelope := veloTx.TxEnvelope()

		_, err := useCase.SetupCredit(context.Background(), envelope)
		assert.IsType(t, nerrors.ErrInternal{}, err)
	})

}

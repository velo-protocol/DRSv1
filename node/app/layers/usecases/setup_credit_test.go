package usecases_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/assert"
	vtxnbuild "gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/node/app/constants"
	nerrors "gitlab.com/velo-labs/cen/node/app/errors"
	"gitlab.com/velo-labs/cen/node/app/testhelpers"
	"testing"
)

func TestUseCase_SetupCredit(t *testing.T) {
	var (
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
			Return(&drsAccountDataEnity, nil)

		// validate trusted partner role
		helper.mockStellarRepo.EXPECT().GetAccountData(testhelpers.TrustedPartnerListKp.Address()).
			Return(map[string]string{publicKey1: base64.StdEncoding.EncodeToString([]byte(publicKey3))}, nil)

		// get trusted partner meta
		helper.mockStellarRepo.EXPECT().GetAccountData(publicKey3).
			Return(map[string]string{"vSGD_GAN6D232HXTF4OHL7J36SAJD3M22H26B2O4QFVRO32OEM523KTMB6Q72": base64.StdEncoding.EncodeToString([]byte("GDWAFY3ZQJVDCKNUUNLVG55NVFBDZVVPYDSFZR3EDPLKIZL344JZLT6U"))}, nil)

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

	t.Run("Error - fail to get tx sender account", func(t *testing.T) {
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
		assert.Contains(t, err.Error(), constants.ErrGetSenderAccount)
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
		assert.Contains(t, err.Error(), constants.ErrGetDrsAccount)
		assert.IsType(t, nerrors.ErrInternal{}, err)
	})

	t.Run("Error - can't get trusted partner list account", func(t *testing.T) {
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
			Return(&drsAccountDataEnity, nil)

		// validate trusted partner role
		helper.mockStellarRepo.EXPECT().GetAccountData(testhelpers.TrustedPartnerListKp.Address()).
			Return(nil, errors.New("stellar return error"))

		signedStellarTxXdr, err := helper.useCase.SetupCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, signedStellarTxXdr)
		assert.Contains(t, err.Error(), constants.ErrGetTrustedPartnerListDataAccount)
		assert.IsType(t, nerrors.ErrInternal{}, err)
	})

	t.Run("Error - permission denied, tx sender is not a trusted partner", func(t *testing.T) {
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
			Return(&drsAccountDataEnity, nil)

		// validate trusted partner role
		helper.mockStellarRepo.EXPECT().GetAccountData(testhelpers.TrustedPartnerListKp.Address()).
			Return(map[string]string{}, nil)

		signedStellarTxXdr, err := helper.useCase.SetupCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, signedStellarTxXdr)
		assert.Equal(t, err.Error(), fmt.Sprintf(constants.ErrFormatSignerNotHavePermission, constants.VeloOpSetupCredit))
		assert.IsType(t, nerrors.ErrPermissionDenied{}, err)
	})

	t.Run("Error - can't decode trusted partner meta address", func(t *testing.T) {
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
			Return(&drsAccountDataEnity, nil)

		// validate trusted partner role
		helper.mockStellarRepo.EXPECT().GetAccountData(testhelpers.TrustedPartnerListKp.Address()).
			Return(map[string]string{publicKey1: "BAD_ENCODED_VALUE"}, nil)

		signedStellarTxXdr, err := helper.useCase.SetupCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, signedStellarTxXdr)
		assert.Contains(t, err.Error(), "fail to decode data")
		assert.IsType(t, nerrors.ErrInternal{}, err)
	})

	t.Run("Error - can't get data of trusted partner meta account", func(t *testing.T) {
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
			Return(&drsAccountDataEnity, nil)

		// validate trusted partner role
		helper.mockStellarRepo.EXPECT().GetAccountData(testhelpers.TrustedPartnerListKp.Address()).
			Return(map[string]string{publicKey1: base64.StdEncoding.EncodeToString([]byte(publicKey3))}, nil)

		// get trusted partner meta
		helper.mockStellarRepo.EXPECT().GetAccountData(publicKey3).
			Return(nil, errors.New("stellar return error"))

		signedStellarTxXdr, err := helper.useCase.SetupCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, signedStellarTxXdr)
		assert.IsType(t, nerrors.ErrInternal{}, err)
	})

	t.Run("Error - asset code has already been used", func(t *testing.T) {
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
			Return(&drsAccountDataEnity, nil)

		// validate trusted partner role
		helper.mockStellarRepo.EXPECT().GetAccountData(testhelpers.TrustedPartnerListKp.Address()).
			Return(map[string]string{publicKey1: base64.StdEncoding.EncodeToString([]byte(publicKey3))}, nil)

		// get trusted partner meta
		helper.mockStellarRepo.EXPECT().GetAccountData(publicKey3).
			Return(map[string]string{
				"vSGD_GAN6D232HXTF4OHL7J36SAJD3M22H26B2O4QFVRO32OEM523KTMB6Q72": base64.StdEncoding.EncodeToString([]byte("GDWAFY3ZQJVDCKNUUNLVG55NVFBDZVVPYDSFZR3EDPLKIZL344JZLT6U")),
				"vTHB_GAHLHUVDHRJ3U3CUOYQRW2TVNRIC6QC6R2MWVCMKYSVYESO5CQMA6PYM": base64.StdEncoding.EncodeToString([]byte("GCSXLKKKEG7CXOVTESER6H4X6I4WYHZABVJMFLEQN62UBLSLVXOPATHV"))}, nil)

		signedStellarTxXdr, err := helper.useCase.SetupCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, signedStellarTxXdr)
		assert.IsType(t, nerrors.ErrInternal{}, err)
		assert.Equal(t, fmt.Sprintf(constants.ErrAssetCodeAlreadyBeenUsed, veloTx.TxEnvelope().VeloTx.VeloOp.Body.SetupCreditOp.AssetCode), err.Error())
	})

}

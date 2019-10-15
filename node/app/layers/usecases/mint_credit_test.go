package usecases_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/errors"
	"gitlab.com/velo-labs/cen/node/app/testhelpers"
	"testing"
)

func TestUseCase_MintCredit(t *testing.T) {

	var (
		collateralAmount  = decimal.NewFromFloat(1000)
		collateralAsset   = "VELO"
		vThbIssuerAccount = "GAN6D232HXTF4OHL7J36SAJD3M22H26B2O4QFVRO32OEM523KTMB6Q72"
		peggedCurrency    = "THB"
		assetToBeIssued   = "vTHB"
		peggedValueStroop = decimal.NewFromFloat(2000)
		peggedValue       = decimal.New(peggedValueStroop.IntPart(), -7)
		medianPrice       = decimal.NewFromFloat(2.5)

		getMockVeloTx = func() *vtxnbuild.VeloTx {
			return &vtxnbuild.VeloTx{
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: publicKey1,
				},
				VeloOp: &vtxnbuild.MintCredit{
					AssetCodeToBeIssued: assetToBeIssued,
					CollateralAssetCode: collateralAsset,
					CollateralAmount:    collateralAmount.String(),
				},
			}
		}
	)

	t.Run("Success", func(t *testing.T) {
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
			Return(map[string]string{"vTHB_" + vThbIssuerAccount: base64.StdEncoding.EncodeToString([]byte("GDWAFY3ZQJVDCKNUUNLVG55NVFBDZVVPYDSFZR3EDPLKIZL344JZLT6U"))}, nil)

		// get issuer account data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(vThbIssuerAccount).
			Return(map[string]string{
				"peggedCurrency": peggedCurrency,
				"peggedValue":    peggedValueStroop.String(),
			}, nil)

		//get median price from price account
		helper.mockStellarRepo.EXPECT().GetMedianPriceFromPriceAccount(drsAccountDataEnity.VeloPriceAddress(vxdr.Currency(peggedCurrency))).
			Return(medianPrice, nil)

		mintAmount := collateralAmount.Mul(medianPrice).Div(peggedValue)

		mintOutput, err := helper.useCase.MintCredit(context.Background(), veloTx)
		assert.NoError(t, err)
		assert.NotNil(t, mintOutput)
		assert.NotEmpty(t, mintOutput.SignedStellarTxXdr)
		assert.Equal(t, collateralAmount.String(), mintOutput.CollateralAmount.String())
		assert.Equal(t, collateralAsset, mintOutput.CollateralAsset)
		assert.Equal(t, mintAmount.String(), mintOutput.MintAmount.String())
		assert.Equal(t, assetToBeIssued, mintOutput.MintCurrency)

	})

	t.Run("Success - large collateral amount", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()
		largeCollateral := decimal.New(92233720368547, -7)

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.MintCredit{
				AssetCodeToBeIssued: assetToBeIssued,
				CollateralAssetCode: collateralAsset,
				CollateralAmount:    largeCollateral.String(),
			},
		}
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
			Return(map[string]string{"vTHB_" + vThbIssuerAccount: base64.StdEncoding.EncodeToString([]byte("GDWAFY3ZQJVDCKNUUNLVG55NVFBDZVVPYDSFZR3EDPLKIZL344JZLT6U"))}, nil)

		// get issuer account data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(vThbIssuerAccount).
			Return(map[string]string{
				"peggedCurrency": peggedCurrency,
				"peggedValue":    peggedValueStroop.String(),
			}, nil)

		//get median price from price account
		helper.mockStellarRepo.EXPECT().GetMedianPriceFromPriceAccount(drsAccountDataEnity.VeloPriceAddress(vxdr.Currency(peggedCurrency))).
			Return(medianPrice, nil)

		mintAmount := largeCollateral.Mul(medianPrice).Div(peggedValue)

		mintOutput, err := helper.useCase.MintCredit(context.Background(), veloTx)
		assert.NoError(t, err)
		assert.NotNil(t, mintOutput)
		assert.NotEmpty(t, mintOutput.SignedStellarTxXdr)
		assert.Equal(t, largeCollateral.String(), mintOutput.CollateralAmount.String())
		assert.Equal(t, collateralAsset, mintOutput.CollateralAsset)
		assert.Equal(t, mintAmount.String(), mintOutput.MintAmount.String())
		assert.Equal(t, assetToBeIssued, mintOutput.MintCurrency)

	})

	t.Run("Error - calculate mint over flow", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()
		largeCollateral := decimal.New(922337203685477, -7)

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.MintCredit{
				AssetCodeToBeIssued: assetToBeIssued,
				CollateralAssetCode: collateralAsset,
				CollateralAmount:    largeCollateral.String(),
			},
		}
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
			Return(map[string]string{"vTHB_" + vThbIssuerAccount: base64.StdEncoding.EncodeToString([]byte("GDWAFY3ZQJVDCKNUUNLVG55NVFBDZVVPYDSFZR3EDPLKIZL344JZLT6U"))}, nil)

		// get issuer account data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(vThbIssuerAccount).
			Return(map[string]string{
				"peggedCurrency": peggedCurrency,
				"peggedValue":    peggedValueStroop.String(),
			}, nil)

		//get median price from price account
		helper.mockStellarRepo.EXPECT().GetMedianPriceFromPriceAccount(drsAccountDataEnity.VeloPriceAddress(vxdr.Currency(peggedCurrency))).
			Return(medianPrice, nil)

		mintOutput, err := helper.useCase.MintCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, mintOutput)
		assert.IsType(t, nerrors.ErrInternal{}, err)
		assert.Contains(t, err.Error(), constants.ErrBuildAndSignTransaction)

	})

	t.Run("Error - velo op validation fail", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.MintCredit{
				CollateralAssetCode: collateralAsset,
				CollateralAmount:    collateralAmount.String(),
			},
		}

		mintOutput, err := helper.useCase.MintCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, mintOutput)
		assert.IsType(t, nerrors.ErrInvalidArgument{}, err)
	})

	t.Run("Error - VeloTx missing signer", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign()

		signedStellarTxXdr, err := helper.useCase.MintCredit(context.Background(), veloTx)

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

		signedStellarTxXdr, err := helper.useCase.MintCredit(context.Background(), veloTx)

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

		signedStellarTxXdr, err := helper.useCase.MintCredit(context.Background(), veloTx)
		assert.Nil(t, signedStellarTxXdr)
		assert.Error(t, err)
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

		signedStellarTxXdr, err := helper.useCase.MintCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, signedStellarTxXdr)
		assert.Contains(t, err.Error(), constants.ErrGetDrsAccountData)
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

		signedStellarTxXdr, err := helper.useCase.MintCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, signedStellarTxXdr)
		assert.Contains(t, err.Error(), constants.ErrGetTrustedPartnerListAccountData)
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

		signedStellarTxXdr, err := helper.useCase.MintCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, signedStellarTxXdr)
		assert.Equal(t, err.Error(), fmt.Sprintf(constants.ErrFormatSignerNotHavePermission, constants.VeloOpMintCredit))
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

		signedStellarTxXdr, err := helper.useCase.MintCredit(context.Background(), veloTx)
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

		signedStellarTxXdr, err := helper.useCase.MintCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, signedStellarTxXdr)
		assert.IsType(t, nerrors.ErrInternal{}, err)
	})

	t.Run("Error - can't decode distribution address", func(t *testing.T) {
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
			Return(map[string]string{"vTHB_" + vThbIssuerAccount: "BAD_ENCODED_VALUE"}, nil)

		mintOutput, err := helper.useCase.MintCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, mintOutput)
		assert.IsType(t, nerrors.ErrInternal{}, err)
		assert.Contains(t, err.Error(), fmt.Sprintf(constants.ErrToDecodeData, "BAD_ENCODED_VALUE"))

	})

	t.Run("Error - empty issuer address", func(t *testing.T) {
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
			Return(map[string]string{}, nil)

		mintOutput, err := helper.useCase.MintCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, mintOutput)

	})

	t.Run("Error - can't get issuer account", func(t *testing.T) {
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
			Return(map[string]string{"vTHB_" + vThbIssuerAccount: base64.StdEncoding.EncodeToString([]byte("GDWAFY3ZQJVDCKNUUNLVG55NVFBDZVVPYDSFZR3EDPLKIZL344JZLT6U"))}, nil)

		// get issuer account data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(vThbIssuerAccount).
			Return(nil, errors.New("stellar return error"))

		mintOutput, err := helper.useCase.MintCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, mintOutput)
		assert.IsType(t, nerrors.ErrInternal{}, err)

	})

	t.Run("Error - pegged currency not found", func(t *testing.T) {
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
			Return(map[string]string{"vTHB_" + vThbIssuerAccount: base64.StdEncoding.EncodeToString([]byte("GDWAFY3ZQJVDCKNUUNLVG55NVFBDZVVPYDSFZR3EDPLKIZL344JZLT6U"))}, nil)

		// get issuer account data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(vThbIssuerAccount).
			Return(map[string]string{
				"peggedValue": peggedValueStroop.String(),
			}, nil)

		mintOutput, err := helper.useCase.MintCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, mintOutput)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Equal(t, fmt.Sprintf(constants.ErrAssetCodeToBeIssuedNotSetup, veloTx.TxEnvelope().VeloTx.VeloOp.Body.MintCreditOp.AssetCodeToBeIssued), err.Error())
	})

	t.Run("Error - pegged value not found", func(t *testing.T) {
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
			Return(map[string]string{"vTHB_" + vThbIssuerAccount: base64.StdEncoding.EncodeToString([]byte("GDWAFY3ZQJVDCKNUUNLVG55NVFBDZVVPYDSFZR3EDPLKIZL344JZLT6U"))}, nil)

		// get issuer account data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(vThbIssuerAccount).
			Return(map[string]string{
				"peggedCurrency": peggedCurrency,
			}, nil)

		mintOutput, err := helper.useCase.MintCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, mintOutput)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Equal(t, fmt.Sprintf(constants.ErrAssetCodeToBeIssuedNotSetup, veloTx.TxEnvelope().VeloTx.VeloOp.Body.MintCreditOp.AssetCodeToBeIssued), err.Error())

	})

	t.Run("Error - can't get median price from price account", func(t *testing.T) {
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
			Return(map[string]string{"vTHB_" + vThbIssuerAccount: base64.StdEncoding.EncodeToString([]byte("GDWAFY3ZQJVDCKNUUNLVG55NVFBDZVVPYDSFZR3EDPLKIZL344JZLT6U"))}, nil)

		// get issuer account data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(vThbIssuerAccount).
			Return(map[string]string{
				"peggedCurrency": peggedCurrency,
				"peggedValue":    peggedValueStroop.String(),
			}, nil)

		//get median price from price account
		helper.mockStellarRepo.EXPECT().GetMedianPriceFromPriceAccount(drsAccountDataEnity.VeloPriceAddress(vxdr.Currency(peggedCurrency))).
			Return(decimal.Zero, errors.New("stellar return error"))

		mintOutput, err := helper.useCase.MintCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, mintOutput)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)

	})

	t.Run("Error - pegged value less than zero", func(t *testing.T) {
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
			Return(map[string]string{"vTHB_" + vThbIssuerAccount: base64.StdEncoding.EncodeToString([]byte("GDWAFY3ZQJVDCKNUUNLVG55NVFBDZVVPYDSFZR3EDPLKIZL344JZLT6U"))}, nil)

		// get issuer account data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(vThbIssuerAccount).
			Return(map[string]string{
				"peggedCurrency": peggedCurrency,
				"peggedValue":    "-1",
			}, nil)

		//get median price from price account
		helper.mockStellarRepo.EXPECT().GetMedianPriceFromPriceAccount(drsAccountDataEnity.VeloPriceAddress(vxdr.Currency(peggedCurrency))).
			Return(medianPrice, nil)

		mintOutput, err := helper.useCase.MintCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Nil(t, mintOutput)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Equal(t, constants.ErrPeggedValueMustBeGreaterThanZero, err.Error())

	})
}

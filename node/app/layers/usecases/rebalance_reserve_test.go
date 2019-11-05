package usecases_test

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/protocols/horizon/base"
	"github.com/stellar/go/support/render/hal"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/environments"
	"gitlab.com/velo-labs/cen/node/app/errors"
	"gitlab.com/velo-labs/cen/node/app/testhelpers"
	"testing"
)

func TestUseCase_RebalanceReserve(t *testing.T) {
	testhelpers.InitEnv()

	var (
		trustedPartnerAddress1     = publicKey2
		trustedPartnerMetaAddress1 = publicKey3
		//trustedPartnerAddress2     = publicKey4
		//trustedPartnerMetaAddress2 = publicKey5

		medianPriceThb = decimal.NewFromFloat(10000000)
		medianPriceUsd = decimal.NewFromFloat(10000000)
		medianPriceSgd = decimal.NewFromFloat(10000000)

		stableCreditAsset1  = "vUSD"
		stableCreditIssuer1 = publicKey1
		//stableCreditAsset2  = "vTHB"
		//stableCreditIssuer2 = publicKey1

		drsHighCollateralAmount   = decimal.NewFromFloat(2050.125)
		drsSmallCollateralAmount  = decimal.NewFromFloat(205.125)
		drsMediumCollateralAmount = decimal.NewFromFloat(1537.5075)
		collateralAssetCode       = string(vxdr.AssetVELO)
		collateralAssetIssuer     = env.VeloIssuerPublicKey

		getMockVeloTx = func() *vtxnbuild.VeloTx {
			return &vtxnbuild.VeloTx{
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: publicKey1,
				},
				VeloOp: &vtxnbuild.RebalanceReserve{},
			}
		}
	)

	t.Run("Success, drs collateral amount greater than drs collateral required amount", func(t *testing.T) {
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

		// get drs account
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)

		// get median price thb
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.New(medianPriceThb.IntPart(), -7), nil)

		// get median price usd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceUsdVeloAddress).
			Return(decimal.New(medianPriceUsd.IntPart(), -7), nil)

		// get median price sgd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceSgdVeloAddress).
			Return(decimal.New(medianPriceSgd.IntPart(), -7), nil)

		// get tp list data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(drsAccountDataEnity.TrustedPartnerListAddress).
			Return(map[string]string{trustedPartnerAddress1: trustedPartnerMetaAddress1}, nil)

		// calculate collateral amount
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(trustedPartnerMetaAddress1).
			Return(map[string]string{fmt.Sprintf("%s_%s", stableCreditAsset1, stableCreditIssuer1): "GCDOC2AYBMEESYXYD3NBPFHWAA44PHQKGTKHRDZXQLWJRWOIW5X2MTFQ"}, nil)

		mockGetIssuerAccountOutput := &entities.GetIssuerAccountOutput{
			Account:        nil,
			PeggedValue:    decimal.NewFromFloat(1.5),
			PeggedCurrency: "USD",
			AssetCode:      stableCreditAsset1,
		}

		helper.mockSubUseCase.EXPECT().GetIssuerAccount(context.Background(), &entities.GetIssuerAccountInput{IssuerAddress: stableCreditIssuer1}).Return(
			mockGetIssuerAccountOutput, nil)

		helper.mockStellarRepo.EXPECT().GetAsset(entities.GetAssetInput{
			AssetCode:   stableCreditAsset1,
			AssetIssuer: stableCreditIssuer1,
		}).Return(&horizon.AssetsPage{
			Links: hal.Links{},
			Embedded: struct {
				Records []horizon.AssetStat
			}{
				Records: []horizon.AssetStat{
					{
						Amount: "1025.0050000",
					},
				},
			},
		}, nil)

		helper.mockStellarRepo.EXPECT().GetAccountBalances(env.DrsPublicKey).Return(
			[]horizon.Balance{
				{
					Balance: drsHighCollateralAmount.String(),
					Asset: base.Asset{
						Code:   collateralAssetCode,
						Issuer: collateralAssetIssuer,
					},
				},
			}, nil)

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, collateralAssetIssuer, output.Collaterals[0].AssetIssuer)
		assert.Equal(t, collateralAssetCode, output.Collaterals[0].AssetCode)
		assert.Equal(t, drsHighCollateralAmount.String(), output.Collaterals[0].PoolAmount.String())
		assert.True(t, output.Collaterals[0].RequiredAmount.GreaterThan(decimal.Zero))

		stellarTx, txErr := txnbuild.TransactionFromXDR(*output.SignedStellarTxXdr)
		assert.NoError(t, txErr)

		txEnv := stellarTx.TxEnvelope()
		assert.Equal(t, kp1.Address(), txEnv.Tx.SourceAccount.Address())
		assert.Equal(t, xdr.OperationTypePayment, txEnv.Tx.Operations[0].Body.Type)
		assert.Equal(t, env.DrsPublicKey, txEnv.Tx.Operations[0].SourceAccount.Address())
		assert.Equal(t, drsAccountDataEnity.DrsReserve, txEnv.Tx.Operations[0].Body.PaymentOp.Destination.Address())
		assert.Equal(t, collateralAssetIssuer, txEnv.Tx.Operations[0].Body.PaymentOp.Asset.AlphaNum4.Issuer.Address())
		assert.Equal(t, collateralAssetCode, func() string {
			bytes, _ := txEnv.Tx.Operations[0].Body.PaymentOp.Asset.AlphaNum4.AssetCode.MarshalBinary()
			return string(bytes)
		}())

	})

	t.Run("Success, drs collateral required amount greater than drs collateral amount", func(t *testing.T) {
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

		// get drs account
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)

		// get median price thb
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.New(medianPriceThb.IntPart(), -7), nil)

		// get median price usd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceUsdVeloAddress).
			Return(decimal.New(medianPriceUsd.IntPart(), -7), nil)

		// get median price sgd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceSgdVeloAddress).
			Return(decimal.New(medianPriceSgd.IntPart(), -7), nil)

		// get tp list data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(drsAccountDataEnity.TrustedPartnerListAddress).
			Return(map[string]string{trustedPartnerAddress1: trustedPartnerMetaAddress1}, nil)

		// calculate collateral amount
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(trustedPartnerMetaAddress1).
			Return(map[string]string{fmt.Sprintf("%s_%s", stableCreditAsset1, stableCreditIssuer1): "GCDOC2AYBMEESYXYD3NBPFHWAA44PHQKGTKHRDZXQLWJRWOIW5X2MTFQ"}, nil)

		mockGetIssuerAccountOutput := &entities.GetIssuerAccountOutput{
			Account:        nil,
			PeggedValue:    decimal.NewFromFloat(1.5),
			PeggedCurrency: "USD",
			AssetCode:      stableCreditAsset1,
		}

		helper.mockSubUseCase.EXPECT().GetIssuerAccount(context.Background(), &entities.GetIssuerAccountInput{IssuerAddress: stableCreditIssuer1}).Return(
			mockGetIssuerAccountOutput, nil)

		helper.mockStellarRepo.EXPECT().GetAsset(entities.GetAssetInput{
			AssetCode:   stableCreditAsset1,
			AssetIssuer: stableCreditIssuer1,
		}).Return(&horizon.AssetsPage{
			Links: hal.Links{},
			Embedded: struct {
				Records []horizon.AssetStat
			}{
				Records: []horizon.AssetStat{
					{
						Amount: "1025.0050000",
					},
				},
			},
		}, nil)

		helper.mockStellarRepo.EXPECT().GetAccountBalances(env.DrsPublicKey).Return(
			[]horizon.Balance{
				{
					Balance: drsSmallCollateralAmount.String(),
					Asset: base.Asset{
						Code:   collateralAssetCode,
						Issuer: collateralAssetIssuer,
					},
				},
			}, nil)

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, collateralAssetIssuer, output.Collaterals[0].AssetIssuer)
		assert.Equal(t, collateralAssetCode, output.Collaterals[0].AssetCode)
		assert.Equal(t, drsSmallCollateralAmount.String(), output.Collaterals[0].PoolAmount.String())
		assert.True(t, output.Collaterals[0].RequiredAmount.GreaterThan(decimal.Zero))

		stellarTx, txErr := txnbuild.TransactionFromXDR(*output.SignedStellarTxXdr)
		assert.NoError(t, txErr)

		txEnv := stellarTx.TxEnvelope()
		assert.Equal(t, kp1.Address(), txEnv.Tx.SourceAccount.Address())
		assert.Equal(t, xdr.OperationTypePayment, txEnv.Tx.Operations[0].Body.Type)
		assert.Equal(t, drsAccountDataEnity.DrsReserve, txEnv.Tx.Operations[0].SourceAccount.Address())
		assert.Equal(t, env.DrsPublicKey, txEnv.Tx.Operations[0].Body.PaymentOp.Destination.Address())
		assert.Equal(t, collateralAssetIssuer, txEnv.Tx.Operations[0].Body.PaymentOp.Asset.AlphaNum4.Issuer.Address())
		assert.Equal(t, collateralAssetCode, func() string {
			bytes, _ := txEnv.Tx.Operations[0].Body.PaymentOp.Asset.AlphaNum4.AssetCode.MarshalBinary()
			return string(bytes)
		}())

	})

	t.Run("Success, drs collateral required amount equal drs collateral amount", func(t *testing.T) {
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

		// get drs account
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)

		// get median price thb
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.New(medianPriceThb.IntPart(), -7), nil)

		// get median price usd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceUsdVeloAddress).
			Return(decimal.New(medianPriceUsd.IntPart(), -7), nil)

		// get median price sgd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceSgdVeloAddress).
			Return(decimal.New(medianPriceSgd.IntPart(), -7), nil)

		// get tp list data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(drsAccountDataEnity.TrustedPartnerListAddress).
			Return(map[string]string{trustedPartnerAddress1: trustedPartnerMetaAddress1}, nil)

		// calculate collateral amount
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(trustedPartnerMetaAddress1).
			Return(map[string]string{fmt.Sprintf("%s_%s", stableCreditAsset1, stableCreditIssuer1): "GCDOC2AYBMEESYXYD3NBPFHWAA44PHQKGTKHRDZXQLWJRWOIW5X2MTFQ"}, nil)

		mockGetIssuerAccountOutput := &entities.GetIssuerAccountOutput{
			Account:        nil,
			PeggedValue:    decimal.NewFromFloat(1.5),
			PeggedCurrency: "USD",
			AssetCode:      stableCreditAsset1,
		}

		helper.mockSubUseCase.EXPECT().GetIssuerAccount(context.Background(), &entities.GetIssuerAccountInput{IssuerAddress: stableCreditIssuer1}).Return(
			mockGetIssuerAccountOutput, nil)

		helper.mockStellarRepo.EXPECT().GetAsset(entities.GetAssetInput{
			AssetCode:   stableCreditAsset1,
			AssetIssuer: stableCreditIssuer1,
		}).Return(&horizon.AssetsPage{
			Links: hal.Links{},
			Embedded: struct {
				Records []horizon.AssetStat
			}{
				Records: []horizon.AssetStat{
					{
						Amount: "1025.0050000",
					},
				},
			},
		}, nil)

		helper.mockStellarRepo.EXPECT().GetAccountBalances(env.DrsPublicKey).Return(
			[]horizon.Balance{
				{
					Balance: drsMediumCollateralAmount.String(),
					Asset: base.Asset{
						Code:   collateralAssetCode,
						Issuer: collateralAssetIssuer,
					},
				},
			}, nil)

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Equal(t, constants.ErrRebalanceIsNotRequired, err.Error())
	})

	t.Run("Success, pegged currency THB", func(t *testing.T) {
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

		// get drs account
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)

		// get median price thb
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.New(medianPriceThb.IntPart(), -7), nil)

		// get median price usd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceUsdVeloAddress).
			Return(decimal.New(medianPriceUsd.IntPart(), -7), nil)

		// get median price sgd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceSgdVeloAddress).
			Return(decimal.New(medianPriceSgd.IntPart(), -7), nil)

		// get tp list data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(drsAccountDataEnity.TrustedPartnerListAddress).
			Return(map[string]string{trustedPartnerAddress1: trustedPartnerMetaAddress1}, nil)

		// calculate collateral amount
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(trustedPartnerMetaAddress1).
			Return(map[string]string{fmt.Sprintf("%s_%s", stableCreditAsset1, stableCreditIssuer1): "GCDOC2AYBMEESYXYD3NBPFHWAA44PHQKGTKHRDZXQLWJRWOIW5X2MTFQ"}, nil)

		mockGetIssuerAccountOutput := &entities.GetIssuerAccountOutput{
			Account:        nil,
			PeggedValue:    decimal.NewFromFloat(1.5),
			PeggedCurrency: "THB",
			AssetCode:      stableCreditAsset1,
		}

		helper.mockSubUseCase.EXPECT().GetIssuerAccount(context.Background(), &entities.GetIssuerAccountInput{IssuerAddress: stableCreditIssuer1}).Return(
			mockGetIssuerAccountOutput, nil)

		helper.mockStellarRepo.EXPECT().GetAsset(entities.GetAssetInput{
			AssetCode:   stableCreditAsset1,
			AssetIssuer: stableCreditIssuer1,
		}).Return(&horizon.AssetsPage{
			Links: hal.Links{},
			Embedded: struct {
				Records []horizon.AssetStat
			}{
				Records: []horizon.AssetStat{
					{
						Amount: "1025.0050000",
					},
				},
			},
		}, nil)

		helper.mockStellarRepo.EXPECT().GetAccountBalances(env.DrsPublicKey).Return(
			[]horizon.Balance{
				{
					Balance: drsHighCollateralAmount.String(),
					Asset: base.Asset{
						Code:   collateralAssetCode,
						Issuer: collateralAssetIssuer,
					},
				},
			}, nil)

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, collateralAssetIssuer, output.Collaterals[0].AssetIssuer)
		assert.Equal(t, collateralAssetCode, output.Collaterals[0].AssetCode)
		assert.Equal(t, drsHighCollateralAmount.String(), output.Collaterals[0].PoolAmount.String())
		assert.True(t, output.Collaterals[0].RequiredAmount.GreaterThan(decimal.Zero))

		stellarTx, txErr := txnbuild.TransactionFromXDR(*output.SignedStellarTxXdr)
		assert.NoError(t, txErr)

		txEnv := stellarTx.TxEnvelope()
		assert.Equal(t, kp1.Address(), txEnv.Tx.SourceAccount.Address())
		assert.Equal(t, xdr.OperationTypePayment, txEnv.Tx.Operations[0].Body.Type)
		assert.Equal(t, env.DrsPublicKey, txEnv.Tx.Operations[0].SourceAccount.Address())
		assert.Equal(t, drsAccountDataEnity.DrsReserve, txEnv.Tx.Operations[0].Body.PaymentOp.Destination.Address())
		assert.Equal(t, collateralAssetIssuer, txEnv.Tx.Operations[0].Body.PaymentOp.Asset.AlphaNum4.Issuer.Address())
		assert.Equal(t, collateralAssetCode, func() string {
			bytes, _ := txEnv.Tx.Operations[0].Body.PaymentOp.Asset.AlphaNum4.AssetCode.MarshalBinary()
			return string(bytes)
		}())

	})

	t.Run("Success, pegged currency SGD", func(t *testing.T) {
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

		// get drs account
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)

		// get median price thb
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.New(medianPriceThb.IntPart(), -7), nil)

		// get median price usd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceUsdVeloAddress).
			Return(decimal.New(medianPriceUsd.IntPart(), -7), nil)

		// get median price sgd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceSgdVeloAddress).
			Return(decimal.New(medianPriceSgd.IntPart(), -7), nil)

		// get tp list data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(drsAccountDataEnity.TrustedPartnerListAddress).
			Return(map[string]string{trustedPartnerAddress1: trustedPartnerMetaAddress1}, nil)

		// calculate collateral amount
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(trustedPartnerMetaAddress1).
			Return(map[string]string{fmt.Sprintf("%s_%s", stableCreditAsset1, stableCreditIssuer1): "GCDOC2AYBMEESYXYD3NBPFHWAA44PHQKGTKHRDZXQLWJRWOIW5X2MTFQ"}, nil)

		mockGetIssuerAccountOutput := &entities.GetIssuerAccountOutput{
			Account:        nil,
			PeggedValue:    decimal.NewFromFloat(1.5),
			PeggedCurrency: "SGD",
			AssetCode:      stableCreditAsset1,
		}

		helper.mockSubUseCase.EXPECT().GetIssuerAccount(context.Background(), &entities.GetIssuerAccountInput{IssuerAddress: stableCreditIssuer1}).Return(
			mockGetIssuerAccountOutput, nil)

		helper.mockStellarRepo.EXPECT().GetAsset(entities.GetAssetInput{
			AssetCode:   stableCreditAsset1,
			AssetIssuer: stableCreditIssuer1,
		}).Return(&horizon.AssetsPage{
			Links: hal.Links{},
			Embedded: struct {
				Records []horizon.AssetStat
			}{
				Records: []horizon.AssetStat{
					{
						Amount: "1025.0050000",
					},
				},
			},
		}, nil)

		helper.mockStellarRepo.EXPECT().GetAccountBalances(env.DrsPublicKey).Return(
			[]horizon.Balance{
				{
					Balance: drsHighCollateralAmount.String(),
					Asset: base.Asset{
						Code:   collateralAssetCode,
						Issuer: collateralAssetIssuer,
					},
				},
			}, nil)

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, collateralAssetIssuer, output.Collaterals[0].AssetIssuer)
		assert.Equal(t, collateralAssetCode, output.Collaterals[0].AssetCode)
		assert.Equal(t, drsHighCollateralAmount.String(), output.Collaterals[0].PoolAmount.String())
		assert.True(t, output.Collaterals[0].RequiredAmount.GreaterThan(decimal.Zero))

		stellarTx, txErr := txnbuild.TransactionFromXDR(*output.SignedStellarTxXdr)
		assert.NoError(t, txErr)

		txEnv := stellarTx.TxEnvelope()
		assert.Equal(t, kp1.Address(), txEnv.Tx.SourceAccount.Address())
		assert.Equal(t, xdr.OperationTypePayment, txEnv.Tx.Operations[0].Body.Type)
		assert.Equal(t, env.DrsPublicKey, txEnv.Tx.Operations[0].SourceAccount.Address())
		assert.Equal(t, drsAccountDataEnity.DrsReserve, txEnv.Tx.Operations[0].Body.PaymentOp.Destination.Address())
		assert.Equal(t, collateralAssetIssuer, txEnv.Tx.Operations[0].Body.PaymentOp.Asset.AlphaNum4.Issuer.Address())
		assert.Equal(t, collateralAssetCode, func() string {
			bytes, _ := txEnv.Tx.Operations[0].Body.PaymentOp.Asset.AlphaNum4.AssetCode.MarshalBinary()
			return string(bytes)
		}())

	})

	t.Run("Success, get empty records of asset", func(t *testing.T) {
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

		// get drs account
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)

		// get median price thb
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.New(medianPriceThb.IntPart(), -7), nil)

		// get median price usd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceUsdVeloAddress).
			Return(decimal.New(medianPriceUsd.IntPart(), -7), nil)

		// get median price sgd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceSgdVeloAddress).
			Return(decimal.New(medianPriceSgd.IntPart(), -7), nil)

		// get tp list data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(drsAccountDataEnity.TrustedPartnerListAddress).
			Return(map[string]string{trustedPartnerAddress1: trustedPartnerMetaAddress1}, nil)

		// calculate collateral amount
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(trustedPartnerMetaAddress1).
			Return(map[string]string{fmt.Sprintf("%s_%s", stableCreditAsset1, stableCreditIssuer1): "GCDOC2AYBMEESYXYD3NBPFHWAA44PHQKGTKHRDZXQLWJRWOIW5X2MTFQ"}, nil)

		mockGetIssuerAccountOutput := &entities.GetIssuerAccountOutput{
			Account:        nil,
			PeggedValue:    decimal.NewFromFloat(1.5),
			PeggedCurrency: "USD",
			AssetCode:      stableCreditAsset1,
		}

		helper.mockSubUseCase.EXPECT().GetIssuerAccount(context.Background(), &entities.GetIssuerAccountInput{IssuerAddress: stableCreditIssuer1}).Return(
			mockGetIssuerAccountOutput, nil)

		helper.mockStellarRepo.EXPECT().GetAsset(entities.GetAssetInput{
			AssetCode:   stableCreditAsset1,
			AssetIssuer: stableCreditIssuer1,
		}).Return(&horizon.AssetsPage{
			Links: hal.Links{},
			Embedded: struct {
				Records []horizon.AssetStat
			}{},
		}, nil)

		helper.mockStellarRepo.EXPECT().GetAccountBalances(env.DrsPublicKey).Return(
			[]horizon.Balance{
				{
					Balance: drsHighCollateralAmount.String(),
					Asset: base.Asset{
						Code:   collateralAssetCode,
						Issuer: collateralAssetIssuer,
					},
				},
			}, nil)

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, collateralAssetIssuer, output.Collaterals[0].AssetIssuer)
		assert.Equal(t, collateralAssetCode, output.Collaterals[0].AssetCode)
		assert.Equal(t, drsHighCollateralAmount.String(), output.Collaterals[0].PoolAmount.String())
		assert.True(t, output.Collaterals[0].RequiredAmount.Equal(decimal.Zero))

		stellarTx, txErr := txnbuild.TransactionFromXDR(*output.SignedStellarTxXdr)
		assert.NoError(t, txErr)

		txEnv := stellarTx.TxEnvelope()
		assert.Equal(t, kp1.Address(), txEnv.Tx.SourceAccount.Address())
		assert.Equal(t, xdr.OperationTypePayment, txEnv.Tx.Operations[0].Body.Type)
		assert.Equal(t, env.DrsPublicKey, txEnv.Tx.Operations[0].SourceAccount.Address())
		assert.Equal(t, drsAccountDataEnity.DrsReserve, txEnv.Tx.Operations[0].Body.PaymentOp.Destination.Address())
		assert.Equal(t, collateralAssetIssuer, txEnv.Tx.Operations[0].Body.PaymentOp.Asset.AlphaNum4.Issuer.Address())
		assert.Equal(t, collateralAssetCode, func() string {
			bytes, _ := txEnv.Tx.Operations[0].Body.PaymentOp.Asset.AlphaNum4.AssetCode.MarshalBinary()
			return string(bytes)
		}())

	})

	t.Run("Error - VeloTx missing signer", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign()

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrUnAuthenticated{}, err)
		assert.Contains(t, err.Error(), constants.ErrSignatureNotFound)
	})

	t.Run("Error - VeloTx wrong signer", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp2)

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrUnAuthenticated{}, err)
		assert.Contains(t, err.Error(), constants.ErrSignatureNotMatchSourceAccount)
	})

	t.Run("Error - can't get tx sender account", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		// get tx sender account
		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(nil, errors.New("stellar return error"))

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrNotFound{}, err)
		assert.Contains(t, err.Error(), constants.ErrGetSenderAccount)
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

		// get drs account
		helper.mockStellarRepo.EXPECT().GetDrsAccountData().
			Return(nil, errors.New("stellar return error"))

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrInternal{}, err)
		assert.Contains(t, err.Error(), constants.ErrGetDrsAccountData)
	})

	t.Run("Error - can't get median Price thb", func(t *testing.T) {
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

		// get drs account
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)

		// get median price thb
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.Zero, errors.New("stellar return error"))

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Contains(t, err.Error(), constants.ErrGetPriceOfPeggedCurrency)
	})

	t.Run("Error - can't get median Price usd", func(t *testing.T) {
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

		// get drs account
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)

		// get median price thb
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.New(medianPriceThb.IntPart(), -7), nil)

		// get median price usd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceUsdVeloAddress).
			Return(decimal.Zero, errors.New("stellar return error"))

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Contains(t, err.Error(), constants.ErrGetPriceOfPeggedCurrency)
	})

	t.Run("Error - can't get median Price sgd", func(t *testing.T) {
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

		// get drs account
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)

		// get median price thb
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.New(medianPriceThb.IntPart(), -7), nil)

		// get median price usd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceUsdVeloAddress).
			Return(decimal.New(medianPriceUsd.IntPart(), -7), nil)

		// get median price sgd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceSgdVeloAddress).
			Return(decimal.Zero, errors.New("stellar return error"))

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Contains(t, err.Error(), constants.ErrGetPriceOfPeggedCurrency)
	})

	t.Run("Error - can't get tp list data", func(t *testing.T) {
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

		// get drs account
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)

		// get median price thb
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.New(medianPriceThb.IntPart(), -7), nil)

		// get median price usd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceUsdVeloAddress).
			Return(decimal.New(medianPriceUsd.IntPart(), -7), nil)

		// get median price sgd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceSgdVeloAddress).
			Return(decimal.New(medianPriceSgd.IntPart(), -7), nil)

		// get tp list data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(drsAccountDataEnity.TrustedPartnerListAddress).
			Return(map[string]string{}, errors.New("cannot decode data"))

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Contains(t, err.Error(), constants.ErrGetTrustedPartnerListAccountData)
	})

	t.Run("Error - can't get tp list meta data", func(t *testing.T) {
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

		// get drs account
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)

		// get median price thb
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.New(medianPriceThb.IntPart(), -7), nil)

		// get median price usd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceUsdVeloAddress).
			Return(decimal.New(medianPriceUsd.IntPart(), -7), nil)

		// get median price sgd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceSgdVeloAddress).
			Return(decimal.New(medianPriceSgd.IntPart(), -7), nil)

		// get tp list data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(drsAccountDataEnity.TrustedPartnerListAddress).
			Return(map[string]string{trustedPartnerAddress1: trustedPartnerMetaAddress1}, nil)

		// calculate collateral amount
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(trustedPartnerMetaAddress1).
			Return(map[string]string{}, errors.New("cannot decode data"))

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Contains(t, err.Error(), constants.ErrGetTrustedPartnerMetaAccountDetail)
	})

	t.Run("Error - can't verify asset code", func(t *testing.T) {
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

		// get drs account
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)

		// get median price thb
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.New(medianPriceThb.IntPart(), -7), nil)

		// get median price usd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceUsdVeloAddress).
			Return(decimal.New(medianPriceUsd.IntPart(), -7), nil)

		// get median price sgd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceSgdVeloAddress).
			Return(decimal.New(medianPriceSgd.IntPart(), -7), nil)

		// get tp list data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(drsAccountDataEnity.TrustedPartnerListAddress).
			Return(map[string]string{trustedPartnerAddress1: trustedPartnerMetaAddress1}, nil)

		// calculate collateral amount
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(trustedPartnerMetaAddress1).
			Return(map[string]string{fmt.Sprintf("%s_%s_WRONGASSET", stableCreditAsset1, stableCreditIssuer1): "GCDOC2AYBMEESYXYD3NBPFHWAA44PHQKGTKHRDZXQLWJRWOIW5X2MTFQ"}, nil)

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Contains(t, err.Error(), constants.ErrVerifyAssetCode)
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

		// get drs account
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)

		// get median price thb
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.New(medianPriceThb.IntPart(), -7), nil)

		// get median price usd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceUsdVeloAddress).
			Return(decimal.New(medianPriceUsd.IntPart(), -7), nil)

		// get median price sgd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceSgdVeloAddress).
			Return(decimal.New(medianPriceSgd.IntPart(), -7), nil)

		// get tp list data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(drsAccountDataEnity.TrustedPartnerListAddress).
			Return(map[string]string{trustedPartnerAddress1: trustedPartnerMetaAddress1}, nil)

		// calculate collateral amount
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(trustedPartnerMetaAddress1).
			Return(map[string]string{fmt.Sprintf("%s_%s", stableCreditAsset1, stableCreditIssuer1): "GCDOC2AYBMEESYXYD3NBPFHWAA44PHQKGTKHRDZXQLWJRWOIW5X2MTFQ"}, nil)

		helper.mockSubUseCase.EXPECT().GetIssuerAccount(context.Background(), &entities.GetIssuerAccountInput{IssuerAddress: stableCreditIssuer1}).Return(
			nil, errors.New("cannot get issuer account"))

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Contains(t, err.Error(), constants.ErrGetIssuerAccount)
	})

	t.Run("Error - can't get asset", func(t *testing.T) {
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

		// get drs account
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)

		// get median price thb
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.New(medianPriceThb.IntPart(), -7), nil)

		// get median price usd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceUsdVeloAddress).
			Return(decimal.New(medianPriceUsd.IntPart(), -7), nil)

		// get median price sgd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceSgdVeloAddress).
			Return(decimal.New(medianPriceSgd.IntPart(), -7), nil)

		// get tp list data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(drsAccountDataEnity.TrustedPartnerListAddress).
			Return(map[string]string{trustedPartnerAddress1: trustedPartnerMetaAddress1}, nil)

		// calculate collateral amount
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(trustedPartnerMetaAddress1).
			Return(map[string]string{fmt.Sprintf("%s_%s", stableCreditAsset1, stableCreditIssuer1): "GCDOC2AYBMEESYXYD3NBPFHWAA44PHQKGTKHRDZXQLWJRWOIW5X2MTFQ"}, nil)

		mockGetIssuerAccountOutput := &entities.GetIssuerAccountOutput{
			Account:        nil,
			PeggedValue:    decimal.NewFromFloat(1.5),
			PeggedCurrency: "USD",
			AssetCode:      stableCreditAsset1,
		}

		helper.mockSubUseCase.EXPECT().GetIssuerAccount(context.Background(), &entities.GetIssuerAccountInput{IssuerAddress: stableCreditIssuer1}).Return(
			mockGetIssuerAccountOutput, nil)

		helper.mockStellarRepo.EXPECT().GetAsset(entities.GetAssetInput{
			AssetCode:   stableCreditAsset1,
			AssetIssuer: stableCreditIssuer1,
		}).Return(nil, errors.New("cannot get asset"))

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Contains(t, err.Error(), fmt.Sprintf(constants.ErrGetAsset, stableCreditAsset1))
	})

	t.Run("Error - invalid stable amount format", func(t *testing.T) {
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

		// get drs account
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)

		// get median price thb
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.New(medianPriceThb.IntPart(), -7), nil)

		// get median price usd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceUsdVeloAddress).
			Return(decimal.New(medianPriceUsd.IntPart(), -7), nil)

		// get median price sgd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceSgdVeloAddress).
			Return(decimal.New(medianPriceSgd.IntPart(), -7), nil)

		// get tp list data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(drsAccountDataEnity.TrustedPartnerListAddress).
			Return(map[string]string{trustedPartnerAddress1: trustedPartnerMetaAddress1}, nil)

		// calculate collateral amount
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(trustedPartnerMetaAddress1).
			Return(map[string]string{fmt.Sprintf("%s_%s", stableCreditAsset1, stableCreditIssuer1): "GCDOC2AYBMEESYXYD3NBPFHWAA44PHQKGTKHRDZXQLWJRWOIW5X2MTFQ"}, nil)

		mockGetIssuerAccountOutput := &entities.GetIssuerAccountOutput{
			Account:        nil,
			PeggedValue:    decimal.NewFromFloat(1.5),
			PeggedCurrency: "USD",
			AssetCode:      stableCreditAsset1,
		}

		helper.mockSubUseCase.EXPECT().GetIssuerAccount(context.Background(), &entities.GetIssuerAccountInput{IssuerAddress: stableCreditIssuer1}).Return(
			mockGetIssuerAccountOutput, nil)

		helper.mockStellarRepo.EXPECT().GetAsset(entities.GetAssetInput{
			AssetCode:   stableCreditAsset1,
			AssetIssuer: stableCreditIssuer1,
		}).Return(&horizon.AssetsPage{
			Links: hal.Links{},
			Embedded: struct {
				Records []horizon.AssetStat
			}{
				Records: []horizon.AssetStat{
					{
						Amount: "amount",
					},
				},
			},
		}, nil)

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Contains(t, err.Error(), "invalid stable amount format")
	})

	t.Run("Error - pegged currency is not support", func(t *testing.T) {
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

		// get drs account
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)

		// get median price thb
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.New(medianPriceThb.IntPart(), -7), nil)

		// get median price usd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceUsdVeloAddress).
			Return(decimal.New(medianPriceUsd.IntPart(), -7), nil)

		// get median price sgd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceSgdVeloAddress).
			Return(decimal.New(medianPriceSgd.IntPart(), -7), nil)

		// get tp list data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(drsAccountDataEnity.TrustedPartnerListAddress).
			Return(map[string]string{trustedPartnerAddress1: trustedPartnerMetaAddress1}, nil)

		// calculate collateral amount
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(trustedPartnerMetaAddress1).
			Return(map[string]string{fmt.Sprintf("%s_%s", stableCreditAsset1, stableCreditIssuer1): "GCDOC2AYBMEESYXYD3NBPFHWAA44PHQKGTKHRDZXQLWJRWOIW5X2MTFQ"}, nil)

		mockGetIssuerAccountOutput := &entities.GetIssuerAccountOutput{
			Account:        nil,
			PeggedValue:    decimal.NewFromFloat(1.5),
			PeggedCurrency: "JPY",
			AssetCode:      stableCreditAsset1,
		}

		helper.mockSubUseCase.EXPECT().GetIssuerAccount(context.Background(), &entities.GetIssuerAccountInput{IssuerAddress: stableCreditIssuer1}).Return(
			mockGetIssuerAccountOutput, nil)

		helper.mockStellarRepo.EXPECT().GetAsset(entities.GetAssetInput{
			AssetCode:   stableCreditAsset1,
			AssetIssuer: stableCreditIssuer1,
		}).Return(&horizon.AssetsPage{
			Links: hal.Links{},
			Embedded: struct {
				Records []horizon.AssetStat
			}{
				Records: []horizon.AssetStat{
					{
						Amount: "1025.0050000",
					},
				},
			},
		}, nil)

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrInternal{}, err)
		assert.Equal(t, constants.ErrPeggedCurrencyIsNotSupport, err.Error())
	})

	t.Run("Error, can't get drs collateral account balances", func(t *testing.T) {
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

		// get drs account
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)

		// get median price thb
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.New(medianPriceThb.IntPart(), -7), nil)

		// get median price usd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceUsdVeloAddress).
			Return(decimal.New(medianPriceUsd.IntPart(), -7), nil)

		// get median price sgd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceSgdVeloAddress).
			Return(decimal.New(medianPriceSgd.IntPart(), -7), nil)

		// get tp list data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(drsAccountDataEnity.TrustedPartnerListAddress).
			Return(map[string]string{trustedPartnerAddress1: trustedPartnerMetaAddress1}, nil)

		// calculate collateral amount
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(trustedPartnerMetaAddress1).
			Return(map[string]string{fmt.Sprintf("%s_%s", stableCreditAsset1, stableCreditIssuer1): "GCDOC2AYBMEESYXYD3NBPFHWAA44PHQKGTKHRDZXQLWJRWOIW5X2MTFQ"}, nil)

		mockGetIssuerAccountOutput := &entities.GetIssuerAccountOutput{
			Account:        nil,
			PeggedValue:    decimal.NewFromFloat(1.5),
			PeggedCurrency: "USD",
			AssetCode:      stableCreditAsset1,
		}

		helper.mockSubUseCase.EXPECT().GetIssuerAccount(context.Background(), &entities.GetIssuerAccountInput{IssuerAddress: stableCreditIssuer1}).Return(
			mockGetIssuerAccountOutput, nil)

		helper.mockStellarRepo.EXPECT().GetAsset(entities.GetAssetInput{
			AssetCode:   stableCreditAsset1,
			AssetIssuer: stableCreditIssuer1,
		}).Return(&horizon.AssetsPage{
			Links: hal.Links{},
			Embedded: struct {
				Records []horizon.AssetStat
			}{
				Records: []horizon.AssetStat{
					{
						Amount: "1025.0050000",
					},
				},
			},
		}, nil)

		helper.mockStellarRepo.EXPECT().GetAccountBalances(env.DrsPublicKey).Return(
			nil, errors.New("stellar return error"))

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrInternal{}, err)
		assert.Contains(t, err.Error(), constants.ErrGetDrsAccountBalance)

	})

	t.Run("Success, drs collateral amount greater than drs collateral required amount", func(t *testing.T) {
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

		// get drs account
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)

		// get median price thb
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.New(medianPriceThb.IntPart(), -7), nil)

		// get median price usd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceUsdVeloAddress).
			Return(decimal.New(medianPriceUsd.IntPart(), -7), nil)

		// get median price sgd
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceSgdVeloAddress).
			Return(decimal.New(medianPriceSgd.IntPart(), -7), nil)

		// get tp list data
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(drsAccountDataEnity.TrustedPartnerListAddress).
			Return(map[string]string{trustedPartnerAddress1: trustedPartnerMetaAddress1}, nil)

		// calculate collateral amount
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(trustedPartnerMetaAddress1).
			Return(map[string]string{fmt.Sprintf("%s_%s", stableCreditAsset1, stableCreditIssuer1): "GCDOC2AYBMEESYXYD3NBPFHWAA44PHQKGTKHRDZXQLWJRWOIW5X2MTFQ"}, nil)

		mockGetIssuerAccountOutput := &entities.GetIssuerAccountOutput{
			Account:        nil,
			PeggedValue:    decimal.NewFromFloat(1.5),
			PeggedCurrency: "USD",
			AssetCode:      stableCreditAsset1,
		}

		helper.mockSubUseCase.EXPECT().GetIssuerAccount(context.Background(), &entities.GetIssuerAccountInput{IssuerAddress: stableCreditIssuer1}).Return(
			mockGetIssuerAccountOutput, nil)

		helper.mockStellarRepo.EXPECT().GetAsset(entities.GetAssetInput{
			AssetCode:   stableCreditAsset1,
			AssetIssuer: stableCreditIssuer1,
		}).Return(&horizon.AssetsPage{
			Links: hal.Links{},
			Embedded: struct {
				Records []horizon.AssetStat
			}{
				Records: []horizon.AssetStat{
					{
						Amount: "1025.0050000",
					},
				},
			},
		}, nil)

		helper.mockStellarRepo.EXPECT().GetAccountBalances(env.DrsPublicKey).Return(
			[]horizon.Balance{
				{
					Balance: drsHighCollateralAmount.String(),
					Asset: base.Asset{
						Code:   "kTONG",
						Issuer: collateralAssetIssuer,
					},
				},
			}, nil)

		output, err := helper.useCase.RebalanceReserve(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Equal(t, constants.ErrDrsCollateralTrustlineNotFound, err.Error())

	})

}

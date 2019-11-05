package usecases_test

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/protocols/horizon/base"
	"github.com/stellar/go/support/render/hal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	env "gitlab.com/velo-labs/cen/node/app/environments"
	"gitlab.com/velo-labs/cen/node/app/errors"
	"gitlab.com/velo-labs/cen/node/app/testhelpers"
	"testing"
)

func TestUseCase_GetCollateralHealthCheck(t *testing.T) {
	testhelpers.InitEnv()

	var (
		trustedPartnerAddress1     = publicKey2
		trustedPartnerMetaAddress1 = publicKey3
		trustedPartnerAddress2     = publicKey4
		trustedPartnerMetaAddress2 = publicKey5

		medianPriceThb = decimal.NewFromFloat(23000000)
		medianPriceUsd = decimal.NewFromFloat(10000000)
		medianPriceSgd = decimal.NewFromFloat(20000000)

		stableCreditAsset1  = "vUSD"
		stableCreditIssuer1 = publicKey1
		stableCreditAsset2  = "vTHB"
		stableCreditIssuer2 = publicKey1

		collateralPoolAmount  = decimal.NewFromFloat(2050.125)
		collateralAssetCode   = string(vxdr.AssetVELO)
		collateralAssetIssuer = env.VeloIssuerPublicKey
	)

	t.Run("success, pegged currency USD", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

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

		helper.mockStellarRepo.EXPECT().GetAccountBalances(env.DrsPublicKey).Return([]horizon.Balance{
			{
				Balance: collateralPoolAmount.String(),
				Asset: base.Asset{
					Code:   collateralAssetCode,
					Issuer: collateralAssetIssuer,
				},
			},
		}, nil)

		output, err := helper.useCase.GetCollateralHealthCheck(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, collateralAssetCode, output.AssetCode)
		assert.Equal(t, collateralAssetIssuer, output.AssetIssuer)
		assert.Equal(t, collateralPoolAmount.String(), output.PoolAmount.String())
		assert.True(t, output.RequiredAmount.GreaterThan(decimal.Zero))
	})

	t.Run("success, pegged currency THB", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

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

		helper.mockStellarRepo.EXPECT().GetAccountBalances(env.DrsPublicKey).Return([]horizon.Balance{
			{
				Balance: collateralPoolAmount.String(),
				Asset: base.Asset{
					Code:   collateralAssetCode,
					Issuer: collateralAssetIssuer,
				},
			},
		}, nil)

		output, err := helper.useCase.GetCollateralHealthCheck(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, collateralAssetCode, output.AssetCode)
		assert.Equal(t, collateralAssetIssuer, output.AssetIssuer)
		assert.Equal(t, collateralPoolAmount.String(), output.PoolAmount.String())
		assert.True(t, output.RequiredAmount.GreaterThan(decimal.Zero))
	})

	t.Run("success, pegged currency SGD", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

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

		helper.mockStellarRepo.EXPECT().GetAccountBalances(env.DrsPublicKey).Return([]horizon.Balance{
			{
				Balance: collateralPoolAmount.String(),
				Asset: base.Asset{
					Code:   collateralAssetCode,
					Issuer: collateralAssetIssuer,
				},
			},
		}, nil)

		output, err := helper.useCase.GetCollateralHealthCheck(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, collateralAssetCode, output.AssetCode)
		assert.Equal(t, collateralAssetIssuer, output.AssetIssuer)
		assert.Equal(t, collateralPoolAmount.String(), output.PoolAmount.String())
		assert.True(t, output.RequiredAmount.GreaterThan(decimal.Zero))
	})

	t.Run("success, multi stable credit in one tp", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

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
			Return(map[string]string{
				fmt.Sprintf("%s_%s", stableCreditAsset1, stableCreditIssuer1): "GCDOC2AYBMEESYXYD3NBPFHWAA44PHQKGTKHRDZXQLWJRWOIW5X2MTFQ",
				fmt.Sprintf("%s_%s", stableCreditAsset2, stableCreditIssuer2): "GCDOC2AYBMEESYXYD3NBPFHWAA44PHQKGTKHRDZXQLWJRWOIW5X2MTFQ",
			}, nil)

		mockGetIssuerAccountOutput1 := &entities.GetIssuerAccountOutput{
			Account:        nil,
			PeggedValue:    decimal.NewFromFloat(1.5),
			PeggedCurrency: "USD",
			AssetCode:      stableCreditAsset1,
		}

		helper.mockSubUseCase.EXPECT().GetIssuerAccount(context.Background(), &entities.GetIssuerAccountInput{IssuerAddress: stableCreditIssuer1}).Return(
			mockGetIssuerAccountOutput1, nil)

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

		mockGetIssuerAccountOutput2 := &entities.GetIssuerAccountOutput{
			Account:        nil,
			PeggedValue:    decimal.NewFromFloat(1.5),
			PeggedCurrency: "THB",
			AssetCode:      stableCreditAsset2,
		}

		helper.mockSubUseCase.EXPECT().GetIssuerAccount(context.Background(), &entities.GetIssuerAccountInput{IssuerAddress: stableCreditIssuer2}).Return(
			mockGetIssuerAccountOutput2, nil)

		helper.mockStellarRepo.EXPECT().GetAsset(entities.GetAssetInput{
			AssetCode:   stableCreditAsset2,
			AssetIssuer: stableCreditIssuer2,
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

		helper.mockStellarRepo.EXPECT().GetAccountBalances(env.DrsPublicKey).Return([]horizon.Balance{
			{
				Balance: collateralPoolAmount.String(),
				Asset: base.Asset{
					Code:   collateralAssetCode,
					Issuer: collateralAssetIssuer,
				},
			},
		}, nil)

		output, err := helper.useCase.GetCollateralHealthCheck(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, collateralAssetCode, output.AssetCode)
		assert.Equal(t, collateralAssetIssuer, output.AssetIssuer)
		assert.Equal(t, collateralPoolAmount.String(), output.PoolAmount.String())
		assert.True(t, output.RequiredAmount.GreaterThan(decimal.Zero))
	})

	t.Run("success, multi tp", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

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
			Return(
				map[string]string{
					trustedPartnerAddress1: trustedPartnerMetaAddress1,
					trustedPartnerAddress2: trustedPartnerMetaAddress2,
				}, nil)

		// calculate collateral amount
		// tp1
		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(trustedPartnerMetaAddress1).
			Return(map[string]string{fmt.Sprintf("%s_%s", stableCreditAsset1, stableCreditIssuer1): "GCDOC2AYBMEESYXYD3NBPFHWAA44PHQKGTKHRDZXQLWJRWOIW5X2MTFQ"}, nil)

		mockGetIssuerAccountOutput1 := &entities.GetIssuerAccountOutput{
			Account:        nil,
			PeggedValue:    decimal.NewFromFloat(1.5),
			PeggedCurrency: "USD",
			AssetCode:      stableCreditAsset1,
		}

		helper.mockSubUseCase.EXPECT().GetIssuerAccount(context.Background(), &entities.GetIssuerAccountInput{IssuerAddress: stableCreditIssuer1}).Return(
			mockGetIssuerAccountOutput1, nil)

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

		// tp 2

		helper.mockStellarRepo.EXPECT().GetAccountDecodedData(trustedPartnerMetaAddress2).
			Return(map[string]string{fmt.Sprintf("%s_%s", stableCreditAsset2, stableCreditIssuer2): "GCDOC2AYBMEESYXYD3NBPFHWAA44PHQKGTKHRDZXQLWJRWOIW5X2MTFQ"}, nil)

		mockGetIssuerAccountOutput2 := &entities.GetIssuerAccountOutput{
			Account:        nil,
			PeggedValue:    decimal.NewFromFloat(1.5),
			PeggedCurrency: "THB",
			AssetCode:      stableCreditAsset2,
		}

		helper.mockSubUseCase.EXPECT().GetIssuerAccount(context.Background(), &entities.GetIssuerAccountInput{IssuerAddress: stableCreditIssuer2}).Return(
			mockGetIssuerAccountOutput2, nil)

		helper.mockStellarRepo.EXPECT().GetAsset(entities.GetAssetInput{
			AssetCode:   stableCreditAsset2,
			AssetIssuer: stableCreditIssuer2,
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

		helper.mockStellarRepo.EXPECT().GetAccountBalances(env.DrsPublicKey).Return([]horizon.Balance{
			{
				Balance: collateralPoolAmount.String(),
				Asset: base.Asset{
					Code:   collateralAssetCode,
					Issuer: collateralAssetIssuer,
				},
			},
		}, nil)

		output, err := helper.useCase.GetCollateralHealthCheck(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, collateralAssetCode, output.AssetCode)
		assert.Equal(t, collateralAssetIssuer, output.AssetIssuer)
		assert.Equal(t, collateralPoolAmount.String(), output.PoolAmount.String())
		assert.True(t, output.RequiredAmount.GreaterThan(decimal.Zero))
	})

	t.Run("Error - can't get DRS account data", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		// get drs account
		helper.mockStellarRepo.EXPECT().GetDrsAccountData().
			Return(nil, errors.New("stellar return error"))

		output, err := helper.useCase.GetCollateralHealthCheck(context.Background())

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrInternal{}, err)
		assert.Contains(t, err.Error(), constants.ErrGetDrsAccountData)
	})

	t.Run("Error - can't get median Price thb", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		// get drs account
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)

		// get median price thb
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.Zero, errors.New("stellar return error"))

		output, err := helper.useCase.GetCollateralHealthCheck(context.Background())

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Contains(t, err.Error(), constants.ErrGetPriceOfPeggedCurrency)
	})

	t.Run("Error - can't get median Price usd", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

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

		output, err := helper.useCase.GetCollateralHealthCheck(context.Background())

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Contains(t, err.Error(), constants.ErrGetPriceOfPeggedCurrency)
	})

	t.Run("Error - can't get median Price sgd", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

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

		output, err := helper.useCase.GetCollateralHealthCheck(context.Background())

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Contains(t, err.Error(), constants.ErrGetPriceOfPeggedCurrency)
	})

	t.Run("Error - can't get tp list data", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

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

		output, err := helper.useCase.GetCollateralHealthCheck(context.Background())

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Contains(t, err.Error(), constants.ErrGetTrustedPartnerListAccountData)
	})

	t.Run("Error - can't get tp list meta data", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

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

		output, err := helper.useCase.GetCollateralHealthCheck(context.Background())

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Contains(t, err.Error(), constants.ErrGetTrustedPartnerMetaAccountDetail)
	})

	t.Run("Error - can't verify asset code", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

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

		output, err := helper.useCase.GetCollateralHealthCheck(context.Background())

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Contains(t, err.Error(), constants.ErrVerifyAssetCode)
	})

	t.Run("Error - can't get issuer account", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

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

		output, err := helper.useCase.GetCollateralHealthCheck(context.Background())

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Contains(t, err.Error(), constants.ErrGetIssuerAccount)
	})

	t.Run("Error - can't get asset empty records", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

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

		output, err := helper.useCase.GetCollateralHealthCheck(context.Background())

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Contains(t, err.Error(), fmt.Sprintf(constants.ErrGetAsset, stableCreditAsset1))
	})

	t.Run("Error - get asset", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

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
				Records: []horizon.AssetStat{},
			},
		}, nil)

		output, err := helper.useCase.GetCollateralHealthCheck(context.Background())

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Equal(t, fmt.Sprintf(constants.ErrGetAsset, stableCreditAsset1), err.Error())
	})

	t.Run("Error - invalid stable amount format", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

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

		output, err := helper.useCase.GetCollateralHealthCheck(context.Background())

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
		assert.Contains(t, err.Error(), "invalid stable amount format")
	})

	t.Run("Error - can,t get drs collateral balances", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

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
						Amount: "100.0",
					},
				},
			},
		}, nil)

		helper.mockStellarRepo.EXPECT().GetAccountBalances(env.DrsPublicKey).Return(nil, errors.New("can't get drs reserve "))

		output, err := helper.useCase.GetCollateralHealthCheck(context.Background())

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrInternal{}, err)
		assert.Contains(t, err.Error(), constants.ErrGetDrsAccountBalance)
	})

	t.Run("Error - invalid drs collateral pool amount format", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

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
						Amount: "100.0",
					},
				},
			},
		}, nil)

		helper.mockStellarRepo.EXPECT().GetAccountBalances(env.DrsPublicKey).Return(
			[]horizon.Balance{
				{
					Balance: "amount",
					Asset: base.Asset{
						Code:   collateralAssetCode,
						Issuer: collateralAssetIssuer,
					},
				},
			}, nil)

		output, err := helper.useCase.GetCollateralHealthCheck(context.Background())

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.IsType(t, nerrors.ErrInternal{}, err)
	})
}

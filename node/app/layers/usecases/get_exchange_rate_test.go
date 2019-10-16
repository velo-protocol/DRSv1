package usecases_test

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	env "gitlab.com/velo-labs/cen/node/app/environments"
	"gitlab.com/velo-labs/cen/node/app/errors"
	"testing"
)

func TestUseCase_GetExchangeRate(t *testing.T) {
	var (
		peggedValue = decimal.NewFromFloat(15000000)
		medianPrice = decimal.NewFromFloat(23000000)

		vThbIssuerAddress = "GAN6D232HXTF4OHL7J36SAJD3M22H26B2O4QFVRO32OEM523KTMB6Q72"
		vThbAsset         = "vTHB"

		trustedPartnerAddress     = publicKey2
		trustedPartnerMetaAddress = publicKey3
		peggedCurrency            = "THB"

		getMockGetExchangeRateInput = &entities.GetExchangeRateInput{
			AssetCode: vThbAsset,
			Issuer:    vThbIssuerAddress,
		}

		getMockGetIssuerAccount = &entities.GetIssuerAccountInput{
			IssuerAddress: vThbIssuerAddress,
		}

		getMockGetTrustedPartnerFromIssuerAccountInput = &entities.GetTrustedPartnerFromIssuerAccountInput{
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
		}
	)

	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockSubUseCase.EXPECT().
			GetIssuerAccount(context.Background(), getMockGetIssuerAccount).
			Return(&entities.GetIssuerAccountOutput{
				Account: &horizon.Account{
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
				PeggedValue:    decimal.New(peggedValue.IntPart(), -7),
				PeggedCurrency: peggedCurrency,
				AssetCode:      vThbAsset,
			}, nil)
		helper.mockSubUseCase.EXPECT().
			GetTrustedPartnerFromIssuerAccount(context.Background(), gomock.AssignableToTypeOf(getMockGetTrustedPartnerFromIssuerAccountInput)).
			Return(&entities.GetTrustedPartnerFromIssuerAccountOutput{
				TrustedPartnerAccount: &horizon.Account{
					AccountID: trustedPartnerAddress,
				},
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountDecodedDataByKey(drsAccountDataEnity.TrustedPartnerListAddress, trustedPartnerAddress).
			Return(trustedPartnerMetaAddress, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountDecodedDataByKey(trustedPartnerMetaAddress, fmt.Sprintf("%s_%s", vThbAsset, vThbIssuerAddress)).
			Return(fmt.Sprintf("%s_%s", vThbAsset, vThbIssuerAddress), nil)
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.New(medianPrice.IntPart(), -7), nil)

		output, err := helper.useCase.GetExchangeRate(context.Background(), getMockGetExchangeRateInput)

		assert.NoError(t, err)
		assert.NotEmpty(t, output)
		assert.Equal(t, vThbIssuerAddress, output.Issuer)
		assert.Equal(t, vThbAsset, output.AssetCode)
		assert.Equal(t, "0.6521739", output.RedeemablePricePerUnit.String())
	})

	t.Run("Error - get exchange rate input invalid argument ", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		_, err := helper.useCase.GetExchangeRate(context.Background(), &entities.GetExchangeRateInput{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrMustNotBeBlank)
		assert.IsType(t, nerrors.ErrInvalidArgument{}, err)
	})

	t.Run("Error - fail to get issuer account", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockSubUseCase.EXPECT().
			GetIssuerAccount(context.Background(), getMockGetIssuerAccount).
			Return(nil, errors.New("fail to get issuer account"))

		_, err := helper.useCase.GetExchangeRate(context.Background(), getMockGetExchangeRateInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrGetIssuerAccount)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
	})

	t.Run("Error - fail to get trusted partner from issuer account", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockSubUseCase.EXPECT().
			GetIssuerAccount(context.Background(), getMockGetIssuerAccount).
			Return(&entities.GetIssuerAccountOutput{
				Account: &horizon.Account{
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
				PeggedValue:    decimal.New(peggedValue.IntPart(), -7),
				PeggedCurrency: peggedCurrency,
				AssetCode:      vThbAsset,
			}, nil)
		helper.mockSubUseCase.EXPECT().
			GetTrustedPartnerFromIssuerAccount(context.Background(), gomock.AssignableToTypeOf(getMockGetTrustedPartnerFromIssuerAccountInput)).
			Return(&entities.GetTrustedPartnerFromIssuerAccountOutput{
				TrustedPartnerAccount: &horizon.Account{
					AccountID: trustedPartnerAddress,
				},
			}, errors.New("some error has occurs"))

		_, err := helper.useCase.GetExchangeRate(context.Background(), getMockGetExchangeRateInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "some error has occurs")
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
	})

	t.Run("Error - fail to get drs account", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockSubUseCase.EXPECT().
			GetIssuerAccount(context.Background(), getMockGetIssuerAccount).
			Return(&entities.GetIssuerAccountOutput{
				Account: &horizon.Account{
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
				PeggedValue:    decimal.New(peggedValue.IntPart(), -7),
				PeggedCurrency: peggedCurrency,
				AssetCode:      vThbAsset,
			}, nil)
		helper.mockSubUseCase.EXPECT().
			GetTrustedPartnerFromIssuerAccount(context.Background(), gomock.AssignableToTypeOf(getMockGetTrustedPartnerFromIssuerAccountInput)).
			Return(&entities.GetTrustedPartnerFromIssuerAccountOutput{
				TrustedPartnerAccount: &horizon.Account{
					AccountID: trustedPartnerAddress,
				},
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(nil, errors.New("some error has occurred"))

		_, err := helper.useCase.GetExchangeRate(context.Background(), getMockGetExchangeRateInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrGetDrsAccountData)
		assert.IsType(t, nerrors.ErrInternal{}, err)
	})

	t.Run("Error - fail to verify trusted partner account", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockSubUseCase.EXPECT().
			GetIssuerAccount(context.Background(), getMockGetIssuerAccount).
			Return(&entities.GetIssuerAccountOutput{
				Account: &horizon.Account{
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
				PeggedValue:    decimal.New(peggedValue.IntPart(), -7),
				PeggedCurrency: peggedCurrency,
				AssetCode:      vThbAsset,
			}, nil)
		helper.mockSubUseCase.EXPECT().
			GetTrustedPartnerFromIssuerAccount(context.Background(), gomock.AssignableToTypeOf(getMockGetTrustedPartnerFromIssuerAccountInput)).
			Return(&entities.GetTrustedPartnerFromIssuerAccountOutput{
				TrustedPartnerAccount: &horizon.Account{
					AccountID: trustedPartnerAddress,
				},
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountDecodedDataByKey(drsAccountDataEnity.TrustedPartnerListAddress, trustedPartnerAddress).
			Return("", errors.New("some error has occurred"))

		_, err := helper.useCase.GetExchangeRate(context.Background(), getMockGetExchangeRateInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrVerifyTrustedPartnerAccount)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
	})

	t.Run("Error - fail to verify asset code", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockSubUseCase.EXPECT().
			GetIssuerAccount(context.Background(), getMockGetIssuerAccount).
			Return(&entities.GetIssuerAccountOutput{
				Account: &horizon.Account{
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
				PeggedValue:    decimal.New(peggedValue.IntPart(), -7),
				PeggedCurrency: peggedCurrency,
				AssetCode:      vThbAsset,
			}, nil)
		helper.mockSubUseCase.EXPECT().
			GetTrustedPartnerFromIssuerAccount(context.Background(), gomock.AssignableToTypeOf(getMockGetTrustedPartnerFromIssuerAccountInput)).
			Return(&entities.GetTrustedPartnerFromIssuerAccountOutput{
				TrustedPartnerAccount: &horizon.Account{
					AccountID: trustedPartnerAddress,
				},
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountDecodedDataByKey(drsAccountDataEnity.TrustedPartnerListAddress, trustedPartnerAddress).
			Return(trustedPartnerMetaAddress, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountDecodedDataByKey(trustedPartnerMetaAddress, fmt.Sprintf("%s_%s", vThbAsset, vThbIssuerAddress)).
			Return("", errors.New("some error has occurred"))

		_, err := helper.useCase.GetExchangeRate(context.Background(), getMockGetExchangeRateInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrVerifyAssetCode)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
	})

	t.Run("Error - fail to get price of pegged currency", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockSubUseCase.EXPECT().
			GetIssuerAccount(context.Background(), getMockGetIssuerAccount).
			Return(&entities.GetIssuerAccountOutput{
				Account: &horizon.Account{
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
				PeggedValue:    decimal.New(peggedValue.IntPart(), -7),
				PeggedCurrency: peggedCurrency,
				AssetCode:      vThbAsset,
			}, nil)
		helper.mockSubUseCase.EXPECT().
			GetTrustedPartnerFromIssuerAccount(context.Background(), gomock.AssignableToTypeOf(getMockGetTrustedPartnerFromIssuerAccountInput)).
			Return(&entities.GetTrustedPartnerFromIssuerAccountOutput{
				TrustedPartnerAccount: &horizon.Account{
					AccountID: trustedPartnerAddress,
				},
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountDecodedDataByKey(drsAccountDataEnity.TrustedPartnerListAddress, trustedPartnerAddress).
			Return(trustedPartnerMetaAddress, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountDecodedDataByKey(trustedPartnerMetaAddress, fmt.Sprintf("%s_%s", vThbAsset, vThbIssuerAddress)).
			Return(fmt.Sprintf("%s_%s", vThbAsset, vThbIssuerAddress), nil)
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.Zero, errors.New("some error has occurred"))

		_, err := helper.useCase.GetExchangeRate(context.Background(), getMockGetExchangeRateInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrGetPriceOfPeggedCurrency)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
	})

	t.Run("Error - median price must be greater than zero", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockSubUseCase.EXPECT().
			GetIssuerAccount(context.Background(), getMockGetIssuerAccount).
			Return(&entities.GetIssuerAccountOutput{
				Account: &horizon.Account{
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
				PeggedValue:    decimal.New(peggedValue.IntPart(), -7),
				PeggedCurrency: peggedCurrency,
				AssetCode:      vThbAsset,
			}, nil)
		helper.mockSubUseCase.EXPECT().
			GetTrustedPartnerFromIssuerAccount(context.Background(), gomock.AssignableToTypeOf(getMockGetTrustedPartnerFromIssuerAccountInput)).
			Return(&entities.GetTrustedPartnerFromIssuerAccountOutput{
				TrustedPartnerAccount: &horizon.Account{
					AccountID: trustedPartnerAddress,
				},
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountDecodedDataByKey(drsAccountDataEnity.TrustedPartnerListAddress, trustedPartnerAddress).
			Return(trustedPartnerMetaAddress, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountDecodedDataByKey(trustedPartnerMetaAddress, fmt.Sprintf("%s_%s", vThbAsset, vThbIssuerAddress)).
			Return(fmt.Sprintf("%s_%s", vThbAsset, vThbIssuerAddress), nil)
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.NewFromFloat(-1.0), nil)

		_, err := helper.useCase.GetExchangeRate(context.Background(), getMockGetExchangeRateInput)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrMedianPriceMustBeGreaterThanZero)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
	})

}

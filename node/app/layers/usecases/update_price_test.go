package usecases_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/assert"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"github.com/velo-protocol/DRSv1/node/app/constants"
	"github.com/velo-protocol/DRSv1/node/app/errors"
	"testing"
)

func TestUseCase_UpdatePrice(t *testing.T) {
	var (
		asset                       = "VELO"
		currency                    = "THB"
		priceInCurrencyPerAssetUnit = "1.5"
	)

	t.Run("Success", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.PriceUpdate{
				Asset:                       asset,
				Currency:                    currency,
				PriceInCurrencyPerAssetUnit: priceInCurrencyPerAssetUnit,
			},
		}
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{AccountID: publicKey1, Sequence: "1"}, nil)
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountData(drsAccountDataEnity.PriceFeederListAddress).
			Return(map[string]string{
				publicKey1: base64.StdEncoding.EncodeToString([]byte(currency)),
			}, nil)

		output, err := helper.useCase.UpdatePrice(context.Background(), veloTx)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.NotNil(t, output.SignedStellarTxXdr)
		assert.Equal(t, asset, output.Asset)
		assert.Equal(t, currency, output.Currency)
		assert.Equal(t, "1.5000000", output.PriceInCurrencyPerAssetUnit)
	})
	t.Run("Error - velo op validation fail", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.PriceUpdate{
				Asset:                       "",
				Currency:                    currency,
				PriceInCurrencyPerAssetUnit: priceInCurrencyPerAssetUnit,
			},
		}

		_, err := helper.useCase.UpdatePrice(context.Background(), veloTx)
		assert.Error(t, err)
		assert.IsType(t, nerrors.ErrInvalidArgument{}, err)
	})

	t.Run("Error - VeloTx missing signer", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.PriceUpdate{
				Asset:                       asset,
				Currency:                    currency,
				PriceInCurrencyPerAssetUnit: priceInCurrencyPerAssetUnit,
			},
		}
		_ = veloTx.Build()
		_ = veloTx.Sign()

		_, err := helper.useCase.UpdatePrice(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrSignatureNotFound)
		assert.IsType(t, nerrors.ErrUnAuthenticated{}, err)
	})

	t.Run("Error - VeloTx wrong signer", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.PriceUpdate{
				Asset:                       asset,
				Currency:                    currency,
				PriceInCurrencyPerAssetUnit: priceInCurrencyPerAssetUnit,
			},
		}
		_ = veloTx.Build()
		_ = veloTx.Sign(kp2)

		_, err := helper.useCase.UpdatePrice(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrSignatureNotMatchSourceAccount)
		assert.IsType(t, nerrors.ErrUnAuthenticated{}, err)
	})
	t.Run("Error - tx sender account not found", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.PriceUpdate{
				Asset:                       asset,
				Currency:                    currency,
				PriceInCurrencyPerAssetUnit: priceInCurrencyPerAssetUnit,
			},
		}
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(nil, errors.New("some error has occurred"))

		_, err := helper.useCase.UpdatePrice(context.Background(), veloTx)

		assert.IsType(t, nerrors.ErrNotFound{}, err)
	})

	t.Run("Error - fail to get drs account data", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.PriceUpdate{
				Asset:                       asset,
				Currency:                    currency,
				PriceInCurrencyPerAssetUnit: priceInCurrencyPerAssetUnit,
			},
		}
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{AccountID: publicKey1, Sequence: "1"}, nil)
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(nil, errors.New("some error has occurred"))

		_, err := helper.useCase.UpdatePrice(context.Background(), veloTx)

		assert.IsType(t, nerrors.ErrInternal{}, err)
		assert.Contains(t, err.Error(), constants.ErrGetDrsAccountData)
	})
	t.Run("Error - fail to get data of price feeder list account", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.PriceUpdate{
				Asset:                       asset,
				Currency:                    currency,
				PriceInCurrencyPerAssetUnit: priceInCurrencyPerAssetUnit,
			},
		}
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{AccountID: publicKey1, Sequence: "1"}, nil)
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountData(drsAccountDataEnity.PriceFeederListAddress).
			Return(nil, errors.New("some error has occurred"))

		_, err := helper.useCase.UpdatePrice(context.Background(), veloTx)
		assert.IsType(t, nerrors.ErrInternal{}, err)
		assert.Contains(t, err.Error(), constants.ErrGetPriceFeederListAccountData)
	})
	t.Run("Error - tx sender has no permission to perform price update", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.PriceUpdate{
				Asset:                       asset,
				Currency:                    currency,
				PriceInCurrencyPerAssetUnit: priceInCurrencyPerAssetUnit,
			},
		}
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{AccountID: publicKey1, Sequence: "1"}, nil)
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountData(drsAccountDataEnity.PriceFeederListAddress).
			Return(map[string]string{}, nil)

		_, err := helper.useCase.UpdatePrice(context.Background(), veloTx)

		assert.EqualError(t, err, fmt.Sprintf(constants.ErrFormatSignerNotHavePermission, constants.VeloOpPriceUpdate))
		assert.IsType(t, nerrors.ErrPermissionDenied{}, err)
	})
	t.Run("Error - currency must match with the registered currency", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.PriceUpdate{
				Asset:                       asset,
				Currency:                    currency,
				PriceInCurrencyPerAssetUnit: priceInCurrencyPerAssetUnit,
			},
		}
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{AccountID: publicKey1, Sequence: "1"}, nil)
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountData(drsAccountDataEnity.PriceFeederListAddress).
			Return(map[string]string{
				publicKey1: base64.StdEncoding.EncodeToString([]byte("USD")),
			}, nil)

		_, err := helper.useCase.UpdatePrice(context.Background(), veloTx)

		assert.IsType(t, nerrors.ErrInvalidArgument{}, err)
		assert.Contains(t, err.Error(), constants.ErrCurrencyMustMatchWithRegisteredCurrency)
	})
}

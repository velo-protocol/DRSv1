package grpc_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/assert"
	spec "github.com/velo-protocol/DRSv1/grpc"
	"github.com/velo-protocol/DRSv1/libs/convert"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"github.com/velo-protocol/DRSv1/node/app/constants"
	"github.com/velo-protocol/DRSv1/node/app/entities"
	"github.com/velo-protocol/DRSv1/node/app/errors"
	"testing"
)

func TestHandler_SubmitVeloTx_PriceUpdate(t *testing.T) {
	var (
		priceFeederKP, _            = vconvert.SecretKeyToKeyPair(secretKey1)
		asset                       = "VELO"
		currency                    = "THB"
		priceInCurrencyPerAssetUnit = "1.5000000"
	)

	mockedOutput := &entities.UpdatePriceOutput{
		SignedStellarTxXdr:          "AAAAA...=",
		Asset:                       asset,
		Currency:                    currency,
		PriceInCurrencyPerAssetUnit: decimal.NewFromFloat(1.5000000),
	}

	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		helper.mockController.Finish()

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.PriceUpdate{
				Asset:                       asset,
				Currency:                    currency,
				PriceInCurrencyPerAssetUnit: priceInCurrencyPerAssetUnit,
			},
		}).BuildSignEncode(priceFeederKP)

		helper.mockUseCase.EXPECT().
			UpdatePrice(context.Background(), gomock.AssignableToTypeOf(&vtxnbuild.VeloTx{})).
			Return(mockedOutput, nil)

		output, err := helper.handler.SubmitVeloTx(context.Background(), &spec.VeloTxRequest{
			SignedVeloTxXdr: veloTxB64,
		})

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, "AAAAA...=", output.SignedStellarTxXdr)
		assert.Equal(t, constants.ReplyPriceUpdateSuccess, output.Message)
		assert.Equal(t, asset, output.PriceUpdateOpResponse.Asset)
		assert.Equal(t, currency, output.PriceUpdateOpResponse.Currency)
		assert.Equal(t, priceInCurrencyPerAssetUnit, output.PriceUpdateOpResponse.PriceInCurrencyPerAssetUnit)
	})

	t.Run("error, use case return error", func(t *testing.T) {
		helper := initTest(t)
		helper.mockController.Finish()

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.PriceUpdate{
				Asset:                       asset,
				Currency:                    currency,
				PriceInCurrencyPerAssetUnit: priceInCurrencyPerAssetUnit,
			},
		}).BuildSignEncode(priceFeederKP)

		helper.mockUseCase.EXPECT().
			UpdatePrice(context.Background(), gomock.AssignableToTypeOf(&vtxnbuild.VeloTx{})).
			Return(nil, nerrors.ErrInternal{Message: "some error has occurred"})

		_, err := helper.handler.SubmitVeloTx(context.Background(), &spec.VeloTxRequest{
			SignedVeloTxXdr: veloTxB64,
		})

		assert.Error(t, err)
		assert.Equal(t, "rpc error: code = Internal desc = some error has occurred", err.Error())
	})
}

package grpc_test

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/golang/mock/gomock"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/assert"
	spec "gitlab.com/velo-labs/cen/grpc"
	"gitlab.com/velo-labs/cen/libs/convert"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/errors"
	"testing"
)

func TestHandler_SubmitVeloTx_PriceUpdate(t *testing.T) {
	var (
		priceFeederKP, _ = vconvert.SecretKeyToKeyPair(secretKey1)
	)
	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		helper.mockController.Finish()

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.PriceUpdate{
				Asset:                       "VELO",
				Currency:                    "THB",
				PriceInCurrencyPerAssetUnit: "1",
			},
		}).BuildSignEncode(priceFeederKP)

		helper.mockUseCase.EXPECT().
			UpdatePrice(context.Background(), gomock.AssignableToTypeOf(&vtxnbuild.VeloTx{})).
			Return(pointer.ToString("AAAAA...="), nil)

		reply, err := helper.handler.SubmitVeloTx(context.Background(), &spec.VeloTxRequest{
			SignedVeloTxXdr: veloTxB64,
		})

		assert.NoError(t, err)
		assert.Equal(t, "AAAAA...=", reply.SignedStellarTxXdr)
		assert.Equal(t, constants.ReplyPriceUpdateSuccess, reply.Message)
	})

	t.Run("error, use case return error", func(t *testing.T) {
		helper := initTest(t)
		helper.mockController.Finish()

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.PriceUpdate{
				Asset:                       "VELO",
				Currency:                    "THB",
				PriceInCurrencyPerAssetUnit: "1",
			},
		}).BuildSignEncode(priceFeederKP)

		helper.mockUseCase.EXPECT().
			UpdatePrice(context.Background(), gomock.AssignableToTypeOf(&vtxnbuild.VeloTx{})).
			Return(nil, nerrors.ErrInternal{Message: "some error has occurred"})

		_, err := helper.handler.SubmitVeloTx(context.Background(), &spec.VeloTxRequest{
			SignedVeloTxXdr: veloTxB64,
		})

		assert.Error(t, err)
	})
}

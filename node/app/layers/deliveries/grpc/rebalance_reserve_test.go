package grpc_test

import (
	"context"
	"github.com/AlekSi/pointer"
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

func TestHandler_RebalanceReserve(t *testing.T) {

	var (
		assetCode           = "VELO"
		assetIssuer         = "GCQDOOHRLBZW2A6COOMMWI5RAKGEZPBXSGZ6L6WA7M7GK3ZMHODDRAS3"
		requiredAmount      = decimal.NewFromFloat(150.2230124)
		poolAmount          = decimal.NewFromFloat(250.0092210)
		equalRequiredAmount = decimal.NewFromFloat(250.0092210)
		equalPoolAmount     = decimal.NewFromFloat(250.0092210)

		mockedVeloTxb = "AAAAAGqNwzi4rQDI2eTalxx56rODZdWROenUGE4mxojW0+y7AAAABQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAHW0+y7AAAAQF02xQ66AyeS6jes18Pvuz3ZADLch3Li60sQpVxAAc3IhNEjveBs7K/U0w7MDRw054lXnIPFROuzXovbz3+C9wM="

		clientKP, _ = vconvert.SecretKeyToKeyPair(secretKey1)
	)

	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		helper.mockController.Finish()

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.RebalanceReserve{},
		}).BuildSignEncode(clientKP)

		helper.mockUseCase.EXPECT().
			RebalanceReserve(context.Background(), gomock.AssignableToTypeOf(&vtxnbuild.VeloTx{})).Return(&entities.RebalanceOutput{
			Collaterals: []*entities.Collateral{{
				AssetCode:      assetCode,
				AssetIssuer:    assetIssuer,
				RequiredAmount: requiredAmount,
				PoolAmount:     poolAmount,
			}},
			SignedStellarTxXdr: &veloTxB64,
		}, nil)

		rebalanceReserveOutput, err := helper.handler.SubmitVeloTx(context.Background(), &spec.VeloTxRequest{
			SignedVeloTxXdr: veloTxB64,
		})

		assert.NoError(t, err)
		assert.Equal(t, mockedVeloTxb, rebalanceReserveOutput.SignedStellarTxXdr)
		assert.Equal(t, constants.ReplyRebalanceReserveSuccess, rebalanceReserveOutput.Message)
	})

	t.Run("success, with equal collateral amount no signed transaction", func(t *testing.T) {
		helper := initTest(t)
		helper.mockController.Finish()

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.RebalanceReserve{},
		}).BuildSignEncode(clientKP)

		helper.mockUseCase.EXPECT().
			RebalanceReserve(context.Background(), gomock.AssignableToTypeOf(&vtxnbuild.VeloTx{})).Return(&entities.RebalanceOutput{
			Collaterals: []*entities.Collateral{{
				AssetCode:      assetCode,
				AssetIssuer:    assetIssuer,
				RequiredAmount: equalRequiredAmount,
				PoolAmount:     equalPoolAmount,
			}},
			SignedStellarTxXdr: pointer.ToString(""),
		}, nil)

		rebalanceReserveOutput, err := helper.handler.SubmitVeloTx(context.Background(), &spec.VeloTxRequest{
			SignedVeloTxXdr: veloTxB64,
		})

		assert.NoError(t, err)
		assert.Equal(t, "", rebalanceReserveOutput.SignedStellarTxXdr)
		assert.Equal(t, constants.ReplyRebalanceReserveSuccess, rebalanceReserveOutput.Message)
	})

	t.Run("error, use case return error", func(t *testing.T) {
		helper := initTest(t)
		helper.mockController.Finish()

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.RebalanceReserve{},
		}).BuildSignEncode(clientKP)

		helper.mockUseCase.EXPECT().RebalanceReserve(context.Background(), gomock.AssignableToTypeOf(&vtxnbuild.VeloTx{})).Return(nil, nerrors.ErrInternal{Message: "some error has occurred"})

		rebalanceReserveOutput, err := helper.handler.SubmitVeloTx(context.Background(), &spec.VeloTxRequest{
			SignedVeloTxXdr: veloTxB64,
		})

		assert.Error(t, err)
		assert.Nil(t, rebalanceReserveOutput)
	})
}

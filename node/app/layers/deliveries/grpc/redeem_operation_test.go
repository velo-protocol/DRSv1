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

func TestHandler_SubmitVeloTx_Redeem(t *testing.T) {

	var (
		clientKP, _ = vconvert.SecretKeyToKeyPair(secretKey1)

		assetCodeToBeRedeemed   = "vTHB"
		assetIssuerToBeRedeemed = "GBVI3QZYXCWQBSGZ4TNJOHDZ5KZYGZOVSE46TVAYJYTMNCGW2PWLWO73"
		assetAmountToBeRedeemed = "1"

		collateralCode   = "VELO"
		collateralIssuer = "GBVI3QZYXCWQBSGZ4TNJOHDZ5KZYGZOVSE46TVAYJYTMNCGW2PWLWO73"
		collateralAmount = "1.0000000"
	)

	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		helper.mockController.Finish()

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.RedeemCredit{
				AssetCode: assetCodeToBeRedeemed,
				Issuer:    assetIssuerToBeRedeemed,
				Amount:    assetAmountToBeRedeemed,
			},
		}).BuildSignEncode(clientKP)

		helper.mockUseCase.EXPECT().
			RedeemCredit(context.Background(), gomock.AssignableToTypeOf(&vtxnbuild.VeloTx{})).
			Return(&entities.RedeemCreditOutput{
				SignedStellarTxXdr:      "AAAAA...=",
				AssetCodeToBeRedeemed:   assetCodeToBeRedeemed,
				AssetIssuerToBeRedeemed: assetIssuerToBeRedeemed,
				AssetAmountToBeRedeemed: decimal.New(1, 0).Truncate(7),
				CollateralCode:          collateralCode,
				CollateralIssuer:        collateralIssuer,
				CollateralAmount:        decimal.New(1, 0).Truncate(7),
			}, nil)

		reply, err := helper.handler.SubmitVeloTx(context.Background(), &spec.VeloTxRequest{
			SignedVeloTxXdr: veloTxB64,
		})

		assert.NoError(t, err)
		assert.Equal(t, "AAAAA...=", reply.SignedStellarTxXdr)
		assert.Equal(t, constants.ReplyRedeemCreditSuccess, reply.Message)
		assert.NotEmpty(t, reply.RedeemCreditOpResponse)

		assert.Equal(t, assetCodeToBeRedeemed, reply.RedeemCreditOpResponse.AssetCodeToBeRedeemed)
		assert.Equal(t, assetIssuerToBeRedeemed, reply.RedeemCreditOpResponse.AssetIssuerToBeRedeemed)
		assert.Equal(t, "1.0000000", reply.RedeemCreditOpResponse.AssetAmountToBeRedeemed)

		assert.Equal(t, collateralCode, reply.RedeemCreditOpResponse.CollateralCode)
		assert.Equal(t, collateralIssuer, reply.RedeemCreditOpResponse.CollateralIssuer)
		assert.Equal(t, collateralAmount, reply.RedeemCreditOpResponse.CollateralAmount)
	})

	t.Run("error, use case return error", func(t *testing.T) {
		helper := initTest(t)
		helper.mockController.Finish()

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.RedeemCredit{
				AssetCode: assetCodeToBeRedeemed,
				Issuer:    assetIssuerToBeRedeemed,
				Amount:    assetAmountToBeRedeemed,
			},
		}).BuildSignEncode(clientKP)

		helper.mockUseCase.EXPECT().
			RedeemCredit(context.Background(), gomock.AssignableToTypeOf(&vtxnbuild.VeloTx{})).
			Return(nil, nerrors.ErrInternal{Message: "some error has occurred"})

		_, err := helper.handler.SubmitVeloTx(context.Background(), &spec.VeloTxRequest{
			SignedVeloTxXdr: veloTxB64,
		})

		assert.Error(t, err)
	})

}

package grpc_test

import (
	"context"
	"fmt"
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

func TestHandler_SubmitVeloTx_Mint(t *testing.T) {

	var (
		trustedPartnerKP, _   = vconvert.SecretKeyToKeyPair(secretKey1)
		assetAmountToBeIssued = decimal.New(52702950798, -8)  // 527.02950798
		collateralAmount      = decimal.New(100045690008, -8) // 1000.45690008
		assetCodeToBeIssued   = "vTHB"
		collateralAssetCode   = "VELO"
	)

	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		helper.mockController.Finish()

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.MintCredit{
				AssetCodeToBeIssued: assetCodeToBeIssued,
				CollateralAssetCode: collateralAssetCode,
				CollateralAmount:    "1000.4569",
			},
		}).BuildSignEncode(trustedPartnerKP)

		helper.mockUseCase.EXPECT().
			MintCredit(context.Background(), gomock.AssignableToTypeOf(&vtxnbuild.VeloTx{})).
			Return(&entities.MintCreditOutput{
				SignedStellarTxXdr:    "AAAAA...=",
				AssetAmountToBeIssued: assetAmountToBeIssued,
				AssetCodeToBeIssued:   assetCodeToBeIssued,
				CollateralAmount:      collateralAmount,
				CollateralAssetCode:   collateralAssetCode,
			}, nil)

		reply, err := helper.handler.SubmitVeloTx(context.Background(), &spec.VeloTxRequest{
			SignedVeloTxXdr: veloTxB64,
		})

		assert.NoError(t, err)
		assert.Equal(t, "AAAAA...=", reply.SignedStellarTxXdr)
		assert.Equal(t,
			fmt.Sprintf(constants.ReplyMintCreditSuccess, "527.0295079", "vTHB", "1000.4569000", "VELO"),
			reply.Message,
		)
		assert.Equal(t, assetCodeToBeIssued, reply.MintCreditOpResponse.AssetCodeToBeIssued)
		assert.Equal(t, collateralAssetCode, reply.MintCreditOpResponse.CollateralAssetCode)
		assert.NotEmpty(t, reply.MintCreditOpResponse.AssetAmountToBeIssued)
		assert.NotEmpty(t, reply.MintCreditOpResponse.CollateralAmount)
	})

	t.Run("error, use case return error", func(t *testing.T) {
		helper := initTest(t)
		helper.mockController.Finish()

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.MintCredit{
				AssetCodeToBeIssued: "vTHB",
				CollateralAssetCode: "VELO",
				CollateralAmount:    "1000.4569",
			},
		}).BuildSignEncode(trustedPartnerKP)

		helper.mockUseCase.EXPECT().
			MintCredit(context.Background(), gomock.AssignableToTypeOf(&vtxnbuild.VeloTx{})).
			Return(nil, nerrors.ErrInternal{Message: "some error has occurred"})

		_, err := helper.handler.SubmitVeloTx(context.Background(), &spec.VeloTxRequest{
			SignedVeloTxXdr: veloTxB64,
		})

		assert.Error(t, err)
	})
}

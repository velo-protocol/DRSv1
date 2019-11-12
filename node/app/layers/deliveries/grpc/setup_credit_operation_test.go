package grpc_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/assert"
	spec "gitlab.com/velo-labs/cen/grpc"
	"gitlab.com/velo-labs/cen/libs/convert"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/errors"
	"testing"
)

func TestHandler_SubmitVeloTx_SetupCredit(t *testing.T) {

	var (
		trustedPartnerKP, _ = vconvert.SecretKeyToKeyPair(secretKey1)
	)

	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		helper.mockController.Finish()

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.SetupCredit{
				PeggedValue:    "1.00",
				PeggedCurrency: "THB",
				AssetCode:      "vTHB",
			},
		}).BuildSignEncode(trustedPartnerKP)

		helper.mockUseCase.EXPECT().
			SetupCredit(context.Background(), gomock.AssignableToTypeOf(&vtxnbuild.VeloTx{})).
			Return(&entities.SetupCreditOutput{
				SignedStellarTxXdr: "AAAAA...=",
			}, nil)

		reply, err := helper.handler.SubmitVeloTx(context.Background(), &spec.VeloTxRequest{
			SignedVeloTxXdr: veloTxB64,
		})

		assert.NoError(t, err)
		assert.Equal(t, "AAAAA...=", reply.SignedStellarTxXdr)
		assert.Equal(t, constants.ReplySetupCreditSuccess, reply.Message)
		assert.NotNil(t, reply.SetupCreditOpResponse)
	})
	t.Run("error, use case return error", func(t *testing.T) {
		helper := initTest(t)
		helper.mockController.Finish()

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.SetupCredit{
				PeggedValue:    "1.00",
				PeggedCurrency: "THB",
				AssetCode:      "vTHB",
			},
		}).BuildSignEncode(trustedPartnerKP)

		helper.mockUseCase.EXPECT().
			SetupCredit(context.Background(), gomock.AssignableToTypeOf(&vtxnbuild.VeloTx{})).
			Return(nil, nerrors.ErrInternal{Message: "some error has occurred"})

		_, err := helper.handler.SubmitVeloTx(context.Background(), &spec.VeloTxRequest{
			SignedVeloTxXdr: veloTxB64,
		})

		assert.Error(t, err)
	})

}

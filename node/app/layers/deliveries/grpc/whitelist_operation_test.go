package grpc_test

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/assert"
	spec "github.com/velo-protocol/DRSv1/grpc"
	"github.com/velo-protocol/DRSv1/libs/convert"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"github.com/velo-protocol/DRSv1/libs/xdr"
	"github.com/velo-protocol/DRSv1/node/app/constants"
	"github.com/velo-protocol/DRSv1/node/app/entities"
	"github.com/velo-protocol/DRSv1/node/app/errors"
	"testing"
)

func TestHandler_SubmitVeloTx_Whitelist(t *testing.T) {

	var (
		regulatorKP, _ = vconvert.SecretKeyToKeyPair(secretKey1)
	)

	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		helper.mockController.Finish()

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.Whitelist{
				Address: publicKey2,
				Role:    string(vxdr.RoleTrustedPartner),
			},
		}).BuildSignEncode(regulatorKP)

		helper.mockUseCase.EXPECT().
			Whitelist(context.Background(), gomock.AssignableToTypeOf(&vtxnbuild.VeloTx{})).
			Return(&entities.WhitelistOutput{
				SignedStellarTxXdr: "AAAAA...=",
				Address:            publicKey2,
				Role:               string(vxdr.RoleTrustedPartner),
			}, nil)

		reply, err := helper.handler.SubmitVeloTx(context.Background(), &spec.VeloTxRequest{
			SignedVeloTxXdr: veloTxB64,
		})

		assert.NoError(t, err)
		assert.Equal(t, "AAAAA...=", reply.SignedStellarTxXdr)
		assert.Equal(t, fmt.Sprintf(constants.ReplyWhitelistSuccess, publicKey2, vxdr.RoleMap[vxdr.RoleTrustedPartner]), reply.Message)
		assert.Equal(t, string(vxdr.RoleTrustedPartner), reply.WhitelistOpResponse.Role)
		assert.Equal(t, publicKey2, reply.WhitelistOpResponse.Address)
	})
	t.Run("error, use case return error", func(t *testing.T) {
		helper := initTest(t)
		helper.mockController.Finish()

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.Whitelist{
				Address: publicKey2,
				Role:    string(vxdr.RoleTrustedPartner),
			},
		}).BuildSignEncode(regulatorKP)

		helper.mockUseCase.EXPECT().
			Whitelist(context.Background(), gomock.AssignableToTypeOf(&vtxnbuild.VeloTx{})).
			Return(nil, nerrors.ErrInternal{Message: "some error has occurred"})

		_, err := helper.handler.SubmitVeloTx(context.Background(), &spec.VeloTxRequest{
			SignedVeloTxXdr: veloTxB64,
		})

		assert.Error(t, err)
	})

}

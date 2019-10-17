package grpc_test

import (
	"context"
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/golang/mock/gomock"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/assert"
	spec "gitlab.com/velo-labs/cen/grpc"
	"gitlab.com/velo-labs/cen/libs/convert"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/errors"
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
			CreateWhitelist(context.Background(), gomock.AssignableToTypeOf(&vtxnbuild.VeloTx{})).
			Return(pointer.ToString("AAAAA...="), nil)

		reply, err := helper.handler.SubmitVeloTx(context.Background(), &spec.VeloTxRequest{
			SignedVeloTxXdr: veloTxB64,
		})

		assert.NoError(t, err)
		assert.Equal(t, "AAAAA...=", reply.SignedStellarTxXdr)
		assert.Equal(t, fmt.Sprintf(constants.ReplyWhitelistSuccess, publicKey2, vxdr.RoleMap[vxdr.RoleTrustedPartner]), reply.Message)
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
			CreateWhitelist(context.Background(), gomock.AssignableToTypeOf(&vtxnbuild.VeloTx{})).
			Return(nil, nerrors.ErrInternal{Message: "some error has occurred"})

		_, err := helper.handler.SubmitVeloTx(context.Background(), &spec.VeloTxRequest{
			SignedVeloTxXdr: veloTxB64,
		})

		assert.Error(t, err)
	})

}

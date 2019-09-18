package grpc

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/assert"
	spec "gitlab.com/velo-labs/cen/grpc"
	"gitlab.com/velo-labs/cen/libs/convert"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/layers/mocks"
	"testing"
)

func TestHandler_SubmitVeloTx(t *testing.T) {
	const (
		publicKey1 = "GBVI3QZYXCWQBSGZ4TNJOHDZ5KZYGZOVSE46TVAYJYTMNCGW2PWLWO73"
		secretKey1 = "SBR25NMQRKQ4RLGNV5XB3MMQB4ADVYSMPGVBODQVJE7KPTDR6KGK3XMX"
		publicKey2 = "GC2ROYZQH5FTVEPQZF7CAB32SCJC7DWVKILDUAT5BCU5O7HEI7HFUB25"
		//secretKey2 = "SCHQI345PYWHM2APNR4MN433HNCBS7VDUROOZKTYHZUBBTHI2YHNCJ4G"
	)

	var (
		kp1, _ = vconvert.SecretKeyToKeyPair(secretKey1)
		//kp2, _ = vconvert.SecretKeyToKeyPair(secretKey2)

		newMockedUseCase = func() (*mocks.MockUseCase, func()) {
			ctrl := gomock.NewController(t)
			mockedUseCase := mocks.NewMockUseCase(ctrl)
			return mockedUseCase, ctrl.Finish
		}
	)

	t.Run("error, cannot unmarshal xdr string to VeloTx", func(t *testing.T) {
		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.WhiteList{},
		}).BuildSignEncode(kp1)

		_, err := (&handler{}).SubmitVeloTx(context.Background(), &spec.VeloTxRequest{
			SignedVeloTxXdr: veloTxB64,
		})

		assert.Error(t, err)
	})

	t.Run("must be able to handle white list operation", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			mockedUseCase, finish := newMockedUseCase()
			defer finish()

			veloTxB64, _ := (&vtxnbuild.VeloTx{
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: publicKey1,
				},
				VeloOp: &vtxnbuild.WhiteList{
					Address: publicKey2,
					Role:    string(vxdr.RoleTrustedPartner),
				},
			}).BuildSignEncode(kp1)

			mockedUseCase.EXPECT().
				CreateWhiteList(context.Background(), gomock.AssignableToTypeOf(&vxdr.VeloTxEnvelope{})).
				Return(nil)

			reply, err := (&handler{mockedUseCase}).SubmitVeloTx(context.Background(), &spec.VeloTxRequest{
				SignedVeloTxXdr: veloTxB64,
			})

			assert.NoError(t, err)
			assert.Equal(t, "", reply.SignedStellarTxXdr)
		})
		t.Run("error, use case return error", func(t *testing.T) {
			mockedUseCase, finish := newMockedUseCase()
			defer finish()

			veloTxB64, _ := (&vtxnbuild.VeloTx{
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: publicKey1,
				},
				VeloOp: &vtxnbuild.WhiteList{
					Address: publicKey2,
					Role:    string(vxdr.RoleTrustedPartner),
				},
			}).BuildSignEncode(kp1)

			mockedUseCase.EXPECT().
				CreateWhiteList(context.Background(), gomock.AssignableToTypeOf(&vxdr.VeloTxEnvelope{})).
				Return(errors.New("some error has occurred"))

			_, err := (&handler{mockedUseCase}).SubmitVeloTx(context.Background(), &spec.VeloTxRequest{
				SignedVeloTxXdr: veloTxB64,
			})

			assert.Error(t, err)
		})
	})
}

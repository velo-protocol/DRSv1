package grpc_test

import (
	"context"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/assert"
	spec "gitlab.com/velo-labs/cen/grpc"
	"gitlab.com/velo-labs/cen/libs/convert"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"testing"
)

func TestHandler_SubmitVeloTx(t *testing.T) {

	var (
		kp1, _ = vconvert.SecretKeyToKeyPair(secretKey1)
	)

	t.Run("error, cannot unmarshal xdr string to VeloTx", func(t *testing.T) {
		helper := initTest(t)

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.Whitelist{},
		}).BuildSignEncode(kp1)

		_, err := helper.handler.SubmitVeloTx(context.Background(), &spec.VeloTxRequest{
			SignedVeloTxXdr: veloTxB64,
		})

		assert.Error(t, err)
	})

}

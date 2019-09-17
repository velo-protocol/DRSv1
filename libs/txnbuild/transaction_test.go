package vtxnbuild

import (
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/assert"

	"gitlab.com/velo-labs/cen/libs/xdr"
	"testing"
)

func TestVeloTx_Build(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		veloTx := VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &WhiteList{
				Address: publicKey2,
				Role:    string(vxdr.RoleTrustedPartner),
			},
		}

		err := veloTx.Build()
		assert.NoError(t, err)

		// assert data in xdr tx
		assert.Equal(t, vxdr.OperationTypeWhiteList, veloTx.veloXdrTx.VeloOp.Body.Type)
		assert.Equal(t, publicKey1, veloTx.veloXdrTx.SourceAccount.Address())
		assert.Equal(t, publicKey2, veloTx.veloXdrTx.VeloOp.Body.WhiteListOp.Address.Address())
		assert.Equal(t, vxdr.RoleTrustedPartner, veloTx.veloXdrTx.VeloOp.Body.WhiteListOp.Role)

		// assert equality of tx in envelope and the velo tx
		assert.Equal(t, veloTx.veloXdrTx.VeloOp, veloTx.veloXdrEnvelope.VeloTx.VeloOp)
	})
	// TODO: Add error test case
}

func TestVeloTx_Sign(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		veloTx := VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &WhiteList{
				Address: publicKey2,
				Role:    string(vxdr.RoleTrustedPartner),
			},
		}
		_ = veloTx.Build()
		err := veloTx.Sign(kp1, kp2, kp1)

		assert.NoError(t, err)

		// assert data in xdr tx
		assert.Equal(t, vxdr.OperationTypeWhiteList, veloTx.veloXdrTx.VeloOp.Body.Type)
		assert.Equal(t, publicKey1, veloTx.veloXdrTx.SourceAccount.Address())
		assert.Equal(t, publicKey2, veloTx.veloXdrTx.VeloOp.Body.WhiteListOp.Address.Address())
		assert.Equal(t, vxdr.RoleTrustedPartner, veloTx.veloXdrTx.VeloOp.Body.WhiteListOp.Role)

		// assert data in xdr envelope
		assert.Len(t, veloTx.veloXdrEnvelope.Signatures, 3)
		// first and third signature should be equal
		assert.True(t,
			string(veloTx.veloXdrEnvelope.Signatures[0].Signature) ==
				string(veloTx.veloXdrEnvelope.Signatures[2].Signature))
		// first and second signature shouldn't be equal
		assert.False(t,
			string(veloTx.veloXdrEnvelope.Signatures[0].Signature) ==
				string(veloTx.veloXdrEnvelope.Signatures[1].Signature))
	})
	// TODO: Add error test case
}

func TestVeloTx_BuildSignEncode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		veloTxB64, err := (&VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &WhiteList{
				Address: publicKey2,
				Role:    string(vxdr.RoleTrustedPartner),
			},
		}).BuildSignEncode(kp1, kp2)

		assert.NoError(t, err)
		assert.NotEmpty(t, veloTxB64)
	})
	// TODO: Add error test case
}

func TestTransactionFromXDR(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		veloTxB64, _ := (&VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &WhiteList{
				Address: publicKey2,
				Role:    string(vxdr.RoleTrustedPartner),
			},
		}).BuildSignEncode(kp1, kp2)

		veloTx, err := TransactionFromXDR(veloTxB64)
		assert.NoError(t, err)

		// assert data in xdr tx
		assert.Equal(t, vxdr.OperationTypeWhiteList, veloTx.veloXdrTx.VeloOp.Body.Type)
		assert.Equal(t, publicKey1, veloTx.veloXdrTx.SourceAccount.Address())
		assert.Equal(t, publicKey2, veloTx.veloXdrTx.VeloOp.Body.WhiteListOp.Address.Address())
		assert.Equal(t, vxdr.RoleTrustedPartner, veloTx.veloXdrTx.VeloOp.Body.WhiteListOp.Role)

		// assert equality of tx in envelope and the velo tx
		assert.Equal(t, veloTx.veloXdrTx.VeloOp, veloTx.veloXdrEnvelope.VeloTx.VeloOp)

		// assert data in xdr envelope
		assert.Len(t, veloTx.veloXdrEnvelope.Signatures, 2)
		assert.False(t,
			string(veloTx.veloXdrEnvelope.Signatures[0].Signature) ==
				string(veloTx.veloXdrEnvelope.Signatures[1].Signature))
	})
	// TODO: Add error test case
}

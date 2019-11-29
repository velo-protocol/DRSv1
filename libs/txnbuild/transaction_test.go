package vtxnbuild

import (
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
	"github.com/stretchr/testify/assert"

	"github.com/velo-protocol/DRSv1/libs/xdr"
	"testing"
)

func TestVeloTx_Build(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		veloTx := VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &Whitelist{
				Address: publicKey2,
				Role:    string(vxdr.RoleTrustedPartner),
			},
		}

		err := veloTx.Build()
		assert.NoError(t, err)

		// assert data in xdr tx
		assert.Equal(t, vxdr.OperationTypeWhitelist, veloTx.veloXdrTx.VeloOp.Body.Type)
		assert.Equal(t, publicKey1, veloTx.veloXdrTx.SourceAccount.Address())
		assert.Equal(t, publicKey2, veloTx.veloXdrTx.VeloOp.Body.WhitelistOp.Address.Address())
		assert.Equal(t, vxdr.RoleTrustedPartner, veloTx.veloXdrTx.VeloOp.Body.WhitelistOp.Role)

		// assert equality of tx in envelope and the velo tx
		assert.Equal(t, veloTx.veloXdrTx.VeloOp, veloTx.veloXdrEnvelope.VeloTx.VeloOp)
	})

	t.Run("success, tx has been assigned with envelope, reset it", func(t *testing.T) {
		veloTx := VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &Whitelist{
				Address: publicKey2,
				Role:    string(vxdr.RoleTrustedPartner),
			},
		}
		veloTx.veloXdrEnvelope = &vxdr.VeloTxEnvelope{}

		err := veloTx.Build()
		assert.NoError(t, err)

		// assert data in xdr tx
		assert.Equal(t, vxdr.OperationTypeWhitelist, veloTx.veloXdrTx.VeloOp.Body.Type)
		assert.Equal(t, publicKey1, veloTx.veloXdrTx.SourceAccount.Address())
		assert.Equal(t, publicKey2, veloTx.veloXdrTx.VeloOp.Body.WhitelistOp.Address.Address())
		assert.Equal(t, vxdr.RoleTrustedPartner, veloTx.veloXdrTx.VeloOp.Body.WhitelistOp.Role)

		// assert equality of tx in envelope and the velo tx
		assert.Equal(t, veloTx.veloXdrTx.VeloOp, veloTx.veloXdrEnvelope.VeloTx.VeloOp)
	})

	t.Run("error, tx has already been signed, cannot be rebuilt", func(t *testing.T) {
		veloTx := VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &Whitelist{
				Address: publicKey2,
				Role:    string(vxdr.RoleTrustedPartner),
			},
		}

		veloTx.veloXdrEnvelope = &vxdr.VeloTxEnvelope{
			Signatures: []xdr.DecoratedSignature{
				{Signature: []byte("SOME_SIGNATURE")},
			},
		}

		err := veloTx.Build()
		assert.Error(t, err)
	})

	t.Run("error, source account cannot be nil", func(t *testing.T) {
		veloTx := VeloTx{
			SourceAccount: nil,
			VeloOp: &Whitelist{
				Address: publicKey2,
				Role:    string(vxdr.RoleTrustedPartner),
			},
		}

		err := veloTx.Build()
		assert.Error(t, err)
	})
	t.Run("error, source account cannot be blank", func(t *testing.T) {
		veloTx := VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: "",
			},
			VeloOp: &Whitelist{
				Address: publicKey2,
				Role:    string(vxdr.RoleTrustedPartner),
			},
		}

		err := veloTx.Build()
		assert.Error(t, err)
	})
	t.Run("error, bad source account", func(t *testing.T) {
		veloTx := VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: "BAD_PUBLIC_KY",
			},
			VeloOp: &Whitelist{
				Address: publicKey2,
				Role:    string(vxdr.RoleTrustedPartner),
			},
		}

		err := veloTx.Build()
		assert.Error(t, err)
	})

	t.Run("error, VeloOp validation fail, bad role", func(t *testing.T) {
		veloTx := VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &Whitelist{
				Address: publicKey2,
				Role:    "BAD_ROLE",
			},
		}

		err := veloTx.Build()
		assert.Error(t, err)
	})
}

func TestVeloTx_Sign(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		veloTx := VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &Whitelist{
				Address: publicKey2,
				Role:    string(vxdr.RoleTrustedPartner),
			},
		}
		_ = veloTx.Build()
		err := veloTx.Sign(kp1, kp2, kp1)

		assert.NoError(t, err)

		// assert data in xdr tx
		assert.Equal(t, vxdr.OperationTypeWhitelist, veloTx.veloXdrTx.VeloOp.Body.Type)
		assert.Equal(t, publicKey1, veloTx.veloXdrTx.SourceAccount.Address())
		assert.Equal(t, publicKey2, veloTx.veloXdrTx.VeloOp.Body.WhitelistOp.Address.Address())
		assert.Equal(t, vxdr.RoleTrustedPartner, veloTx.veloXdrTx.VeloOp.Body.WhitelistOp.Role)

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
}

func TestVeloTx_BuildSignEncode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		veloTxB64, err := (&VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &Whitelist{
				Address: publicKey2,
				Role:    string(vxdr.RoleTrustedPartner),
			},
		}).BuildSignEncode(kp1, kp2)

		assert.NoError(t, err)
		assert.NotEmpty(t, veloTxB64)
	})
}

func TestTransactionFromXDR(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		veloTxB64, _ := (&VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &Whitelist{
				Address: publicKey2,
				Role:    string(vxdr.RoleTrustedPartner),
			},
		}).BuildSignEncode(kp1, kp2)

		veloTx, err := TransactionFromXDR(veloTxB64)
		assert.NoError(t, err)

		// assert data in xdr tx
		assert.Equal(t, vxdr.OperationTypeWhitelist, veloTx.veloXdrTx.VeloOp.Body.Type)
		assert.Equal(t, publicKey1, veloTx.veloXdrTx.SourceAccount.Address())
		assert.Equal(t, publicKey2, veloTx.veloXdrTx.VeloOp.Body.WhitelistOp.Address.Address())
		assert.Equal(t, vxdr.RoleTrustedPartner, veloTx.veloXdrTx.VeloOp.Body.WhitelistOp.Role)

		// assert equality of tx in envelope and the velo tx
		assert.Equal(t, veloTx.veloXdrTx.VeloOp, veloTx.veloXdrEnvelope.VeloTx.VeloOp)

		// assert data in xdr envelope
		assert.Len(t, veloTx.veloXdrEnvelope.Signatures, 2)
		assert.False(t,
			string(veloTx.veloXdrEnvelope.Signatures[0].Signature) ==
				string(veloTx.veloXdrEnvelope.Signatures[1].Signature))
	})

	t.Run("error, unable to unmarshal velo xdr", func(t *testing.T) {
		_, err := TransactionFromXDR("AAAAAAAABAD_XDR")
		assert.Error(t, err)
	})
}

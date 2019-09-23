package vtxnbuild

import (
	"github.com/stellar/go/xdr"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"testing"
)

func TestWhiteList_BuildXDR(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		veloXdrOp, err := (&WhiteList{
			Address: publicKey1,
			Role:    string(vxdr.RoleRegulator),
		}).BuildXDR()

		assert.NoError(t, err)
		assert.Equal(t, vxdr.OperationTypeWhiteList, veloXdrOp.Body.Type)
		assert.Equal(t, vxdr.RoleRegulator, veloXdrOp.Body.WhiteListOp.Role)
		assert.Equal(t, publicKey1, veloXdrOp.Body.WhiteListOp.Address.Address())
	})
	t.Run("error, address cannot be blank", func(t *testing.T) {
		_, err := (&WhiteList{
			Address: "",
			Role:    string(vxdr.RoleRegulator),
		}).BuildXDR()

		assert.Error(t, err)
	})
	t.Run("error, role cannot be blank", func(t *testing.T) {
		_, err := (&WhiteList{
			Address: publicKey1,
			Role:    "",
		}).BuildXDR()

		assert.Error(t, err)
	})
	t.Run("error, bad address format", func(t *testing.T) {
		_, err := (&WhiteList{
			Address: "BAD_PK",
			Role:    string(vxdr.RoleRegulator),
		}).BuildXDR()

		assert.Error(t, err)
	})
	t.Run("error, invalid role", func(t *testing.T) {
		_, err := (&WhiteList{
			Address: publicKey1,
			Role:    "BAD_ROLE",
		}).BuildXDR()

		assert.Error(t, err)
	})
}

func TestWhiteList_FromXDR(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var account xdr.AccountId
		_ = account.SetAddress(publicKey1)
		veloXdrOp := vxdr.VeloOp{
			Body: vxdr.OperationBody{
				WhiteListOp: &vxdr.WhiteListOp{
					Address: account,
					Role:    vxdr.RoleRegulator,
				},
			},
		}

		var newVeloWhiteListOp WhiteList
		err := newVeloWhiteListOp.FromXDR(veloXdrOp)

		assert.NoError(t, err)
		assert.Equal(t, publicKey1, newVeloWhiteListOp.Address)
		assert.Equal(t, string(vxdr.RoleRegulator), newVeloWhiteListOp.Role)
	})
	t.Run("error, empty WhiteListOp", func(t *testing.T) {
		veloXdrOp := vxdr.VeloOp{
			Body: vxdr.OperationBody{
				WhiteListOp: nil,
			},
		}

		var newVeloWhiteListOp WhiteList
		err := newVeloWhiteListOp.FromXDR(veloXdrOp)
		assert.Error(t, err)
	})
}

func TestWhiteList_Validate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		err := (&WhiteList{
			Address: publicKey1,
			Role:    string(vxdr.RoleRegulator),
		}).Validate()

		assert.NoError(t, err)
	})
	t.Run("error, invalid public key format", func(t *testing.T) {
		err := (&WhiteList{
			Address: "BAD_PK",
			Role:    string(vxdr.RoleRegulator),
		}).Validate()

		assert.Error(t, err)
	})
	t.Run("error, empty public key format", func(t *testing.T) {
		err := (&WhiteList{
			Address: "",
			Role:    string(vxdr.RoleRegulator),
		}).Validate()

		assert.Error(t, err)
	})
	t.Run("error, unknown role", func(t *testing.T) {
		err := (&WhiteList{
			Address: publicKey1,
			Role:    "BAD_ROLE",
		}).Validate()

		assert.Error(t, err)
	})
}

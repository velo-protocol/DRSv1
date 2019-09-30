package vtxnbuild

import (
	"fmt"
	"github.com/stellar/go/xdr"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"testing"
)

func TestWhiteList_BuildXDR(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		veloXdrOp, err := (&WhiteList{
			Address:  publicKey1,
			Role:     string(vxdr.RoleRegulator),
			Currency: string(vxdr.CurrencyTHB),
		}).BuildXDR()

		assert.NoError(t, err)
		assert.Equal(t, vxdr.OperationTypeWhiteList, veloXdrOp.Body.Type)
		assert.Equal(t, vxdr.RoleRegulator, veloXdrOp.Body.WhiteListOp.Role)
		assert.Equal(t, publicKey1, veloXdrOp.Body.WhiteListOp.Address.Address())
		assert.Equal(t, vxdr.CurrencyTHB, veloXdrOp.Body.WhiteListOp.Currency)
	})
	t.Run("error, validation fail", func(t *testing.T) {
		_, err := (&WhiteList{
			Address:  publicKey1,
			Role:     "",
			Currency: string(vxdr.CurrencyTHB),
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
					Address:  account,
					Role:     vxdr.RoleRegulator,
					Currency: vxdr.CurrencyTHB,
				},
			},
		}

		var newVeloWhiteListOp WhiteList
		err := newVeloWhiteListOp.FromXDR(veloXdrOp)

		assert.NoError(t, err)
		assert.Equal(t, publicKey1, newVeloWhiteListOp.Address)
		assert.Equal(t, string(vxdr.RoleRegulator), newVeloWhiteListOp.Role)
		assert.Equal(t, string(vxdr.CurrencyTHB), newVeloWhiteListOp.Currency)
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
			Address:  publicKey1,
			Role:     string(vxdr.RoleRegulator),
			Currency: string(vxdr.CurrencyTHB),
		}).Validate()

		assert.NoError(t, err)
	})
	t.Run("success, without currency", func(t *testing.T) {
		err := (&WhiteList{
			Address: publicKey1,
			Role:    string(vxdr.RoleRegulator),
		}).Validate()

		assert.NoError(t, err)
	})

	t.Run("error, address cannot be blank", func(t *testing.T) {
		err := (&WhiteList{
			Address:  "",
			Role:     string(vxdr.RoleRegulator),
			Currency: string(vxdr.CurrencyTHB),
		}).Validate()

		assert.EqualError(t, err, "address must not be blank")
	})
	t.Run("error, role cannot be blank", func(t *testing.T) {
		err := (&WhiteList{
			Address:  publicKey1,
			Role:     "",
			Currency: string(vxdr.CurrencyTHB),
		}).Validate()

		assert.EqualError(t, err, "role must not be blank")
	})

	t.Run("error, invalid public key format", func(t *testing.T) {
		err := (&WhiteList{
			Address:  "BAD_PK",
			Role:     string(vxdr.RoleRegulator),
			Currency: string(vxdr.CurrencyTHB),
		}).Validate()

		assert.EqualError(t, err, "invalid address format")
	})
	t.Run("error, unknown role", func(t *testing.T) {
		err := (&WhiteList{
			Address:  publicKey1,
			Role:     "BAD_ROLE",
			Currency: string(vxdr.CurrencyTHB),
		}).Validate()

		assert.EqualError(t, err, "role specified does not exist")
	})
	t.Run("error, unknown currency", func(t *testing.T) {
		err := (&WhiteList{
			Address:  publicKey1,
			Role:     string(vxdr.RoleRegulator),
			Currency: "BAD_CURRENCY",
		}).Validate()

		assert.EqualError(t, err, fmt.Sprintf("currency %s does not exist", "BAD_CURRENCY"))
	})
}

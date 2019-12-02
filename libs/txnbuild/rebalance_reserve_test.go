package vtxnbuild

import (
	"github.com/stellar/go/xdr"
	"github.com/stretchr/testify/assert"
	"github.com/velo-protocol/DRSv1/libs/xdr"
	"testing"
)

func TestRebalanceReserve_BuildXDR(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		veloXdrOp, err := (&RebalanceReserve{}).BuildXDR()

		assert.NoError(t, err)
		assert.Equal(t, vxdr.OperationTypeRebalanceReserve, veloXdrOp.Body.Type)
	})

}

func TestRebalanceReserve_FromXDR(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		var account xdr.AccountId
		_ = account.SetAddress(publicKey1)
		veloXdrOp := vxdr.VeloOp{
			Body: vxdr.OperationBody{
				RebalanceReserveOp: &vxdr.RebalanceReserveOp{},
			},
		}

		var newVeloRebalanceReserveOp RebalanceReserve
		err := newVeloRebalanceReserveOp.FromXDR(veloXdrOp)

		assert.NoError(t, err)
	})

	t.Run("error, empty RebalanceReserveOp", func(t *testing.T) {
		var account xdr.AccountId
		_ = account.SetAddress(publicKey1)
		veloXdrOp := vxdr.VeloOp{
			Body: vxdr.OperationBody{
				RebalanceReserveOp: nil,
			},
		}

		var newVeloRebalanceReserveOp RebalanceReserve
		err := newVeloRebalanceReserveOp.FromXDR(veloXdrOp)

		assert.Error(t, err)
		assert.Empty(t, newVeloRebalanceReserveOp)
	})
}

func TestRebalanceReserve_Validate(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		err := (&RebalanceReserve{}).Validate()
		assert.NoError(t, err)
	})
}

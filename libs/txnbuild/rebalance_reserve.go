package vtxnbuild

import (
	"github.com/pkg/errors"
	"github.com/velo-protocol/DRSv1/libs/xdr"
)

// RebalanceReserve represents the Velo rebalance reserve Operation.
type RebalanceReserve struct{}

// BuildXDR for RebalanceReserve returns a fully configured XDR Operation.
func (rebalanceReserve *RebalanceReserve) BuildXDR() (vxdr.VeloOp, error) {

	// xdr op
	vXdrOp := vxdr.RebalanceReserveOp{}

	body, err := vxdr.NewOperationBody(vxdr.OperationTypeRebalanceReserve, vXdrOp)
	if err != nil {
		return vxdr.VeloOp{}, errors.Wrap(err, "failed to build XDR operation body")
	}

	return vxdr.VeloOp{Body: body}, nil
}

// FromXDR for RebalanceReserve initialises the vtxnbuild struct from the corresponding XDR Operation.
func (rebalanceReserve *RebalanceReserve) FromXDR(vXdrOp vxdr.VeloOp) error {
	redeemOp := vXdrOp.Body.RebalanceReserveOp
	if redeemOp == nil {
		return errors.New("error parsing rebalanceReserve operation from xdr")
	}

	return nil
}

// Validation function for RebalanceReserve. Validates the required struct fields. It returns an error if any of the fields are
// invalid. Otherwise, it returns nil.
func (rebalanceReserve *RebalanceReserve) Validate() error {
	return nil
}

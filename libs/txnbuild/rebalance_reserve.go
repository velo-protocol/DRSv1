package vtxnbuild

import (
	"github.com/pkg/errors"
	"github.com/velo-protocol/DRSv1/libs/xdr"
)

type RebalanceReserve struct{}

func (rebalanceReserve *RebalanceReserve) BuildXDR() (vxdr.VeloOp, error) {

	// xdr op
	vXdrOp := vxdr.RebalanceReserveOp{}

	body, err := vxdr.NewOperationBody(vxdr.OperationTypeRebalanceReserve, vXdrOp)
	if err != nil {
		return vxdr.VeloOp{}, errors.Wrap(err, "failed to build XDR operation body")
	}

	return vxdr.VeloOp{Body: body}, nil
}

func (rebalanceReserve *RebalanceReserve) FromXDR(vXdrOp vxdr.VeloOp) error {
	redeemOp := vXdrOp.Body.RebalanceReserveOp
	if redeemOp == nil {
		return errors.New("error parsing rebalanceReserve operation from xdr")
	}

	return nil
}

func (rebalanceReserve *RebalanceReserve) Validate() error {
	return nil
}

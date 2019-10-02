package vtxnbuild

import (
	"gitlab.com/velo-labs/cen/libs/xdr"
)

type VeloOp interface {
	BuildXDR() (vxdr.VeloOp, error)
	FromXDR(vXdrOp vxdr.VeloOp) error
	Validate() error
}

func operationFromXDR(vXdrOp vxdr.VeloOp) (VeloOp, error) {
	var newVeloOp VeloOp
	switch vXdrOp.Body.Type {
	case vxdr.OperationTypeWhitelist:
		newVeloOp = &Whitelist{}
	case vxdr.OperationTypeSetupCredit:
		newVeloOp = &SetupCredit{}
	case vxdr.OperationTypePriceUpdate:
		newVeloOp = &PriceUpdate{}
	}

	err := newVeloOp.FromXDR(vXdrOp)
	return newVeloOp, err
}

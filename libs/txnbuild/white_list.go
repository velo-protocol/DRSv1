package vtxnbuild

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/strkey"
	"gitlab.com/velo-labs/cen/libs/xdr"
)

type WhiteList struct {
	Address string
	Role    string
}

type SetUpCredit struct {
	Address string
	Role    string
}

func (whiteList *WhiteList) BuildXDR() (vxdr.VeloOp, error) {
	var vXdrOp vxdr.WhiteListOp
	err := vXdrOp.Address.SetAddress(whiteList.Address)
	if err != nil {
		return vxdr.VeloOp{}, errors.Wrap(err, "failed to set whiteList address")
	}

	vXdrOp.Role = vxdr.Role(whiteList.Role)
	if !vXdrOp.Role.IsValid() {
		return vxdr.VeloOp{}, errors.New("failed to set whiteList role")
	}

	body, err := vxdr.NewOperationBody(vxdr.OperationTypeWhiteList, vXdrOp)
	if err != nil {
		return vxdr.VeloOp{}, errors.Wrap(err, "failed to build XDR OperationBody")
	}

	return vxdr.VeloOp{Body: body}, nil
}

func (whiteList *WhiteList) FromXDR(vXdrOp vxdr.VeloOp) error {
	whiteListOp := vXdrOp.Body.WhiteListOp
	if whiteListOp == nil {
		return errors.New("error parsing whiteList operation from xdr")
	}

	whiteList.Role = string(whiteListOp.Role)
	whiteList.Address = whiteListOp.Address.Address()

	return nil
}

func (whiteList *WhiteList) Validate() error {
	if whiteList.Address == "" || !strkey.IsValidEd25519PublicKey(whiteList.Address) {
		return errors.Errorf("%s is not a valid stellar public key", whiteList.Address)
	}

	if !vxdr.Role(whiteList.Role).IsValid() {
		return errors.New("invalid role")
	}

	return nil
}

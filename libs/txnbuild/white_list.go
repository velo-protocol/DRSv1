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


func (whiteList *WhiteList) BuildXDR() (vxdr.VeloOp, error) {
	if err := whiteList.Validate(); err != nil {
		return vxdr.VeloOp{}, err
	}

	var vXdrOp vxdr.WhiteListOp
	err := vXdrOp.Address.SetAddress(whiteList.Address)
	if err != nil {
		return vxdr.VeloOp{}, errors.Wrap(err, "failed to set white list address")
	}

	vXdrOp.Role = vxdr.Role(whiteList.Role)

	body, err := vxdr.NewOperationBody(vxdr.OperationTypeWhiteList, vXdrOp)
	if err != nil {
		return vxdr.VeloOp{}, errors.Wrap(err, "failed to build XDR operation body")
	}

	return vxdr.VeloOp{Body: body}, nil
}

func (whiteList *WhiteList) FromXDR(vXdrOp vxdr.VeloOp) error {
	whiteListOp := vXdrOp.Body.WhiteListOp
	if whiteListOp == nil {
		return errors.New("error parsing white list operation from xdr")
	}

	whiteList.Role = string(whiteListOp.Role)
	whiteList.Address = whiteListOp.Address.Address()

	return nil
}

func (whiteList *WhiteList) Validate() error {
	if whiteList.Address == "" {
		return errors.New("address parameter cannot be blank")
	}

	if whiteList.Role == "" {
		return errors.New("role parameter cannot be blank")
	}

	if !strkey.IsValidEd25519PublicKey(whiteList.Address) {
		return errors.Errorf("%s is not a valid stellar public key", whiteList.Address)
	}

	if !vxdr.Role(whiteList.Role).IsValid() {
		return errors.New("role specified does not exist")
	}

	return nil
}

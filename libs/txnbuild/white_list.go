package vtxnbuild

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/strkey"
	"gitlab.com/velo-labs/cen/libs/xdr"
)

type WhiteList struct {
	Address  string
	Role     string
	Currency string
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
	if whiteList.Currency != "" {
		vXdrOp.Currency = vxdr.Currency(whiteList.Currency)
	}

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
	whiteList.Currency = string(whiteListOp.Currency)

	return nil
}

func (whiteList *WhiteList) Validate() error {
	if whiteList.Address == "" {
		return errors.New("address must not be blank")
	}

	if whiteList.Role == "" {
		return errors.New("role must not be blank")
	}

	if !strkey.IsValidEd25519PublicKey(whiteList.Address) {
		return errors.New("invalid address format")
	}

	if !vxdr.Role(whiteList.Role).IsValid() {
		return errors.New("role specified does not exist")
	}

	if whiteList.Currency != "" {
		if !vxdr.Currency(whiteList.Currency).IsValid() {
			return errors.Errorf("currency %s does not exist", whiteList.Currency)
		}
	}

	return nil
}

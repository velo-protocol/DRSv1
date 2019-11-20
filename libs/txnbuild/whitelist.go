package vtxnbuild

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/strkey"
	"github.com/velo-protocol/DRSv1/libs/xdr"
)

type Whitelist struct {
	Address  string
	Role     string
	Currency string
}

func (whitelist *Whitelist) BuildXDR() (vxdr.VeloOp, error) {
	if err := whitelist.Validate(); err != nil {
		return vxdr.VeloOp{}, err
	}

	var vXdrOp vxdr.WhitelistOp
	err := vXdrOp.Address.SetAddress(whitelist.Address)
	if err != nil {
		return vxdr.VeloOp{}, errors.Wrap(err, "failed to set whitelist address")
	}

	vXdrOp.Role = vxdr.Role(whitelist.Role)
	if whitelist.Currency != "" {
		vXdrOp.Currency = vxdr.Currency(whitelist.Currency)
	}

	body, err := vxdr.NewOperationBody(vxdr.OperationTypeWhitelist, vXdrOp)
	if err != nil {
		return vxdr.VeloOp{}, errors.Wrap(err, "failed to build XDR operation body")
	}

	return vxdr.VeloOp{Body: body}, nil
}

func (whitelist *Whitelist) FromXDR(vXdrOp vxdr.VeloOp) error {
	whitelistOp := vXdrOp.Body.WhitelistOp
	if whitelistOp == nil {
		return errors.New("error parsing whitelist operation from xdr")
	}

	whitelist.Role = string(whitelistOp.Role)
	whitelist.Address = whitelistOp.Address.Address()
	whitelist.Currency = string(whitelistOp.Currency)

	return nil
}

func (whitelist *Whitelist) Validate() error {
	if whitelist.Address == "" {
		return errors.New("address must not be blank")
	}

	if whitelist.Role == "" {
		return errors.New("role must not be blank")
	}

	if !strkey.IsValidEd25519PublicKey(whitelist.Address) {
		return errors.New("invalid address format")
	}

	if !vxdr.Role(whitelist.Role).IsValid() {
		return errors.New("role specified does not exist")
	}

	if whitelist.Currency != "" && !vxdr.Currency(whitelist.Currency).IsValid() {
		return errors.Errorf("currency %s does not exist", whitelist.Currency)
	}

	return nil
}

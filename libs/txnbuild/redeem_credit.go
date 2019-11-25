package vtxnbuild

import (
	"github.com/pkg/errors"
	_amount "github.com/stellar/go/amount"
	"github.com/stellar/go/protocols/horizon"
	"github.com/velo-protocol/DRSv1/libs/xdr"
	"regexp"
)

// RedeemCredit represents the Velo redeem credit Operation.
type RedeemCredit struct {
	AssetCode string
	Issuer    string
	Amount    string
}

// BuildXDR for RedeemCredit returns a fully configured XDR Operation.
func (redeemCredit *RedeemCredit) BuildXDR() (vxdr.VeloOp, error) {
	if err := redeemCredit.Validate(); err != nil {
		return vxdr.VeloOp{}, err
	}

	amount, err := _amount.Parse(redeemCredit.Amount)
	if err != nil {
		return vxdr.VeloOp{}, errors.Wrap(err, "failed to parse amount")
	}

	// xdr op
	vXdrOp := vxdr.RedeemCreditOp{
		AssetCode: redeemCredit.AssetCode,
		Amount:    amount,
	}
	err = vXdrOp.Issuer.SetAddress(redeemCredit.Issuer)
	if err != nil {
		return vxdr.VeloOp{}, errors.Wrap(err, "failed to set redeem credit issuer address")
	}

	body, err := vxdr.NewOperationBody(vxdr.OperationTypeRedeemCredit, vXdrOp)
	if err != nil {
		return vxdr.VeloOp{}, errors.Wrap(err, "failed to build XDR operation body")
	}

	return vxdr.VeloOp{Body: body}, nil
}

// FromXDR for RedeemCredit initialises the vtxnbuild struct from the corresponding xdr Operation.
func (redeemCredit *RedeemCredit) FromXDR(vXdrOp vxdr.VeloOp) error {
	redeemOp := vXdrOp.Body.RedeemCreditOp
	if redeemOp == nil {
		return errors.New("error parsing redeemCredit operation from xdr")
	}

	redeemCredit.AssetCode = redeemOp.AssetCode
	redeemCredit.Amount = _amount.String(redeemOp.Amount)
	redeemCredit.Issuer = redeemOp.Issuer.Address()

	return nil
}

// Validate for RedeemCredit validates the required struct fields. It returns an error if any of the fields are
// invalid. Otherwise, it returns nil.
func (redeemCredit *RedeemCredit) Validate() error {
	if redeemCredit.AssetCode == "" {
		return errors.New("assetCode must not be blank")
	}
	if redeemCredit.Issuer == "" {
		return errors.New("issuer must not be blank")
	}
	if redeemCredit.Amount == "" {
		return errors.New("amount must not be blank")
	}

	if matched, _ := regexp.MatchString(`^[A-Za-z0-9]{1,7}$`, redeemCredit.AssetCode); !matched {
		return errors.New("invalid assetCode format")
	}

	_, err := horizon.KeyTypeFromAddress(redeemCredit.Issuer)
	if err != nil {
		return errors.New("invalid issuer format")
	}

	amount, err := _amount.Parse(redeemCredit.Amount)
	if err != nil {
		return errors.New("invalid amount format")
	}
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	return nil
}

package vtxnbuild

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/amount"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"regexp"
)

type SetupCredit struct {
	PeggedValue    string
	PeggedCurrency string
	AssetCode      string
}

func (setupCredit *SetupCredit) BuildXDR() (vxdr.VeloOp, error) {
	if err := setupCredit.Validate(); err != nil {
		return vxdr.VeloOp{}, err
	}

	peggedValue, err := amount.Parse(setupCredit.PeggedValue)
	if err != nil {
		return vxdr.VeloOp{}, errors.Wrap(err, "failed to parse pegged value")
	}

	// xdr op
	vXdrOp := vxdr.SetupCreditOp{
		AssetCode:      setupCredit.AssetCode,
		PeggedValue:    peggedValue,
		PeggedCurrency: setupCredit.PeggedCurrency,
	}

	body, err := vxdr.NewOperationBody(vxdr.OperationTypeSetupCredit, vXdrOp)
	if err != nil {
		return vxdr.VeloOp{}, errors.Wrap(err, "failed to build XDR operation body")
	}

	return vxdr.VeloOp{Body: body}, nil
}

func (setupCredit *SetupCredit) FromXDR(vXdrOp vxdr.VeloOp) error {
	setupCreditOp := vXdrOp.Body.SetupCreditOp
	if setupCreditOp == nil {
		return errors.New("error parsing setupCredit operation from xdr")
	}

	setupCredit.PeggedValue = amount.String(setupCreditOp.PeggedValue)
	setupCredit.PeggedCurrency = setupCreditOp.PeggedCurrency
	setupCredit.AssetCode = setupCreditOp.AssetCode

	return nil
}

func (setupCredit *SetupCredit) Validate() error {
	if setupCredit.AssetCode == "" {
		return errors.New("assetCode parameter cannot be blank")
	}
	if setupCredit.PeggedValue == "" {
		return errors.New("peggedValue parameter cannot be blank")
	}
	if setupCredit.PeggedCurrency == "" {
		return errors.New("peggedCurrency parameter cannot be blank")
	}

	peggedValue, err := amount.Parse(setupCredit.PeggedValue)
	if err != nil {
		return errors.New("peggedValue parameter is not a number")
	}
	if peggedValue <= 0 {
		return errors.New("peggedValue must be greater than zero")
	}

	if matched, _ := regexp.MatchString(`^[A-Za-z0-9]{1,12}$`, setupCredit.AssetCode); !matched {
		return errors.New("invalid format of asset code")
	}

	if !vxdr.Currency(setupCredit.PeggedCurrency).IsValid() {
		return errors.Errorf("the pegged currency %s does not exist", setupCredit.PeggedCurrency)
	}

	return nil
}

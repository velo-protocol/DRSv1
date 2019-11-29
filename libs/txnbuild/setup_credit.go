package vtxnbuild

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/amount"
	"github.com/velo-protocol/DRSv1/libs/xdr"
	"regexp"
)

// SetupCredit represents the Velo setup credit Operation.
type SetupCredit struct {
	PeggedValue    string
	PeggedCurrency string
	AssetCode      string
}

// BuildXDR for SetupCredit returns a fully configured XDR Operation.
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

// FromXDR for SetupCredit initialises the vtxnbuild struct from the corresponding XDR Operation.
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

// Validation function for SetupCredit. Validates the required struct fields. It returns an error if any of the fields are
// invalid. Otherwise, it returns nil.
func (setupCredit *SetupCredit) Validate() error {
	if setupCredit.AssetCode == "" {
		return errors.New("assetCode must not be blank")
	}
	if setupCredit.PeggedValue == "" {
		return errors.New("peggedValue must not be blank")
	}
	if setupCredit.PeggedCurrency == "" {
		return errors.New("peggedCurrency must not be blank")
	}

	peggedValue, err := amount.Parse(setupCredit.PeggedValue)
	if err != nil {
		return errors.New("invalid peggedValue format")
	}
	if peggedValue <= 0 {
		return errors.New("peggedValue must be greater than zero")
	}

	if matched, _ := regexp.MatchString(`^[A-Za-z0-9]{1,7}$`, setupCredit.AssetCode); !matched {
		return errors.New("invalid assetCode format")
	}

	if !vxdr.Currency(setupCredit.PeggedCurrency).IsValid() {
		return errors.Errorf("peggedCurrency %s does not exist", setupCredit.PeggedCurrency)
	}

	return nil
}

package vtxnbuild

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/amount"
	"github.com/velo-protocol/DRSv1/libs/xdr"
)

// PriceUpdate represents the Velo price update Operation.
type PriceUpdate struct {
	Asset                       string
	Currency                    string
	PriceInCurrencyPerAssetUnit string
}

// BuildXDR for PriceUpdate returns a fully configured XDR Operation.
func (priceUpdate *PriceUpdate) BuildXDR() (vxdr.VeloOp, error) {
	if err := priceUpdate.Validate(); err != nil {
		return vxdr.VeloOp{}, err
	}

	priceInCurrencyPerAssetUnit, err := amount.Parse(priceUpdate.PriceInCurrencyPerAssetUnit)
	if err != nil {
		return vxdr.VeloOp{}, errors.Wrap(err, "failed to parse priceInCurrencyPerAssetUnit")
	}

	// xdr op
	vXdrOp := vxdr.PriceUpdateOp{
		Asset:                       priceUpdate.Asset,
		Currency:                    vxdr.Currency(priceUpdate.Currency),
		PriceInCurrencyPerAssetUnit: priceInCurrencyPerAssetUnit,
	}

	body, err := vxdr.NewOperationBody(vxdr.OperationTypePriceUpdate, vXdrOp)
	if err != nil {
		return vxdr.VeloOp{}, errors.Wrap(err, "failed to build XDR operation body")
	}

	return vxdr.VeloOp{Body: body}, nil
}

// FromXDR for PriceUpdate initialises the vtxnbuild struct from the corresponding xdr Operation.
func (priceUpdate *PriceUpdate) FromXDR(vXdrOp vxdr.VeloOp) error {
	priceUpdateOp := vXdrOp.Body.PriceUpdateOp
	if priceUpdateOp == nil {
		return errors.New("error parsing priceUpdate operation from xdr")
	}

	priceUpdate.Asset = priceUpdateOp.Asset
	priceUpdate.Currency = string(priceUpdateOp.Currency)
	priceUpdate.PriceInCurrencyPerAssetUnit = amount.String(priceUpdateOp.PriceInCurrencyPerAssetUnit)

	return nil
}

// Validate for PriceUpdate validates the required struct fields. It returns an error if any of the fields are
// invalid. Otherwise, it returns nil.
func (priceUpdate *PriceUpdate) Validate() error {
	if priceUpdate.Asset == "" {
		return errors.New("asset must not be blank")
	}
	if priceUpdate.Currency == "" {
		return errors.New("currency must not be blank")
	}
	if priceUpdate.PriceInCurrencyPerAssetUnit == "" {
		return errors.New("priceInCurrencyPerAssetUnit must not be blank")
	}

	priceInCurrencyPerAssetUnit, err := amount.Parse(priceUpdate.PriceInCurrencyPerAssetUnit)
	if err != nil {
		return errors.New("invalid priceInCurrencyPerAssetUnit format")
	}
	if priceInCurrencyPerAssetUnit <= 0 {
		return errors.New("priceInCurrencyPerAssetUnit must be greater than zero")
	}

	if !vxdr.Currency(priceUpdate.Currency).IsValid() {
		return errors.Errorf("the currency %s does not exist", priceUpdate.Currency)
	}

	if !vxdr.Asset(priceUpdate.Asset).IsValid() {
		return errors.Errorf("asset %s does not exist", priceUpdate.Asset)
	}

	return nil
}

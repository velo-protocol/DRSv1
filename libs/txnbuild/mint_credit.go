package vtxnbuild

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/amount"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"regexp"
)

type MintCredit struct {
	AssetCodeToBeIssued string
	CollateralAssetCode string
	CollateralAmount    string
}

func (mintCredit *MintCredit) BuildXDR() (vxdr.VeloOp, error) {
	if err := mintCredit.Validate(); err != nil {
		return vxdr.VeloOp{}, err
	}

	collateralAmount, err := amount.Parse(mintCredit.CollateralAmount)
	if err != nil {
		return vxdr.VeloOp{}, errors.Wrap(err, "failed to parse collateralAmount")
	}

	// xdr op
	vXdrOp := vxdr.MintCreditOp{
		AssetCodeToBeIssued: mintCredit.AssetCodeToBeIssued,
		CollateralAssetCode: vxdr.Asset(mintCredit.CollateralAssetCode),
		CollateralAmount:    collateralAmount,
	}

	body, err := vxdr.NewOperationBody(vxdr.OperationTypeMintCredit, vXdrOp)
	if err != nil {
		return vxdr.VeloOp{}, errors.Wrap(err, "failed to build XDR operation body")
	}

	return vxdr.VeloOp{Body: body}, nil
}

func (mintCredit *MintCredit) FromXDR(vXdrOp vxdr.VeloOp) error {
	mintCreditOp := vXdrOp.Body.MintCreditOp
	if mintCreditOp == nil {
		return errors.New("error parsing mintCredit operation from xdr")
	}

	mintCredit.AssetCodeToBeIssued = mintCreditOp.AssetCodeToBeIssued
	mintCredit.CollateralAssetCode = string(mintCreditOp.CollateralAssetCode)
	mintCredit.CollateralAmount = amount.String(mintCreditOp.CollateralAmount)
	return nil
}

func (mintCredit *MintCredit) Validate() error {
	if mintCredit.AssetCodeToBeIssued == "" {
		return errors.New("assetCodeToBeIssued must not be blank")
	}
	if mintCredit.CollateralAssetCode == "" {
		return errors.New("collateralAssetCode must not be blank")
	}
	if mintCredit.CollateralAmount == "" {
		return errors.New("collateralAmount must not be blank")
	}

	if matched, _ := regexp.MatchString(`^[A-Za-z0-9]{1,7}$`, mintCredit.AssetCodeToBeIssued); !matched {
		return errors.New("invalid assetCode format")
	}

	if !vxdr.Asset(mintCredit.CollateralAssetCode).IsValid() {
		return errors.New("collateralAssetCode value must be VELO")
	}

	collateralAmount, err := amount.Parse(mintCredit.CollateralAmount)
	if err != nil {
		return errors.New("invalid collateralAmount format")
	}
	if collateralAmount <= 0 {
		return errors.New("collateralAmount must be greater than zero")
	}

	return nil
}

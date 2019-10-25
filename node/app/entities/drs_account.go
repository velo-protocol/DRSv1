package entities

import (
	"github.com/pkg/errors"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/utils"
)

type DrsAccountData struct {
	DrsReserve                string `mapstructure:"DrsReserve"`
	TrustedPartnerListAddress string `mapstructure:"TrustedPartnerList"`
	RegulatorListAddress      string `mapstructure:"RegulatorList"`
	PriceFeederListAddress    string `mapstructure:"PriceFeederList"`
	PriceUsdVeloAddress       string `mapstructure:"Price[USD-VELO]"`
	PriceThbVeloAddress       string `mapstructure:"Price[THB-VELO]"`
	PriceSgdVeloAddress       string `mapstructure:"Price[SGD-VELO]"`
	Base64Decoded             bool
}

func (drsAccountData *DrsAccountData) VeloPriceAddress(currency vxdr.Currency) string {
	switch currency {
	case vxdr.CurrencyTHB:
		return drsAccountData.PriceThbVeloAddress
	case vxdr.CurrencyUSD:
		return drsAccountData.PriceUsdVeloAddress
	case vxdr.CurrencySGD:
		return drsAccountData.PriceSgdVeloAddress
	default:
		return ""
	}
}

func (drsAccountData *DrsAccountData) RoleListAddress(role vxdr.Role) string {
	switch role {
	case vxdr.RoleRegulator:
		return drsAccountData.RegulatorListAddress
	case vxdr.RolePriceFeeder:
		return drsAccountData.PriceFeederListAddress
	case vxdr.RoleTrustedPartner:
		return drsAccountData.TrustedPartnerListAddress
	default:
		return ""
	}
}

func (drsAccountData *DrsAccountData) DecodeBase64() error {
	if drsAccountData.Base64Decoded {
		return nil
	}

	var err error
	drsAccountData.DrsReserve, err = utils.DecodeBase64(drsAccountData.DrsReserve)
	if err != nil {
		return errors.Wrap(err, "fail to decode drs reserve address")
	}

	drsAccountData.TrustedPartnerListAddress, err = utils.DecodeBase64(drsAccountData.TrustedPartnerListAddress)
	if err != nil {
		return errors.Wrap(err, "fail to decode trusted partner list address")
	}
	drsAccountData.PriceFeederListAddress, err = utils.DecodeBase64(drsAccountData.PriceFeederListAddress)
	if err != nil {
		return errors.Wrap(err, "fail to decode price feeder list address")
	}
	drsAccountData.RegulatorListAddress, err = utils.DecodeBase64(drsAccountData.RegulatorListAddress)
	if err != nil {
		return errors.Wrap(err, "fail to decode regulator list address")
	}
	drsAccountData.PriceSgdVeloAddress, err = utils.DecodeBase64(drsAccountData.PriceSgdVeloAddress)
	if err != nil {
		return errors.Wrap(err, "fail to decode price sgd-velo address")
	}
	drsAccountData.PriceUsdVeloAddress, err = utils.DecodeBase64(drsAccountData.PriceUsdVeloAddress)
	if err != nil {
		return errors.Wrap(err, "fail to decode price usd-velo address")
	}
	drsAccountData.PriceThbVeloAddress, err = utils.DecodeBase64(drsAccountData.PriceThbVeloAddress)
	if err != nil {
		return errors.Wrap(err, "fail to decode price thb-velo address")
	}

	drsAccountData.Base64Decoded = true
	return nil
}

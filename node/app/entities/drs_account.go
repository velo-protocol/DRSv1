package entities

type DrsAccountData struct {
	TrustedPartnerListAddress string `mapstructure:"TrustedPartnerList"`
	RegulatorListAddress      string `mapstructure:"RegulatorList"`
	PriceFeederListAddress    string `mapstructure:"PriceFeederList"`
	PriceUsdVeloAddress       string `mapstructure:"Price[USD-VELO]"`
	PriceThbVeloAddress       string `mapstructure:"Price[THB-VELO]"`
	PriceSgdVeloAddress       string `mapstructure:"Price[SGD-VELO]"`
}

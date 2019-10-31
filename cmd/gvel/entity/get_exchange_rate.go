package entity

type GetExchangeRateInput struct {
	AssetCode string
	Issuer    string
}

type GetExchangeRateOutput struct {
	AssetCode              string
	Issuer                 string
	RedeemablePricePerUnit string
	RedeemableCollateral   string
}

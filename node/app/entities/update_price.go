package entities

type UpdatePriceOutput struct {
	SignedStellarTxXdr          string
	Asset                       string
	Currency                    string
	PriceInCurrencyPerAssetUnit string
}

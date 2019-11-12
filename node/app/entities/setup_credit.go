package entities

type SetupCreditOutput struct {
	SignedStellarTxXdr string
	AssetIssuer        string
	AssetDistributor   string
	AssetCode          string
	PeggedValue        string
	PeggedCurrency     string
}

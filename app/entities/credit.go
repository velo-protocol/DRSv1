package entities

type Credit struct {
	CreditOwnerAddress   string
	IssuerCreationTxHash string
	SetupTxHash          string
	IssuerAddress        string
	DistributorAddress   string
	AssetName            string
	PeggedValue          string
	PeggedCurrency       string
}

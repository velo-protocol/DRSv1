package entities

type RebalanceOutput struct {
	Collaterals        []Collateral
	SignedStellarTxXdr string
}

type Collateral struct {
	AssetCode      string
	AssetIssuer    string
	RequiredAmount string
	PoolAmount     string
}

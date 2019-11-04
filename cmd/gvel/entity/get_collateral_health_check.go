package entity

type GetCollateralHealthCheckOutput struct {
	AssetCode      string
	AssetIssuer    string
	RequiredAmount string
	PoolAmount     string
}

package entities

type GetAssetInput struct {
	AssetCode   string
	AssetIssuer string
	// Optional cursor default null, A paging token, specifying where to start returning records from.
	Cursor *string
	// Optional order default asc, The order in which to return rows, "asc" or "desc", ordered by asset_code then by asset_issuer.
	Order *string
	// Optional limit default 10, Maximum number of records to return 200.
	Limit *uint
}

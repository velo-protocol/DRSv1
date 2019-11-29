package vxdr

// Asset is a constant which defined supported collateral asset.
type Asset string

const AssetVELO Asset = "VELO"

// IsValid checks if the given asset is supported/valid or not.
func (asset Asset) IsValid() bool {
	return asset == AssetVELO
}

package vxdr

type Asset string

const AssetVELO Asset = "VELO"

func (asset Asset) IsValid() bool {
	return asset == AssetVELO
}

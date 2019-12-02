package constants

import (
	"github.com/stellar/go/amount"
	"math"
)

var MaxTrustlineLimit = amount.StringFromInt64(math.MaxInt64)

const (
	AssetCode = "assetCode"
	Issuer    = "issuer"
)

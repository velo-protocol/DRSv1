package env

import (
	"fmt"
	"syscall"
)

var (
	DrsAddress        string
	DrsPrivkey        string
	VeloIssuerAddress string
	NetworkPassphrase string
	LevelDBPath       string
	HorizonUrl        string
	CreditPrefix      string
	PriceFeedPrefix   string
	PricePrefix       string
)

// Init Initialize env variables
func Init() {
	DrsAddress = requireEnv("DRS_ADDRESS")
	DrsPrivkey = requireEnv("DRS_PRIVKEY")
	VeloIssuerAddress = requireEnv("VELO_ISSUER_ADDRESS")
	NetworkPassphrase = requireEnv("NETWORK_PASSPHRASE")
	LevelDBPath = requireEnv("LEVEL_DB_PATH")
	HorizonUrl = requireEnv("HORIZON_URL")
	CreditPrefix = requireEnv("CREDIT_PREFIX")
	PriceFeedPrefix = requireEnv("PRICE_FEED_PREFIX")
	PricePrefix = requireEnv("PRICE_PREFIX")
}

func requireEnv(envName string) string {
	value, found := syscall.Getenv(envName)

	if !found {
		panic(fmt.Sprintf("%s env is required", envName))
	}

	return value
}

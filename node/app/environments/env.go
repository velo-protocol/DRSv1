package env

import (
	"fmt"
	"syscall"
)

var (
	DrsAddress        string
	DrsPrivateKey     string
	VeloIssuerAddress string
	NetworkPassphrase string
	HorizonUrl        string
	CreditPrefix      string
	PriceFeedPrefix   string
	PricePrefix       string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	Port string
)

// Init Initialize env variables
func Init() {
	DrsAddress = requireEnv("DRS_ADDRESS")
	DrsPrivateKey = requireEnv("DRS_PRIVKEY")
	VeloIssuerAddress = requireEnv("VELO_ISSUER_ADDRESS")
	NetworkPassphrase = requireEnv("NETWORK_PASSPHRASE")
	HorizonUrl = requireEnv("HORIZON_URL")
	CreditPrefix = requireEnv("CREDIT_PREFIX")
	PriceFeedPrefix = requireEnv("PRICE_FEED_PREFIX")
	PricePrefix = requireEnv("PRICE_PREFIX")

	PostgresHost = requireEnv("POSTGRES_HOST")
	PostgresPort = requireEnv("POSTGRES_PORT")
	PostgresUser = requireEnv("POSTGRES_USER")
	PostgresPassword = requireEnv("POSTGRES_PASSWORD")
	PostgresDB = requireEnv("POSTGRES_DB")

	Port = "8080"
}

func requireEnv(envName string) string {
	value, found := syscall.Getenv(envName)

	if !found {
		panic(fmt.Sprintf("%s env is required", envName))
	}

	return value
}

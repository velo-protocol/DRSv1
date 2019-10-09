package env

import (
	"fmt"
	"syscall"
	"time"
)

var (
	DrsPublicKey                string
	DrsSecretKey                string
	VeloIssuerPublicKey         string
	NetworkPassphrase           string
	HorizonURL                  string
	StellarTxTimeBoundInMinutes int64
	ValidPriceBoundary          time.Duration

	Port                int64
	HealthPort          int64
	EnableReflectionAPI bool
)

// Init Initialize env variables
func Init() {
	DrsPublicKey = requireEnv("DRS_PUBLIC_KEY")
	DrsSecretKey = requireEnv("DRS_SECRET_KEY")
	VeloIssuerPublicKey = requireEnv("VELO_ISSUER_PUBLIC_KEY")
	NetworkPassphrase = requireEnv("NETWORK_PASSPHRASE")
	HorizonURL = requireEnv("HORIZON_URL")
	StellarTxTimeBoundInMinutes = 15 * 60 // 15 minute * 60 seconds
	ValidPriceBoundary = 15 * time.Minute

	Port = 8080
	HealthPort = 6666
	EnableReflectionAPI = true
}

func requireEnv(envName string) string {
	value, found := syscall.Getenv(envName)

	if !found {
		panic(fmt.Sprintf("%s env is required", envName))
	}

	return value
}

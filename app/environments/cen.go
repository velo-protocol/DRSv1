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
	HorizonUrl string
)

// Init Initialize env variables
func Init() {
	DrsAddress = requireEnv("DRS_ADDRESS")
	DrsPrivkey = requireEnv("DRS_PRIVKEY")
	VeloIssuerAddress = requireEnv("VELO_ISSUER_ADDRESS")
	NetworkPassphrase = requireEnv("NETWORK_PASSPHRASE")
	LevelDBPath = requireEnv("LEVEL_DB_PATH")
	HorizonUrl = requireEnv("HORIZON_URL")
}

func requireEnv(envName string) string {
	value, found := syscall.Getenv(envName)

	if !found {
		panic(fmt.Sprintf("%s env is required", envName))
	}

	return value
}

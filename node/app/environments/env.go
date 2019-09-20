package env

import (
	"fmt"
	"syscall"
)

var (
	DrsPublicKey        string
	DrsPrivateKey       string
	VeloIssuerPublicKey string
	NetworkPassphrase   string
	HorizonURL          string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	Port                string
	EnableReflectionAPI bool
)

// Init Initialize env variables
func Init() {
	DrsPublicKey = requireEnv("DRS_PUBLIC_KEY")
	DrsPrivateKey = requireEnv("DRS_SECRET_KEY")
	VeloIssuerPublicKey = requireEnv("VELO_ISSUER_PUBLIC_KEY")
	NetworkPassphrase = requireEnv("NETWORK_PASSPHRASE")
	HorizonURL = requireEnv("HORIZON_URL")

	PostgresHost = requireEnv("POSTGRES_HOST")
	PostgresPort = requireEnv("POSTGRES_PORT")
	PostgresUser = requireEnv("POSTGRES_USER")
	PostgresPassword = requireEnv("POSTGRES_PASSWORD")
	PostgresDB = requireEnv("POSTGRES_DB")

	Port = "8080"
	EnableReflectionAPI = true
}

func requireEnv(envName string) string {
	value, found := syscall.Getenv(envName)

	if !found {
		panic(fmt.Sprintf("%s env is required", envName))
	}

	return value
}

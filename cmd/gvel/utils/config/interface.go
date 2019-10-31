package config

type Configuration interface {
	InitConfigFile(configFilePath string) error
	Exists() bool
	SetDefaultAccount(account string) error
	GetDefaultAccount() string
	GetAccountDbPath() string
	GetHorizonUrl() string
	GetVeloNodeUrl() string
	GetNetworkPassphrase() string
	GetIsTestNet() bool
}

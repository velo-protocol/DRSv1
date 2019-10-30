package config

type Configuration interface {
	InitConfigFile(configFilePath string) error
	Exists() bool
	SetDefaultAccount(account string) error
	GetDefaultAccount() string
	GetAccountDbPath() string
	GetFriendBotUrl() string
	GetHorizonUrl() string
	GetVeloNodeUrl() string
	GetNetworkPassphrase() string
}

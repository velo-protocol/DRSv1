package config

type Configuration interface {
	InitConfigFile(configFilePath string) error
	Exists() bool
	GetDefaultAccount() string
	SetDefaultAccount(account string) error
	GetAccountDbPath() string
}

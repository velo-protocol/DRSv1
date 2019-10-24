package config

type Configuration interface {
	GetDefaultAccount() string
	SetDefaultAccount(account string) error
}

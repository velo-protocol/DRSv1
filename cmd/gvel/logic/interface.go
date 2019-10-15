package logic

import "github.com/stellar/go/keypair"

type Logic interface {
	Init(configFilePath string) error
	CreateAccount(passphrase string) (*keypair.Full, error)
}

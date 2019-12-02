package entity

import "github.com/stellar/go/keypair"

type CreateAccountInput struct {
	Passphrase          string
	SetAsDefaultAccount bool
}

type CreateAccountOutput struct {
	GeneratedKeyPair *keypair.Full
	IsDefault        bool
}

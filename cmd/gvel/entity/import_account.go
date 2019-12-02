package entity

import "github.com/stellar/go/keypair"

type ImportAccountInput struct {
	Passphrase   string
	SeedKey      string
	SetAsDefault bool
}

type ImportAccountOutput struct {
	ImportedKeyPair *keypair.Full
	IsDefault       bool
}

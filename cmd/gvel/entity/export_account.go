package entity

import "github.com/stellar/go/keypair"

type ExportAccountInput struct {
	PublicKey  string
	Passphrase string
}

type ExportAccountOutput struct {
	ExportedKeyPair *keypair.Full
}

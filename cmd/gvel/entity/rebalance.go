package entity

import "github.com/stellar/go/protocols/horizon"

type RebalanceInput struct {
	Passphrase string
}

type RebalanceOutput struct {
	TxResult *horizon.TransactionSuccess
}

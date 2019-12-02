package stellar

import "github.com/stellar/go/protocols/horizon"

type Repository interface {
	GetFreeLumens(stellarAddress string) error
	GetStellarAccount(stellarAddress string) (*horizon.Account, error)
}

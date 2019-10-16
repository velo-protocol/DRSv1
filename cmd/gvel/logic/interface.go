package logic

import (
	"github.com/stellar/go/keypair"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
)

type Logic interface {
	Init(configFilePath string) error
	CreateAccount(passphrase string) (*keypair.Full, error)
	ListAccount() (*[]entity.StellarAccount, error)
}

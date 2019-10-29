package logic

import (
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
)

type Logic interface {
	Init(configFilePath string) error
	CreateAccount(input *entity.CreateAccountInput) (*entity.CreateAccountOutput, error)
	ListAccount() (*[]entity.StellarAccount, error)
	SetDefaultAccount(input *entity.SetDefaultAccountInput) (*entity.SetDefaultAccountOutput, error)
}

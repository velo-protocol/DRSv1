package logic

import (
	"github.com/velo-protocol/DRSv1/cmd/gvel/entity"
)

type Logic interface {
	Init(configFilePath string) error
	CreateAccount(input *entity.CreateAccountInput) (*entity.CreateAccountOutput, error)
	ListAccount() (*[]entity.StellarAccount, error)
	SetDefaultAccount(input *entity.SetDefaultAccountInput) (*entity.SetDefaultAccountOutput, error)
	ImportAccount(input *entity.ImportAccountInput) (*entity.ImportAccountOutput, error)
	ExportAccount(input *entity.ExportAccountInput) (*entity.ExportAccountOutput, error)

	SetupCredit(input *entity.SetupCreditInput) (*entity.SetupCreditOutput, error)
	MintCredit(input *entity.MintCreditInput) (*entity.MintCreditOutput, error)
	RedeemCredit(input *entity.RedeemCreditInput) (*entity.RedeemCreditOutput, error)

	GetExchangeRate(input *entity.GetExchangeRateInput) (*entity.GetExchangeRateOutput, error)
	GetCollateralHealthCheck() (*entity.GetCollateralHealthCheckOutput, error)
	RebalanceReserve(input *entity.RebalanceInput) (*entity.RebalanceOutput, error)
}

package test_helpers

import (
	"gitlab.com/velo-labs/cen/app/entities"
)

func GetCreditEntity() entities.Credit {
	return entities.Credit{
		CreditOwnerAddress: GetRandStellarAccount(),
	}
}

package test_helpers

import (
	"gitlab.com/velo-labs/cen/node/app/entities"
)

func GetCreditEntity() entities.Credit {
	return entities.Credit{
		CreditOwnerAddress: GetRandStellarAccount().Address(),
	}
}

package logic

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
)

func (lo *logic) ListAccount() (*[]entity.StellarAccount, error) {
	accountsBytes, err := lo.DB.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account from db")
	}

	var accounts []entity.StellarAccount
	for _, accountBytes := range accountsBytes {
		var tmpAccount entity.StellarAccount

		err := json.Unmarshal(accountBytes, &tmpAccount)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal account")
		}

		accounts = append(accounts, tmpAccount)
	}

	return &accounts, nil
}

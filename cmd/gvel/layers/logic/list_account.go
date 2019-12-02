package logic

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/velo-protocol/DRSv1/cmd/gvel/entity"
)

func (lo *logic) ListAccount() (*[]entity.StellarAccount, error) {
	accountsBytes, err := lo.DB.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account from db")
	}

	defaultAccount := lo.AppConfig.GetDefaultAccount()

	var accounts []entity.StellarAccount
	for _, accountBytes := range accountsBytes {
		var tmpAccount entity.StellarAccount

		err := json.Unmarshal(accountBytes, &tmpAccount)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal account")
		}

		tmpAccount.IsDefault = defaultAccount == tmpAccount.Address

		accounts = append(accounts, tmpAccount)
	}

	return &accounts, nil
}

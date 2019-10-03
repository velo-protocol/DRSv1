package stellar

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/utils"
)

func (repo *repo) GetAccountDecodedDataByKey(stellarAddress string, key string) (string, error) {
	account, err := repo.HorizonClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: stellarAddress,
	})
	if err != nil {
		return "", errors.Wrapf(err, constants.ErrGetAccountDetail, stellarAddress)
	}

	value, err := utils.DecodeBase64(account.Data[key])
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf(constants.ErrToDecodeData, key))
	}
	return value, nil
}

package stellar

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/velo-protocol/DRSv1/node/app/constants"
	"github.com/velo-protocol/DRSv1/node/app/utils"
)

func (repo *repo) GetAccountDecodedData(stellarAddress string) (map[string]string, error) {
	account, err := repo.HorizonClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: stellarAddress,
	})
	if err != nil {
		return nil, errors.Wrapf(err, constants.ErrGetAccountDetail, stellarAddress)
	}

	for key, encodedValue := range account.Data {
		account.Data[key], err = utils.DecodeBase64(encodedValue)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf(constants.ErrToDecodeData, key))
		}
	}
	return account.Data, nil
}

package stellar

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"gitlab.com/velo-labs/cen/node/app/utils"
)

func (repo *repo) GetAccountDecodedData(stellarAddress string) (map[string]string, error) {
	account, err := repo.HorizonClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: stellarAddress,
	})
	if err != nil {
		return nil, errors.Wrap(err, "fail to get DRS account details")
	}

	for key, encodedValue := range account.Data {
		account.Data[key], err = utils.DecodeBase64(encodedValue)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf(`fail to decode data "%s"`, key))
		}
	}
	return account.Data, nil
}

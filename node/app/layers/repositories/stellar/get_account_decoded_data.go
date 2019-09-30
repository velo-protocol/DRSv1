package stellar

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
)

func (repo *repo) GetAccountDecodedData(stellarAddress string) (map[string]string, error) {
	account, err := repo.HorizonClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: stellarAddress,
	})
	if err != nil {
		return nil, errors.Wrap(err, "fail to get DRS account details")
	}

	for key, encodedValue := range account.Data {
		decodedValue, err := base64.StdEncoding.DecodeString(encodedValue)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf(`fail to decode data "%s"`, key))
		}
		account.Data[key] = string(decodedValue)
	}
	return account.Data, nil
}

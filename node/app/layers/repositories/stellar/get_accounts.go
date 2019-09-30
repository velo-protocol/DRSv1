package stellar

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"sync"
)

type getAccountsResult struct {
	Account horizon.Account
	Error   error
}

func (repo *repo) GetAccounts(stellarAddresses ...string) ([]horizon.Account, error) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(stellarAddresses))

	getAccountsResults := make([]getAccountsResult, len(stellarAddresses))
	for i, stellarAddress := range stellarAddresses {
		go func(stellarAddress string, i int) {
			defer func() {
				if r := recover(); r != nil {
					err := r.(error)
					getAccountsResults[i] = getAccountsResult{
						Error: errors.Wrapf(err, "fail to get account detail of %s", stellarAddress),
					}
				}
				waitGroup.Done()
			}()

			account, err := repo.HorizonClient.AccountDetail(horizonclient.AccountRequest{
				AccountID: stellarAddress,
			})
			if err != nil {
				panic(err)
			}

			getAccountsResults[i] = getAccountsResult{
				Account: account,
			}
		}(stellarAddress, i)
	}
	waitGroup.Wait()

	accounts := make([]horizon.Account, len(stellarAddresses))
	for i, result := range getAccountsResults {
		if result.Error != nil {
			return nil, result.Error
		}
		accounts[i] = result.Account
	}

	return accounts, nil
}

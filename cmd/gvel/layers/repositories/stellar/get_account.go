package stellar

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
)

func (s *stellar) GetStellarAccount(stellarAddress string) (*horizon.Account, error) {
	account, err := s.HorizonClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: stellarAddress,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "fail to get account detail of %s", stellarAddress)
	}

	return &account, nil
}

package repository

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"gitlab.com/velo-labs/cen/app/modules/stellar"
)

type horizonStellarRepository struct {
	Horizonclient *horizonclient.Client
}

func NewHorizonStellarRepository(client *horizonclient.Client) stellar.Repository {
	return &horizonStellarRepository{
		Horizonclient: client,
	}
}

func (h *horizonStellarRepository) SubmitTransaction(txB64 string) (*horizon.TransactionSuccess, error) {
	txSuccess, err := h.Horizonclient.SubmitTransactionXDR(txB64)
	if err != nil {
		herr, isHorizonError := err.(*horizonclient.Error)
		if !isHorizonError {
			return nil, errors.Wrap(err, "fail to confirm with stellar")
		}
		herrString, _ := herr.ResultString()
		return nil, errors.Wrap(err, fmt.Sprintf("herr result string: %s", herrString))
	}

	return &txSuccess, nil
}

func (h *horizonStellarRepository) LoadAccount(stellarAddress string) (*horizon.Account, error) {
	account, err := h.Horizonclient.AccountDetail(horizonclient.AccountRequest{
		AccountID: stellarAddress,
	})
	if err != nil {
		return nil, errors.Wrap(err, "fail to get account details")
	}

	return &account, nil
}

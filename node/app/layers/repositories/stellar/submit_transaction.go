package stellar

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
)

func (repo *repo) SubmitTransaction(txB64 string) (*horizon.TransactionSuccess, error) {
	txSuccess, err := repo.HorizonClient.SubmitTransactionXDR(txB64)
	if err != nil {
		herr, ok := err.(horizonclient.Error)
		if !ok {
			return nil, errors.Wrap(err, "fail to confirm with stellar")
		}
		herrString, _ := herr.ResultString()
		return nil, errors.Wrap(err, fmt.Sprintf(`horizon err "%s"`, herrString))
	}

	return &txSuccess, nil
}

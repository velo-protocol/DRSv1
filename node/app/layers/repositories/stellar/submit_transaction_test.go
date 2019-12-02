package stellar_test

import (
	"errors"
	"fmt"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/support/render/problem"
	"github.com/stretchr/testify/assert"
	"github.com/velo-protocol/DRSv1/node/app/testhelpers"
	"testing"
)

func TestRepo_SubmitTransaction(t *testing.T) {
	testhelpers.InitEnv()

	t.Run("success", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("SubmitTransactionXDR", "AAAXDRSTRING").
			Return(horizon.TransactionSuccess{}, nil)

		txSuccess, err := helper.repo.SubmitTransaction("AAAXDRSTRING")

		assert.NoError(t, err)
		assert.NotNil(t, txSuccess)

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "SubmitTransactionXDR", 1)
	})

	t.Run("error, fail to confirm with stellar", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("SubmitTransactionXDR", "AAAXDRSTRING").
			Return(horizon.TransactionSuccess{}, errors.New("some error has occurred"))

		_, err := helper.repo.SubmitTransaction("AAAXDRSTRING")

		assert.Error(t, err)
		//assert.Contains(t, err.Error(), fmt.Sprintf("fail to get account detail of %s", testhelpers.PublicKey1))
		assert.Contains(t, err.Error(), "fail to confirm with stellar")

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "SubmitTransactionXDR", 1)
	})

	t.Run("error, horizon reply with error", func(t *testing.T) {
		helper := initTest()

		helper.mockedHorizonClient.
			On("SubmitTransactionXDR", "AAAXDRSTRING").
			Return(horizon.TransactionSuccess{}, horizonclient.Error{
				Response: nil,
				Problem: problem.P{
					Type: "transaction_failed",
					Extras: map[string]interface{}{
						"result_codes": map[string]interface{}{
							"transaction": "tx_failed",
							"operations":  []string{"op_underfunded", "op_already_exists"},
						},
						"result_xdr": "AAARESULTXDRSTRING",
					},
				},
			})

		_, err := helper.repo.SubmitTransaction("AAAXDRSTRING")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf(`horizon err "%s"`, "AAARESULTXDRSTRING"))

		helper.mockedHorizonClient.
			AssertNumberOfCalls(t, "SubmitTransactionXDR", 1)
	})
}

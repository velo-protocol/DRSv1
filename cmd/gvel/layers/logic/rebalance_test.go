package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stretchr/testify/assert"
	"github.com/velo-protocol/DRSv1/cmd/gvel/entity"
	"github.com/velo-protocol/DRSv1/libs/client"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"testing"
)

func TestLogic_RebalanceReserve(t *testing.T) {
	var (
		passPhrase = "password"
	)

	t.Run("happy", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockConfiguration.EXPECT().
			GetDefaultAccount().
			Return(stellarAccountEntity().Address)
		helper.mockDB.EXPECT().
			Get([]byte(stellarAccountEntity().Address)).
			Return(stellarAccountsBytes(), nil)
		helper.mockVelo.EXPECT().
			Client(helper.keyPair).
			Return(helper.mockVeloClient)
		helper.mockVeloClient.EXPECT().
			RebalanceReserve(context.Background(), gomock.AssignableToTypeOf(vtxnbuild.RebalanceReserve{})).
			Return(vclient.RebalanceReserveResult{
				HorizonResult: &horizon.TransactionSuccess{},
			}, nil)

		output, err := helper.logic.RebalanceReserve(&entity.RebalanceInput{
			Passphrase: passPhrase,
		})

		assert.NoError(t, err)
		assert.NotNil(t, output)
	})

	t.Run("error, database returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockConfiguration.EXPECT().
			GetDefaultAccount().
			Return(stellarAccountEntity().Address)
		helper.mockDB.EXPECT().
			Get([]byte(stellarAccountEntity().Address)).
			Return(nil, errors.New("some error has occurred"))

		output, err := helper.logic.RebalanceReserve(&entity.RebalanceInput{
			Passphrase: passPhrase,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "failed to get account from db")
	})

	t.Run("error, bad passphrase", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockConfiguration.EXPECT().
			GetDefaultAccount().
			Return(stellarAccountEntity().Address)
		helper.mockDB.EXPECT().
			Get([]byte(stellarAccountEntity().Address)).
			Return(stellarAccountsBytes(), nil)

		output, err := helper.logic.RebalanceReserve(&entity.RebalanceInput{
			Passphrase: "bad passphrase",
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "failed to decrypt the seed")
	})

	t.Run("error, velo node client returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockConfiguration.EXPECT().
			GetDefaultAccount().
			Return(stellarAccountEntity().Address)
		helper.mockDB.EXPECT().
			Get([]byte(stellarAccountEntity().Address)).
			Return(stellarAccountsBytes(), nil)
		helper.mockVelo.EXPECT().
			Client(helper.keyPair).
			Return(helper.mockVeloClient)
		helper.mockConfiguration.EXPECT().
			GetHorizonUrl().
			Return("https://fake-horizon.com")
		helper.mockConfiguration.EXPECT().
			GetNetworkPassphrase().
			Return("fake-network")
		helper.mockVeloClient.EXPECT().
			RebalanceReserve(context.Background(), gomock.AssignableToTypeOf(vtxnbuild.RebalanceReserve{})).
			Return(vclient.RebalanceReserveResult{}, errors.New("some error has occurred"))

		output, err := helper.logic.RebalanceReserve(&entity.RebalanceInput{
			Passphrase: passPhrase,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "failed to rebalance reserve")
	})
}

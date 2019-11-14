package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	vclient "gitlab.com/velo-labs/cen/libs/client"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"testing"
)

func TestLogic_SetupCredit(t *testing.T) {
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
			SetupCredit(context.Background(), gomock.AssignableToTypeOf(vtxnbuild.SetupCredit{})).
			Return(vclient.SetupCreditResult{
				HorizonResult: &horizon.TransactionSuccess{},
			}, nil)

		output, err := helper.logic.SetupCredit(&entity.SetupCreditInput{
			Passphrase:     "password",
			PeggedCurrency: "THB",
			PeggedValue:    "1",
			AssetCode:      "vTHB",
		})

		assert.NoError(t, err)
		assert.Equal(t, "THB", output.PeggedCurrency)
		assert.Equal(t, "1", output.PeggedValue)
		assert.Equal(t, "vTHB", output.AssetCode)
		assert.Equal(t, stellarAccountEntity().Address, output.SourceAddress)
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

		_, err := helper.logic.SetupCredit(&entity.SetupCreditInput{
			Passphrase:     "password",
			PeggedCurrency: "THB",
			PeggedValue:    "1",
			AssetCode:      "vTHB",
		})

		assert.Error(t, err)
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

		_, err := helper.logic.SetupCredit(&entity.SetupCreditInput{
			Passphrase:     "badPassword",
			PeggedCurrency: "THB",
			PeggedValue:    "1",
			AssetCode:      "vTHB",
		})

		assert.Error(t, err)
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
			SetupCredit(context.Background(), gomock.AssignableToTypeOf(vtxnbuild.SetupCredit{})).
			Return(vclient.SetupCreditResult{}, errors.New("some error has occurred"))

		_, err := helper.logic.SetupCredit(&entity.SetupCreditInput{
			Passphrase:     "password",
			PeggedCurrency: "THB",
			PeggedValue:    "1",
			AssetCode:      "vTHB",
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to setup credit")
	})
}

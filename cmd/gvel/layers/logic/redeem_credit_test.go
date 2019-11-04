package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"testing"
)

func TestLogic_RedeemCredit(t *testing.T) {
	var (
		assetCode   = "kBEAM"
		assetIssuer = "GC3COBQESTRET2AXK5ADR63L7LOMEZWDPODW4F2Z7Y44TTEOTRBSKXQ3"
		amount      = "100"
		passPhrase  = "password"
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
			RedeemCredit(context.Background(), gomock.AssignableToTypeOf(vtxnbuild.RedeemCredit{})).
			Return(&horizon.TransactionSuccess{}, nil)

		output, err := helper.logic.RedeemCredit(&entity.RedeemCreditInput{
			AssetCode:   assetCode,
			AssetIssuer: assetIssuer,
			Amount:      amount,
			Passphrase:  passPhrase,
		})

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, assetCode, output.AssetCode)
		assert.Equal(t, assetIssuer, output.AssetIssuer)
		assert.Equal(t, amount, output.Amount)
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

		output, err := helper.logic.RedeemCredit(&entity.RedeemCreditInput{
			AssetCode:   assetCode,
			AssetIssuer: assetIssuer,
			Amount:      amount,
			Passphrase:  passPhrase,
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

		output, err := helper.logic.RedeemCredit(&entity.RedeemCreditInput{
			AssetCode:   assetCode,
			AssetIssuer: assetIssuer,
			Amount:      amount,
			Passphrase:  "bad passphrase",
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
			RedeemCredit(context.Background(), gomock.AssignableToTypeOf(vtxnbuild.RedeemCredit{})).
			Return(nil, errors.New("some error has occurred"))

		output, err := helper.logic.RedeemCredit(&entity.RedeemCreditInput{
			AssetCode:   assetCode,
			AssetIssuer: assetIssuer,
			Amount:      amount,
			Passphrase:  passPhrase,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "failed to redeem stable credit")
	})
}
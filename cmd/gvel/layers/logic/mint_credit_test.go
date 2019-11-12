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

func TestLogic_MintCredit(t *testing.T) {
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
			MintCredit(context.Background(), gomock.AssignableToTypeOf(vtxnbuild.MintCredit{})).
			Return(vclient.MintCreditResult{
				HorizonResult: &horizon.TransactionSuccess{},
			}, nil)

		output, err := helper.logic.MintCredit(&entity.MintCreditInput{
			AssetToBeMinted:     "kBeam",
			CollateralAssetCode: "THB",
			CollateralAmount:    "100",
			Passphrase:          "password",
		})

		assert.NoError(t, err)
		assert.Equal(t, "THB", output.CollateralAssetCode)
		assert.Equal(t, "100", output.CollateralAmount)
		assert.Equal(t, "kBeam", output.AssetToBeMinted)
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

		output, err := helper.logic.MintCredit(&entity.MintCreditInput{
			AssetToBeMinted:     "kBeam",
			CollateralAssetCode: "THB",
			CollateralAmount:    "100",
			Passphrase:          "strong_password!",
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

		output, err := helper.logic.MintCredit(&entity.MintCreditInput{
			AssetToBeMinted:     "kBeam",
			CollateralAssetCode: "THB",
			CollateralAmount:    "100",
			Passphrase:          "bad passphrase",
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
			MintCredit(context.Background(), gomock.AssignableToTypeOf(vtxnbuild.MintCredit{})).
			Return(vclient.MintCreditResult{}, errors.New("some error has occurred"))

		output, err := helper.logic.MintCredit(&entity.MintCreditInput{
			AssetToBeMinted:     "kBeam",
			CollateralAssetCode: "THB",
			CollateralAmount:    "100",
			Passphrase:          "password",
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "failed to mint credit")
	})
}

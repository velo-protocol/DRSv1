package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stretchr/testify/assert"
	"github.com/velo-protocol/DRSv1/cmd/gvel/entity"
	cenGrpc "github.com/velo-protocol/DRSv1/grpc"
	"github.com/velo-protocol/DRSv1/libs/client"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"testing"
)

func TestLogic_RedeemCredit(t *testing.T) {
	var (
		assetCodeToBeRedeemed   = "kBEAM"
		assetIssuerToBeRedeemed = "GC3COBQESTRET2AXK5ADR63L7LOMEZWDPODW4F2Z7Y44TTEOTRBSKXQ3"
		amountToBeRedeemed      = "100"
		collateralCode          = "VELO"
		collateralIssuer        = "GVI..."
		collateralAmount        = "100"
		passPhrase              = "password"
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
			Return(vclient.RedeemCreditResult{
				HorizonResult: &horizon.TransactionSuccess{},
				VeloNodeResult: &cenGrpc.RedeemCreditOpResponse{
					AssetCodeToBeRedeemed:   assetCodeToBeRedeemed,
					AssetIssuerToBeRedeemed: assetIssuerToBeRedeemed,
					AssetAmountToBeRedeemed: amountToBeRedeemed,
					CollateralCode:          collateralCode,
					CollateralIssuer:        collateralIssuer,
					CollateralAmount:        collateralAmount,
				},
			}, nil)

		output, err := helper.logic.RedeemCredit(&entity.RedeemCreditInput{
			AssetCodeToBeRedeemed:   assetCodeToBeRedeemed,
			AssetIssuerToBeRedeemed: assetIssuerToBeRedeemed,
			AmountToBeRedeemed:      amountToBeRedeemed,
			Passphrase:              passPhrase,
		})

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, assetCodeToBeRedeemed, output.AssetCodeToBeRedeemed)
		assert.Equal(t, assetIssuerToBeRedeemed, output.AssetIssuerToBeRedeemed)
		assert.Equal(t, amountToBeRedeemed, output.AmountToBeRedeemed)
		assert.Equal(t, collateralCode, output.CollateralCode)
		assert.Equal(t, collateralIssuer, output.CollateralIssuer)
		assert.Equal(t, collateralAmount, output.CollateralAmount)
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
			AssetCodeToBeRedeemed:   assetCodeToBeRedeemed,
			AssetIssuerToBeRedeemed: assetIssuerToBeRedeemed,
			AmountToBeRedeemed:      amountToBeRedeemed,
			Passphrase:              passPhrase,
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
			AssetCodeToBeRedeemed:   assetCodeToBeRedeemed,
			AssetIssuerToBeRedeemed: assetIssuerToBeRedeemed,
			AmountToBeRedeemed:      amountToBeRedeemed,
			Passphrase:              "bad passphrase",
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
			Return(vclient.RedeemCreditResult{}, errors.New("some error has occurred"))

		output, err := helper.logic.RedeemCredit(&entity.RedeemCreditInput{
			AssetCodeToBeRedeemed:   assetCodeToBeRedeemed,
			AssetIssuerToBeRedeemed: assetIssuerToBeRedeemed,
			AmountToBeRedeemed:      amountToBeRedeemed,
			Passphrase:              passPhrase,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "failed to redeem stable credit")
	})
}

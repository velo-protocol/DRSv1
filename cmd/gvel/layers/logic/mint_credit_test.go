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

func TestLogic_MintCredit(t *testing.T) {
	var (
		assetCodeToBeMint   = "kBEAM"
		collateralAssetCode = "VELO"
		collateralAmount    = "100"
		passphrase          = "password"

		assetAmountToBeIssued      = "100"
		assetCodeToBeIssued        = "kBEAM"
		assetIssuerToBeIssued      = "GBI..."
		assetDistributorToBeIssued = "GAD..."
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
			MintCredit(context.Background(), gomock.AssignableToTypeOf(vtxnbuild.MintCredit{})).
			Return(vclient.MintCreditResult{
				HorizonResult: &horizon.TransactionSuccess{},
				VeloNodeResult: &cenGrpc.MintCreditOpResponse{
					AssetAmountToBeIssued:      assetAmountToBeIssued,
					AssetCodeToBeIssued:        assetCodeToBeIssued,
					AssetIssuerToBeIssued:      assetIssuerToBeIssued,
					AssetDistributorToBeIssued: assetDistributorToBeIssued,
					CollateralAmount:           collateralAmount,
					CollateralAssetCode:        collateralAssetCode,
				},
			}, nil)

		output, err := helper.logic.MintCredit(&entity.MintCreditInput{
			AssetCodeToBeMinted: assetCodeToBeMint,
			CollateralAssetCode: collateralAssetCode,
			CollateralAmount:    collateralAmount,
			Passphrase:          passphrase,
		})

		assert.NoError(t, err)
		assert.Equal(t, collateralAssetCode, output.CollateralAssetCode)
		assert.Equal(t, collateralAmount, output.CollateralAmount)
		assert.Equal(t, assetCodeToBeMint, output.AssetCodeToBeMinted)
		assert.Equal(t, assetCodeToBeIssued, output.AssetCodeToBeMinted)
		assert.Equal(t, assetIssuerToBeIssued, output.AssetIssuerToBeIssued)
		assert.Equal(t, assetDistributorToBeIssued, output.AssetDistributorToBeIssued)
		assert.Equal(t, assetAmountToBeIssued, output.AssetAmountToBeIssued)
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
			AssetCodeToBeMinted: assetCodeToBeMint,
			CollateralAssetCode: collateralAssetCode,
			CollateralAmount:    collateralAmount,
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
			AssetCodeToBeMinted: assetCodeToBeMint,
			CollateralAssetCode: collateralAssetCode,
			CollateralAmount:    collateralAmount,
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
			AssetCodeToBeMinted: assetCodeToBeMint,
			CollateralAssetCode: collateralAssetCode,
			CollateralAmount:    collateralAmount,
			Passphrase:          passphrase,
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "failed to mint credit")
	})
}

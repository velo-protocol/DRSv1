package credit_test

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stretchr/testify/assert"
	"github.com/velo-protocol/DRSv1/cmd/gvel/entity"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/console"
	"testing"
)

func TestCommandHandler_Mint(t *testing.T) {
	var (
		assetCodeToBeMint   = "kBEAM"
		collateralAssetCode = "VELO"
		collateralAmount    = "100"
		passphrase          = "password"

		assetIssuerToBeIssued      = "GBI..."
		assetDistributorToBeIssued = "GAD..."
	)

	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockPrompt.EXPECT().
			RequestString("Please input asset code of credit to be minted", nil).
			Return(assetCodeToBeMint)
		helper.mockPrompt.EXPECT().
			RequestString("Please input asset code of collateral", nil).
			Return(collateralAmount)
		helper.mockPrompt.EXPECT().
			RequestString("Please input amount of collateral", nil).
			Return(collateralAssetCode)
		helper.mockPrompt.EXPECT().
			RequestHiddenString("🔑 Please input passphrase", nil).
			Return(passphrase)
		helper.mockLogic.EXPECT().
			MintCredit(gomock.AssignableToTypeOf(&entity.MintCreditInput{})).
			Return(&entity.MintCreditOutput{
				AssetCodeToBeMinted:        assetCodeToBeMint,
				CollateralAssetCode:        collateralAssetCode,
				CollateralAmount:           collateralAmount,
				AssetIssuerToBeIssued:      assetIssuerToBeIssued,
				AssetDistributorToBeIssued: assetDistributorToBeIssued,
				SourceAddress:              "GA...",
				TxResult: &horizon.TransactionSuccess{
					Hash: "264226cb06af3b86299031884175155e67a02e0a8ad0b3ab3a88b409a8c09d5c",
				},
			}, nil)

		helper.creditCommandHandler.Mint(helper.mintCmd, nil)

		logEntries := helper.logHook.Entries
		assert.Equal(t, "100 kBEAM minted successfully. The stable credit is in GAD...", logEntries[0].Message)
		assert.Equal(t, fmt.Sprintf("🔗 Stellar Transaction Hash: %s", "264226cb06af3b86299031884175155e67a02e0a8ad0b3ab3a88b409a8c09d5c"), logEntries[1].Message)
	})

	t.Run("error, logic.MintCredit returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockPrompt.EXPECT().
			RequestString("Please input asset code of credit to be minted", nil).
			Return(assetCodeToBeMint)
		helper.mockPrompt.EXPECT().
			RequestString("Please input asset code of collateral", nil).
			Return(collateralAmount)
		helper.mockPrompt.EXPECT().
			RequestString("Please input amount of collateral", nil).
			Return(collateralAssetCode)
		helper.mockPrompt.EXPECT().
			RequestHiddenString("🔑 Please input passphrase", nil).
			Return(passphrase)

		helper.mockLogic.EXPECT().
			MintCredit(gomock.AssignableToTypeOf(&entity.MintCreditInput{})).
			Return(nil, errors.New("some error has occurred"))

		assert.PanicsWithValue(t, console.ExitError, func() {
			helper.creditCommandHandler.Mint(helper.mintCmd, nil)
		})

	})

}

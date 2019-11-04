package credit_test

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
	"testing"
)

func TestCommandHandler_Redeem(t *testing.T) {
	var (
		assetCode   = "kBeam"
		assetIssuer = "GC3COBQESTRET2AXK5ADR63L7LOMEZWDPODW4F2Z7Y44TTEOTRBSKXQ3"
		amount      = "100"
		passPhrase  = "password"
	)

	t.Run("happy", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockPrompt.EXPECT().
			RequestString("Please input asset code of the stable credit to be redeemed", nil).
			Return(assetCode)
		helper.mockPrompt.EXPECT().
			RequestString("Please input issuing account of the stable credit", nil).
			Return(assetIssuer)
		helper.mockPrompt.EXPECT().
			RequestString("Please input the amount of stable credit", nil).
			Return(amount)
		helper.mockPrompt.EXPECT().
			RequestHiddenString("Please input passphrase", nil).
			Return(passPhrase)
		helper.mockLogic.EXPECT().
			RedeemCredit(gomock.AssignableToTypeOf(&entity.RedeemCreditInput{})).
			Return(&entity.RedeemCreditOutput{
				AssetCode:   assetCode,
				AssetIssuer: assetIssuer,
				Amount:      amount,
				TxResult: &horizon.TransactionSuccess{
					Hash: "264226cb06af3b86299031884175155e67a02e0a8ad0b3ab3a88b409a8c09d5c",
				},
			}, nil)

		helper.creditCommandHandler.Redeem(helper.mintCmd, nil)

		logEntries := helper.logHook.Entries
		msgExpected := fmt.Sprintf("Redemption was successfully. You got %s VELO.", amount)
		assert.Equal(t, msgExpected, logEntries[0].Message)
		assert.Equal(t, fmt.Sprintf("🔗 Stellar Transaction Hash: %s", "264226cb06af3b86299031884175155e67a02e0a8ad0b3ab3a88b409a8c09d5c"), logEntries[1].Message)
	})

	t.Run("error, logic.MintCredit returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockPrompt.EXPECT().
			RequestString("Please input asset code of the stable credit to be redeemed", nil).
			Return(assetCode)
		helper.mockPrompt.EXPECT().
			RequestString("Please input issuing account of the stable credit", nil).
			Return(assetIssuer)
		helper.mockPrompt.EXPECT().
			RequestString("Please input the amount of stable credit", nil).
			Return(amount)
		helper.mockPrompt.EXPECT().
			RequestHiddenString("Please input passphrase", nil).
			Return(passPhrase)

		helper.mockLogic.EXPECT().
			RedeemCredit(gomock.AssignableToTypeOf(&entity.RedeemCreditInput{})).
			Return(nil, errors.New("some error has occurred"))

		assert.PanicsWithValue(t, console.ExitError, func() {
			helper.creditCommandHandler.Redeem(helper.redeemCmd, nil)
		})

		logEntries := helper.logHook.Entries

		assert.Equal(t, errors.New("some error has occurred").Error(), logEntries[0].Message)
	})

}

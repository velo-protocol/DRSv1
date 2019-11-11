package collateral_test

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

func TestCommandHandler_RebalanceReserve(t *testing.T) {
	var (
		passPhrase = "password"
	)

	t.Run("happy", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockPrompt.EXPECT().
			RequestHiddenString("🔑 Please input passphrase", nil).
			Return(passPhrase)
		helper.mockLogic.EXPECT().
			RebalanceReserve(gomock.AssignableToTypeOf(&entity.RebalanceInput{})).
			Return(&entity.RebalanceOutput{
				TxResult: &horizon.TransactionSuccess{
					Hash: "264226cb06af3b86299031884175155e67a02e0a8ad0b3ab3a88b409a8c09d5c",
				},
			}, nil)

		helper.collateralCommandHandler.RebalanceReserve(helper.rebalanceCmd, nil)

		logEntries := helper.logHook.Entries
		assert.Equal(t, "Rebalancing completed.", logEntries[0].Message)
		assert.Equal(t, fmt.Sprintf("🔗 Stellar Transaction Hash: %s", "264226cb06af3b86299031884175155e67a02e0a8ad0b3ab3a88b409a8c09d5c"), logEntries[1].Message)
	})

	t.Run("error, logic.RebalanceReserve returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockPrompt.EXPECT().
			RequestHiddenString("🔑 Please input passphrase", nil).
			Return(passPhrase)

		helper.mockLogic.EXPECT().
			RebalanceReserve(gomock.AssignableToTypeOf(&entity.RebalanceInput{})).
			Return(nil, errors.New("some error has occurred"))

		assert.PanicsWithValue(t, console.ExitError, func() {
			helper.collateralCommandHandler.RebalanceReserve(helper.rebalanceCmd, nil)
		})

		logEntries := helper.logHook.Entries

		assert.Equal(t, errors.New("some error has occurred").Error(), logEntries[0].Message)
	})

}

package commands

import (
	"github.com/stretchr/testify/assert"
	"github.com/velo-protocol/DRSv1/cmd/gvel/constants"
	"testing"
)

func TestGvelHandler_Init(t *testing.T) {
	gvelHandler := NewGvelHandler(nil, nil)
	gvelHandler.Init()

	assert.Equal(t, constants.CmdRootGvel, gvelHandler.RootCommand.Use)

	assert.True(t, gvelHandler.AccountCommand == gvelHandler.RootCommand.Commands()[0])
	assert.True(t, gvelHandler.CollateralCommand == gvelHandler.RootCommand.Commands()[1])
	assert.True(t, gvelHandler.CreditCommand == gvelHandler.RootCommand.Commands()[2])
	assert.True(t, gvelHandler.InitCommand == gvelHandler.RootCommand.Commands()[3])

	assert.Len(t, gvelHandler.RootCommand.Commands(), 4)
	assert.Contains(t, gvelHandler.AccountCommand.Use, constants.CmdAccount)
	assert.Contains(t, gvelHandler.CollateralCommand.Use, constants.CmdCollateral)
	assert.Contains(t, gvelHandler.CreditCommand.Use, constants.CmdCredit)
	assert.Equal(t, constants.CmdInit, gvelHandler.InitCommand.Use)

	assert.Len(t, gvelHandler.AccountCommand.Commands(), 5)
	assert.Equal(t, constants.CmdAccountCreate, gvelHandler.AccountCommand.Commands()[0].Use)
	assert.Equal(t, constants.CmdAccountDefault, gvelHandler.AccountCommand.Commands()[1].Use)
	assert.Equal(t, constants.CmdAccountExport, gvelHandler.AccountCommand.Commands()[2].Use)
	assert.Equal(t, constants.CmdAccountImport, gvelHandler.AccountCommand.Commands()[3].Use)
	assert.Equal(t, constants.CmdAccountList, gvelHandler.AccountCommand.Commands()[4].Use)

	assert.Len(t, gvelHandler.CollateralCommand.Commands(), 2)
	assert.Equal(t, constants.CmdCollateralHealthCheck, gvelHandler.CollateralCommand.Commands()[0].Use)
	assert.Equal(t, constants.CmdCollateralRebalance, gvelHandler.CollateralCommand.Commands()[1].Use)

	assert.Len(t, gvelHandler.CreditCommand.Commands(), 4)
	assert.Equal(t, constants.CmdCreditGetExchange, gvelHandler.CreditCommand.Commands()[0].Use)
	assert.Equal(t, constants.CmdCreditMint, gvelHandler.CreditCommand.Commands()[1].Use)
	assert.Equal(t, constants.CmdCreditRedeem, gvelHandler.CreditCommand.Commands()[2].Use)
	assert.Equal(t, constants.CmdCreditSetup, gvelHandler.CreditCommand.Commands()[3].Use)

	assert.Len(t, gvelHandler.InitCommand.Commands(), 0)
}

package vxdr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewOperationBody(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		_, err := NewOperationBody(OperationTypeWhitelist, WhitelistOp{})
		assert.NoError(t, err)
		_, err = NewOperationBody(OperationTypeSetupCredit, SetupCreditOp{})
		assert.NoError(t, err)
		_, err = NewOperationBody(OperationTypePriceUpdate, PriceUpdateOp{})
		assert.NoError(t, err)
		_, err = NewOperationBody(OperationTypeMintCredit, MintCreditOp{})
		assert.NoError(t, err)
		_, err = NewOperationBody(OperationTypeRedeemCredit, RedeemCreditOp{})
		assert.NoError(t, err)
		_, err = NewOperationBody(OperationTypeRebalanceReserve, RebalanceReserveOp{})
		assert.NoError(t, err)
	})

	t.Run("error, bad value", func(t *testing.T) {
		_, err := NewOperationBody(OperationTypeWhitelist, nil)
		assert.Error(t, err)
		_, err = NewOperationBody(OperationTypeSetupCredit, nil)
		assert.Error(t, err)
		_, err = NewOperationBody(OperationTypePriceUpdate, nil)
		assert.Error(t, err)
		_, err = NewOperationBody(OperationTypeMintCredit, nil)
		assert.Error(t, err)
		_, err = NewOperationBody(OperationTypeRedeemCredit, nil)
		assert.Error(t, err)
		_, err = NewOperationBody(OperationTypeRebalanceReserve, nil)
		assert.Error(t, err)
	})

	t.Run("error, unknown operation type", func(t *testing.T) {
		_, err := NewOperationBody(99, nil)
		assert.Error(t, err)
	})
}

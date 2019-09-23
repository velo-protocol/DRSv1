package vtxnbuild

import (
	"github.com/stellar/go/xdr"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"testing"
)

func TestSetUpCredit_BuildXDR(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		veloXdrOp, err := (&SetupCredit{
			PeggedCurrency: "THB",
			PeggedValue:    "1.00",
			AssetCode:      "vTHB",
		}).BuildXDR()

		assert.NoError(t, err)
		assert.Equal(t, vxdr.OperationTypeSetupCredit, veloXdrOp.Body.Type)
		assert.Equal(t, "THB", veloXdrOp.Body.SetupCreditOp.PeggedCurrency)
		assert.Equal(t, xdr.Int64(10000000), veloXdrOp.Body.SetupCreditOp.PeggedValue)
		assert.Equal(t, "vTHB", veloXdrOp.Body.SetupCreditOp.AssetCode)
	})
	t.Run("error, failed to parse pegged value", func(t *testing.T) {
		_, err := (&SetupCredit{
			PeggedCurrency: "THB",
			PeggedValue:    "BAD_VALUE",
			AssetCode:      "vTHB",
		}).BuildXDR()

		assert.Error(t, err)
	})
}

func TestSetUpCredit_FromXDR(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		veloXdrOp := vxdr.VeloOp{
			Body: vxdr.OperationBody{
				SetupCreditOp: &vxdr.SetupCreditOp{
					PeggedCurrency: "THB",
					PeggedValue:    xdr.Int64(10000000),
					AssetCode:      "vTHB",
				},
			},
		}

		var newVeloSetupCreditOp SetupCredit
		err := newVeloSetupCreditOp.FromXDR(veloXdrOp)

		assert.NoError(t, err)
		assert.Equal(t, "THB", newVeloSetupCreditOp.PeggedCurrency)
		assert.Equal(t, "1.0000000", newVeloSetupCreditOp.PeggedValue)
		assert.Equal(t, "vTHB", newVeloSetupCreditOp.AssetCode)
	})
	t.Run("error, empty SetupCreditOp", func(t *testing.T) {
		veloXdrOp := vxdr.VeloOp{
			Body: vxdr.OperationBody{
				SetupCreditOp: nil,
			},
		}
		var newVeloSetupCreditOp SetupCredit
		err := newVeloSetupCreditOp.FromXDR(veloXdrOp)
		assert.Error(t, err)
	})
}

func TestSetUpCredit_Validate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		err := (&SetupCredit{
			PeggedCurrency: "THB1",
			PeggedValue:    "1.00",
			AssetCode:      "vTHB",
		}).Validate()

		assert.NoError(t, err)
	})
	t.Run("error, empty pegged currency", func(t *testing.T) {
		err := (&SetupCredit{
			PeggedCurrency: "",
			PeggedValue:    "1.00",
			AssetCode:      "vTHB1",
		}).Validate()

		assert.Error(t, err)
	})
	t.Run("error, pegged currency too long", func(t *testing.T) {
		err := (&SetupCredit{
			PeggedCurrency: "1234567890XXX",
			PeggedValue:    "1.00",
			AssetCode:      "vTHB",
		}).Validate()

		assert.Error(t, err)
	})
	t.Run("error, pegged currency cannot have special character", func(t *testing.T) {
		err := (&SetupCredit{
			PeggedCurrency: "_THB",
			PeggedValue:    "1.00",
			AssetCode:      "vTHB",
		}).Validate()

		assert.Error(t, err)
	})
	t.Run("error, empty asset name", func(t *testing.T) {
		err := (&SetupCredit{
			PeggedCurrency: "THB1",
			PeggedValue:    "1.00",
			AssetCode:      "",
		}).Validate()

		assert.Error(t, err)
	})
	t.Run("error, asset name too long", func(t *testing.T) {
		err := (&SetupCredit{
			PeggedCurrency: "THB1",
			PeggedValue:    "1.00",
			AssetCode:      "1234567890XXX",
		}).Validate()

		assert.Error(t, err)
	})
	t.Run("error, asset name cannot have special character", func(t *testing.T) {
		err := (&SetupCredit{
			PeggedCurrency: "THB1",
			PeggedValue:    "1.00",
			AssetCode:      "_vTHB",
		}).Validate()

		assert.Error(t, err)
	})
	t.Run("error, pegged value must be number", func(t *testing.T) {
		err := (&SetupCredit{
			PeggedCurrency: "THB1",
			PeggedValue:    "1.00XXX",
			AssetCode:      "vTHB",
		}).Validate()

		assert.Error(t, err)
	})
	t.Run("error, pegged value must be greater than zero", func(t *testing.T) {
		err := (&SetupCredit{
			PeggedCurrency: "THB1",
			PeggedValue:    "-1.00",
			AssetCode:      "vTHB",
		}).Validate()

		assert.Error(t, err)
	})
}

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
			AssetCode:      "VELO",
		}).BuildXDR()

		assert.NoError(t, err)
		assert.Equal(t, vxdr.OperationTypeSetupCredit, veloXdrOp.Body.Type)
		assert.Equal(t, "THB", veloXdrOp.Body.SetupCreditOp.PeggedCurrency)
		assert.Equal(t, xdr.Int64(10000000), veloXdrOp.Body.SetupCreditOp.PeggedValue)
		assert.Equal(t, "VELO", veloXdrOp.Body.SetupCreditOp.AssetCode)
	})
	t.Run("error, validation fail", func(t *testing.T) {
		_, err := (&SetupCredit{
			PeggedCurrency: "THB",
			PeggedValue:    "1.00",
			AssetCode:      "",
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
					AssetCode:      "VELO",
				},
			},
		}

		var newVeloSetupCreditOp SetupCredit
		err := newVeloSetupCreditOp.FromXDR(veloXdrOp)

		assert.NoError(t, err)
		assert.Equal(t, "THB", newVeloSetupCreditOp.PeggedCurrency)
		assert.Equal(t, "1.0000000", newVeloSetupCreditOp.PeggedValue)
		assert.Equal(t, "VELO", newVeloSetupCreditOp.AssetCode)
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
			PeggedCurrency: "THB",
			PeggedValue:    "1.00",
			AssetCode:      "VELO",
		}).Validate()

		assert.NoError(t, err)
	})
	t.Run("error, empty asset code", func(t *testing.T) {
		err := (&SetupCredit{
			PeggedCurrency: "THB",
			PeggedValue:    "1.00",
			AssetCode:      "",
		}).Validate()

		assert.Error(t, err, "assetCode parameter cannot be blank")
	})
	t.Run("error, empty pegged value", func(t *testing.T) {
		err := (&SetupCredit{
			PeggedCurrency: "THB",
			PeggedValue:    "",
			AssetCode:      "VELO",
		}).Validate()

		assert.Error(t, err, "peggedValue parameter cannot be blank")
	})
	t.Run("error, empty pegged currency", func(t *testing.T) {
		err := (&SetupCredit{
			PeggedCurrency: "",
			PeggedValue:    "1.00",
			AssetCode:      "VELO",
		}).Validate()

		assert.Error(t, err, "peggedCurrency parameter cannot be blank")
	})
	t.Run("error, pegged value must be a number", func(t *testing.T) {
		err := (&SetupCredit{
			PeggedCurrency: "THB",
			PeggedValue:    "1.0000XXX",
			AssetCode:      "VELO",
		}).Validate()

		assert.Error(t, err, "peggedValue parameter is not a number")
	})
	t.Run("error, pegged value must greater than zero", func(t *testing.T) {
		err := (&SetupCredit{
			PeggedCurrency: "THB",
			PeggedValue:    "-1.000",
			AssetCode:      "VELO",
		}).Validate()

		assert.Error(t, err, "peggedValue must be greater than zero")
	})
	t.Run("error, invalid format of asset code, asset code too long", func(t *testing.T) {
		err := (&SetupCredit{
			PeggedCurrency: "THB",
			PeggedValue:    "1.00",
			AssetCode:      "12345678",
		}).Validate()

		assert.Error(t, err, "invalid format of asset name")
	})
	t.Run("error, invalid format of asset code, asset code cannot have special character", func(t *testing.T) {
		err := (&SetupCredit{
			PeggedCurrency: "THB",
			PeggedValue:    "1.00",
			AssetCode:      "_VELO",
		}).Validate()

		assert.Error(t, err, "invalid format of asset name")
	})
	t.Run("error, invalid format of pegged currency, pegged currency too long", func(t *testing.T) {
		err := (&SetupCredit{
			PeggedCurrency: "1234567890XXX",
			PeggedValue:    "1.00",
			AssetCode:      "vTHB",
		}).Validate()

		assert.Error(t, err, "invalid format of pegged currency")
	})
	t.Run("error, invalid format of pegged currency, pegged currency cannot have special character", func(t *testing.T) {
		err := (&SetupCredit{
			PeggedCurrency: "_THB",
			PeggedValue:    "1.00",
			AssetCode:      "VELO",
		}).Validate()

		assert.Error(t, err, "invalid format of pegged currency")
	})
	t.Run("error, asset is not VELO", func(t *testing.T) {
		err := (&SetupCredit{
			PeggedCurrency: "THB",
			PeggedValue:    "1.00",
			AssetCode:      "vTHB",
		}).Validate()

		assert.Error(t, err)
		assert.Errorf(t, err, "assetCode %s does not exist", "vTHB")
	})
}

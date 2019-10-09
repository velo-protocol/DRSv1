package vtxnbuild

import (
	"github.com/stellar/go/xdr"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"testing"
)

func TestPriceUpdate_BuildXDR(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		veloXdrOp, err := (&PriceUpdate{
			Asset:                       "VELO",
			Currency:                    "THB",
			PriceInCurrencyPerAssetUnit: "0.5",
		}).BuildXDR()

		assert.NoError(t, err)
		assert.Equal(t, vxdr.OperationTypePriceUpdate, veloXdrOp.Body.Type)
		assert.Equal(t, "VELO", veloXdrOp.Body.PriceUpdateOp.Asset)
		assert.Equal(t, vxdr.CurrencyTHB, veloXdrOp.Body.PriceUpdateOp.Currency)
		assert.Equal(t, xdr.Int64(5000000), veloXdrOp.Body.PriceUpdateOp.PriceInCurrencyPerAssetUnit)
	})
	t.Run("error, validation fail", func(t *testing.T) {
		_, err := (&PriceUpdate{
			Asset:                       "",
			Currency:                    "THB",
			PriceInCurrencyPerAssetUnit: "0.5",
		}).BuildXDR()

		assert.Error(t, err)
	})
}

func TestPriceUpdate_FromXDR(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		veloXdrOp := vxdr.VeloOp{
			Body: vxdr.OperationBody{
				PriceUpdateOp: &vxdr.PriceUpdateOp{
					Asset:                       "VELO",
					Currency:                    vxdr.CurrencyTHB,
					PriceInCurrencyPerAssetUnit: xdr.Int64(5000000),
				},
			},
		}

		var newVeloPriceUpdateOp PriceUpdate
		err := newVeloPriceUpdateOp.FromXDR(veloXdrOp)

		assert.NoError(t, err)
		assert.Equal(t, "VELO", newVeloPriceUpdateOp.Asset)
		assert.Equal(t, "THB", newVeloPriceUpdateOp.Currency)
		assert.Equal(t, "0.5000000", newVeloPriceUpdateOp.PriceInCurrencyPerAssetUnit)
	})
	t.Run("error, empty PriceUpdateOp", func(t *testing.T) {
		veloXdrOp := vxdr.VeloOp{
			Body: vxdr.OperationBody{
				PriceUpdateOp: nil,
			},
		}

		var newVeloPriceUpdateOp PriceUpdate
		err := newVeloPriceUpdateOp.FromXDR(veloXdrOp)
		assert.Error(t, err)
	})
}

func TestPriceUpdate_Validate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		err := (&PriceUpdate{
			Asset:                       "VELO",
			Currency:                    "THB",
			PriceInCurrencyPerAssetUnit: "0.5",
		}).Validate()

		assert.NoError(t, err)
	})
	t.Run("error, empty asset", func(t *testing.T) {
		err := (&PriceUpdate{
			Asset:                       "",
			Currency:                    "THB",
			PriceInCurrencyPerAssetUnit: "0.5",
		}).Validate()

		assert.EqualError(t, err, "asset must not be blank")
	})
	t.Run("error, empty currency", func(t *testing.T) {
		err := (&PriceUpdate{
			Asset:                       "VELO",
			Currency:                    "",
			PriceInCurrencyPerAssetUnit: "0.5",
		}).Validate()

		assert.EqualError(t, err, "currency must not be blank")
	})
	t.Run("error, empty priceInCurrencyPerAssetUnit", func(t *testing.T) {
		err := (&PriceUpdate{
			Asset:                       "VELO",
			Currency:                    "THB",
			PriceInCurrencyPerAssetUnit: "",
		}).Validate()

		assert.EqualError(t, err, "priceInCurrencyPerAssetUnit must not be blank")
	})
	t.Run("error, priceInCurrencyPerAssetUnit is not a number", func(t *testing.T) {
		err := (&PriceUpdate{
			Asset:                       "VELO",
			Currency:                    "THB",
			PriceInCurrencyPerAssetUnit: "1.XXX",
		}).Validate()

		assert.EqualError(t, err, "invalid priceInCurrencyPerAssetUnit format")
	})
	t.Run("error, priceInCurrencyPerAssetUnit is negative", func(t *testing.T) {
		err := (&PriceUpdate{
			Asset:                       "VELO",
			Currency:                    "THB",
			PriceInCurrencyPerAssetUnit: "-0.5",
		}).Validate()

		assert.EqualError(t, err, "priceInCurrencyPerAssetUnit must be greater than zero")
	})
	t.Run("error, currency does not exist", func(t *testing.T) {
		err := (&PriceUpdate{
			Asset:                       "VELO",
			Currency:                    "JPY",
			PriceInCurrencyPerAssetUnit: "0.5",
		}).Validate()

		assert.EqualError(t, err, "the currency JPY does not exist")
	})
	t.Run("error, invalid asset", func(t *testing.T) {
		err := (&PriceUpdate{
			Asset:                       "vTHB",
			Currency:                    "THB",
			PriceInCurrencyPerAssetUnit: "0.5",
		}).Validate()

		assert.EqualError(t, err, "asset vTHB does not exist")
	})
}

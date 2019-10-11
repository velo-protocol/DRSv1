package vtxnbuild

import (
	"github.com/stellar/go/xdr"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"testing"
)

func TestMintCredit_BuildXDR(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		veloXdrOp, err := (&MintCredit{
			AssetCodeToBeIssued: "vUSD",
			CollateralAssetCode: "VELO",
			CollateralAmount:    "1000.00",
		}).BuildXDR()

		assert.NoError(t, err)
		assert.Equal(t, vxdr.OperationTypeMintCredit, veloXdrOp.Body.Type)
		assert.Equal(t, vxdr.AssetVELO, veloXdrOp.Body.MintCreditOp.CollateralAssetCode)
		assert.Equal(t, "vUSD", veloXdrOp.Body.MintCreditOp.AssetCodeToBeIssued)
		assert.Equal(t, xdr.Int64(10000000000), veloXdrOp.Body.MintCreditOp.CollateralAmount)
	})
	t.Run("error, validation fail", func(t *testing.T) {
		_, err := (&MintCredit{
			AssetCodeToBeIssued: "vUSD",
			CollateralAssetCode: "BITCOIN",
			CollateralAmount:    "1000.00",
		}).BuildXDR()

		assert.Error(t, err)
	})
}

func TestMintCredit_FromXDR(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		veloXdrOp := vxdr.VeloOp{
			Body: vxdr.OperationBody{
				MintCreditOp: &vxdr.MintCreditOp{
					AssetCodeToBeIssued: "vUSD",
					CollateralAssetCode: vxdr.AssetVELO,
					CollateralAmount:    xdr.Int64(10000000000),
				},
			},
		}

		var newVeloMintCreditOp MintCredit
		err := newVeloMintCreditOp.FromXDR(veloXdrOp)

		assert.NoError(t, err)
		assert.Equal(t, "vUSD", newVeloMintCreditOp.AssetCodeToBeIssued)
		assert.Equal(t, "VELO", newVeloMintCreditOp.CollateralAssetCode)
		assert.Equal(t, "1000.0000000", newVeloMintCreditOp.CollateralAmount)
	})
	t.Run("error, empty MintCreditOp", func(t *testing.T) {
		veloXdrOp := vxdr.VeloOp{
			Body: vxdr.OperationBody{
				MintCreditOp: nil,
			},
		}

		var newVeloMintCreditOp MintCredit
		err := newVeloMintCreditOp.FromXDR(veloXdrOp)
		assert.Error(t, err)
	})
}

func TestMintCredit_Validate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		err := (&MintCredit{
			AssetCodeToBeIssued: "vUSD",
			CollateralAssetCode: "VELO",
			CollateralAmount:    "1000.00",
		}).Validate()

		assert.NoError(t, err)
	})
	t.Run("error, empty assetCodeToBeIssued", func(t *testing.T) {
		err := (&MintCredit{
			AssetCodeToBeIssued: "",
			CollateralAssetCode: "VELO",
			CollateralAmount:    "1000.00",
		}).Validate()

		assert.EqualError(t, err, "assetCodeToBeIssued must not be blank")
	})
	t.Run("error, empty collateralAssetCode", func(t *testing.T) {
		err := (&MintCredit{
			AssetCodeToBeIssued: "vUSD",
			CollateralAssetCode: "",
			CollateralAmount:    "1000.00",
		}).Validate()

		assert.EqualError(t, err, "collateralAssetCode must not be blank")
	})
	t.Run("error, empty collateralAmount", func(t *testing.T) {
		err := (&MintCredit{
			AssetCodeToBeIssued: "vUSD",
			CollateralAssetCode: "VELO",
			CollateralAmount:    "",
		}).Validate()

		assert.EqualError(t, err, "collateralAmount must not be blank")
	})
	t.Run("error, invalid assetCodeToBeIssued format", func(t *testing.T) {
		err := (&MintCredit{
			AssetCodeToBeIssued: "tooLongAsset",
			CollateralAssetCode: "VELO",
			CollateralAmount:    "1000.00",
		}).Validate()

		assert.EqualError(t, err, "invalid assetCodeToBeIssued format")
	})
	t.Run("error, collateralAssetCode value must be VELO", func(t *testing.T) {
		err := (&MintCredit{
			AssetCodeToBeIssued: "vUSD",
			CollateralAssetCode: "BITCOIN",
			CollateralAmount:    "1000.00",
		}).Validate()

		assert.EqualError(t, err, "collateralAssetCode value must be VELO")
	})
	t.Run("error, invalid collateralAmount format", func(t *testing.T) {
		err := (&MintCredit{
			AssetCodeToBeIssued: "vUSD",
			CollateralAssetCode: "VELO",
			CollateralAmount:    "1000.xx",
		}).Validate()

		assert.EqualError(t, err, "invalid collateralAmount format")
	})
	t.Run("error, collateralAmount must be greater than zero", func(t *testing.T) {
		err := (&MintCredit{
			AssetCodeToBeIssued: "vUSD",
			CollateralAssetCode: "VELO",
			CollateralAmount:    "-1000.00",
		}).Validate()

		assert.EqualError(t, err, "collateralAmount must be greater than zero")
	})
}

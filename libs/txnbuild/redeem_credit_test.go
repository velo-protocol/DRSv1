package vtxnbuild

import (
	"github.com/stellar/go/xdr"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"testing"
)

func TestRedeemCredit_BuildXDR(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		veloXdrOp, err := (&RedeemCredit{
			AssetCode: "vTHB",
			Issuer:    publicKey1,
			Amount:    "0.5",
		}).BuildXDR()

		assert.NoError(t, err)
		assert.Equal(t, vxdr.OperationTypeRedeemCredit, veloXdrOp.Body.Type)
		assert.Equal(t, "vTHB", veloXdrOp.Body.RedeemCreditOp.AssetCode)
		assert.Equal(t, publicKey1, veloXdrOp.Body.RedeemCreditOp.Issuer.Address())
		assert.Equal(t, xdr.Int64(5000000), veloXdrOp.Body.RedeemCreditOp.Amount)
	})
	t.Run("error, validation fail", func(t *testing.T) {
		_, err := (&RedeemCredit{
			AssetCode: "",
			Issuer:    publicKey1,
			Amount:    "0.5",
		}).BuildXDR()

		assert.Error(t, err)
	})
}

func TestRedeemCredit_FromXDR(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var account xdr.AccountId
		_ = account.SetAddress(publicKey1)
		veloXdrOp := vxdr.VeloOp{
			Body: vxdr.OperationBody{
				RedeemCreditOp: &vxdr.RedeemCreditOp{
					AssetCode: "vTHB",
					Issuer:    account,
					Amount:    xdr.Int64(5000000),
				},
			},
		}

		var newVeloRedeemCreditOp RedeemCredit
		err := newVeloRedeemCreditOp.FromXDR(veloXdrOp)

		assert.NoError(t, err)
		assert.Equal(t, "vTHB", newVeloRedeemCreditOp.AssetCode)
		assert.Equal(t, publicKey1, newVeloRedeemCreditOp.Issuer)
		assert.Equal(t, "0.5000000", newVeloRedeemCreditOp.Amount)
	})
	t.Run("error, empty RedeemCreditOp", func(t *testing.T) {
		veloXdrOp := vxdr.VeloOp{
			Body: vxdr.OperationBody{
				RedeemCreditOp: nil,
			},
		}

		var newVeloRedeemCreditOp RedeemCredit
		err := newVeloRedeemCreditOp.FromXDR(veloXdrOp)
		assert.Error(t, err)
	})
}

func TestRedeemCredit_Validate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		err := (&RedeemCredit{
			AssetCode: "vTHB",
			Issuer:    publicKey1,
			Amount:    "0.5",
		}).Validate()

		assert.NoError(t, err)
	})
	t.Run("error, empty assetCode", func(t *testing.T) {
		err := (&RedeemCredit{
			AssetCode: "",
			Issuer:    publicKey1,
			Amount:    "0.5",
		}).Validate()

		assert.Error(t, err)
		assert.EqualError(t, err, "assetCode must not be blank")
	})
	t.Run("error, empty issuer", func(t *testing.T) {
		err := (&RedeemCredit{
			AssetCode: "vTHB",
			Issuer:    "",
			Amount:    "0.5",
		}).Validate()

		assert.Error(t, err)
		assert.EqualError(t, err, "issuer must not be blank")
	})
	t.Run("error, empty amount", func(t *testing.T) {
		err := (&RedeemCredit{
			AssetCode: "vTHB",
			Issuer:    publicKey1,
			Amount:    "",
		}).Validate()

		assert.Error(t, err)
		assert.EqualError(t, err, "amount must not be blank")
	})
	t.Run("error, invalid assetCode format", func(t *testing.T) {
		err := (&RedeemCredit{
			AssetCode: "tooLongAsset",
			Issuer:    publicKey1,
			Amount:    "1000.00",
		}).Validate()

		assert.EqualError(t, err, "invalid assetCode format")
	})
	t.Run("error, invalid issuer format", func(t *testing.T) {
		err := (&RedeemCredit{
			AssetCode: "vTHB",
			Issuer:    "WRONG_ISSUER_ADDRESS",
			Amount:    "1000.00",
		}).Validate()

		assert.EqualError(t, err, "invalid issuer format")
	})
	t.Run("error, invalid amount format", func(t *testing.T) {
		err := (&RedeemCredit{
			AssetCode: "vTHB",
			Issuer:    publicKey1,
			Amount:    "amount",
		}).Validate()

		assert.EqualError(t, err, "invalid amount format")
	})
	t.Run("error, amount less than zero", func(t *testing.T) {
		err := (&RedeemCredit{
			AssetCode: "vTHB",
			Issuer:    publicKey1,
			Amount:    "-100",
		}).Validate()

		assert.EqualError(t, err, "amount must be greater than zero")
	})
}

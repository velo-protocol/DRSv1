package logic_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/velo-protocol/DRSv1/cmd/gvel/entity"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/crypto"
	"testing"
)

func TestLogic_ExportAccount(t *testing.T) {

	var (
		publicKey  = "GBVI3QZYXCWQBSGZ4TNJOHDZ5KZYGZOVSE46TVAYJYTMNCGW2PWLWO73"
		publicKey2 = "GBMGIOVP376BHT4I4VF3C5AUO2J4DGVXKFXOCNDW77OLMV2KGQV3N7KN"
		seedKey    = "SBR25NMQRKQ4RLGNV5XB3MMQB4ADVYSMPGVBODQVJE7KPTDR6KGK3XMX"
	)

	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockDB.EXPECT().
			Get([]byte(publicKey)).
			Return(stellarAccountsBytes(), nil)

		output, err := helper.logic.ExportAccount(&entity.ExportAccountInput{
			PublicKey:  publicKey,
			Passphrase: "password",
		})

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.NotEmpty(t, output.ExportedKeyPair)
		assert.Equal(t, seedKey, output.ExportedKeyPair.Seed())
	})

	t.Run("error, database returns error", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockDB.EXPECT().
			Get([]byte(publicKey)).
			Return(nil, errors.New("some error has occurred"))

		output, err := helper.logic.ExportAccount(&entity.ExportAccountInput{
			PublicKey:  publicKey,
			Passphrase: "password",
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "failed to get account from db")

	})

	t.Run("error, account does not exist", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockDB.EXPECT().
			Get([]byte(publicKey2)).
			Return(nil, nil)

		output, err := helper.logic.ExportAccount(&entity.ExportAccountInput{
			PublicKey:  publicKey2,
			Passphrase: "password",
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Equal(t, fmt.Sprintf("account %s does not exist", publicKey2), err.Error())
	})

	t.Run("error, failed to convert seed key to key pair", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		helper.mockDB.EXPECT().
			Get([]byte(publicKey)).
			Return(stellarAccountsBytes(), nil)

		output, err := helper.logic.ExportAccount(&entity.ExportAccountInput{
			PublicKey:  publicKey,
			Passphrase: "bad password",
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), fmt.Sprintf("failed to decrypt the seed of %s with given passphrase", publicKey))
	})

	t.Run("error, failed to derive keypair from seed", func(t *testing.T) {
		helper := initTest(t)
		defer helper.done()

		encryptedSeed, nonce, _ := crypto.Encrypt([]byte("SEEDKEY"), "password")
		badAccountByte := entity.StellarAccount{
			Address:       publicKey,
			EncryptedSeed: encryptedSeed,
			Nonce:         nonce,
		}
		mockedStellarAccountBytes, _ := json.Marshal(badAccountByte)

		helper.mockDB.EXPECT().
			Get([]byte(publicKey)).
			Return(mockedStellarAccountBytes, nil)

		output, err := helper.logic.ExportAccount(&entity.ExportAccountInput{
			PublicKey:  publicKey,
			Passphrase: "password",
		})

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "failed to derive keypair from seed")
	})
}

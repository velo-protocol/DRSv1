package logic

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/crypto"
	"gitlab.com/velo-labs/cen/libs/convert"
	"strings"
)

func (lo *logic) ImportAccount(input *entity.ImportAccountInput) (*entity.ImportAccountOutput, error) {

	kp, err := vconvert.SecretKeyToKeyPair(input.SeedKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert seed key to key pair")
	}

	account, err := lo.DB.Get([]byte(kp.Address()))
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return nil, errors.Wrapf(err, "failed to get account from db")
	}
	if account != nil {
		return nil, errors.Errorf("account %s is already exist", kp.Address())
	}

	_, err = lo.Stellar.GetStellarAccount(kp.Address())
	if err != nil {
		return nil, errors.Wrap(err, "failed to verify account with stellar")
	}

	encryptedSeed, nonce, err := crypto.Encrypt([]byte(kp.Seed()), input.Passphrase)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt seed key")
	}

	stellarAccount := entity.StellarAccount{
		Address:       kp.Address(),
		EncryptedSeed: encryptedSeed,
		Nonce:         nonce,
	}

	stellarAccountBytes, err := json.Marshal(stellarAccount)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal entity")
	}

	dbKey := kp.Address()

	err = lo.DB.Save([]byte(dbKey), stellarAccountBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save stellar account")
	}

	// set default account
	mustSetDefault := lo.AppConfig.GetDefaultAccount() == "" || input.SetAsDefault
	if mustSetDefault {
		err = lo.AppConfig.SetDefaultAccount(kp.Address())
		if err != nil {
			return nil, errors.Wrap(err, "failed to write config file")
		}
	}

	return &entity.ImportAccountOutput{
		ImportedKeyPair: kp,
		IsDefault:       mustSetDefault,
	}, nil
}

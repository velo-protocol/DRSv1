package logic

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/keypair"
	"github.com/velo-protocol/DRSv1/cmd/gvel/entity"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/console"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/crypto"
)

func (lo *logic) CreateAccount(input *entity.CreateAccountInput) (*entity.CreateAccountOutput, error) {
	newKP, err := keypair.Random()
	if err != nil {
		return nil, errors.Wrap(err, "failed to random a new key pair")
	}

	console.StartLoading("Creating account with %s with starting balance 10000 XLM.", newKP.Address())
	defer console.StopLoading()

	err = lo.Stellar.GetFreeLumens(newKP.Address())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a stellar account")
	}

	dbKey := fmt.Sprintf("%s", newKP.Address())

	encryptedSeed, nonce, err := crypto.Encrypt([]byte(newKP.Seed()), input.Passphrase)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt seed key")
	}

	stellarAccount := entity.StellarAccount{
		Address:       newKP.Address(),
		EncryptedSeed: encryptedSeed,
		Nonce:         nonce,
	}

	stellarAccountBytes, err := json.Marshal(stellarAccount)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal entity")
	}

	err = lo.DB.Save([]byte(dbKey), stellarAccountBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save stellar account")
	}

	// set default account
	mustSetDefault := lo.AppConfig.GetDefaultAccount() == "" || input.SetAsDefaultAccount
	if mustSetDefault {
		err = lo.AppConfig.SetDefaultAccount(newKP.Address())
		if err != nil {
			return nil, errors.Wrap(err, "failed to write config file")
		}
	}

	return &entity.CreateAccountOutput{
		GeneratedKeyPair: newKP,
		IsDefault:        mustSetDefault,
	}, nil
}

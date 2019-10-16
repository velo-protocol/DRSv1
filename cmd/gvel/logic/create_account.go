package logic

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/keypair"
	"gitlab.com/velo-labs/cen/cmd/gvel/crypto"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
)

func (lo *logic) CreateAccount(passphrase string) (*keypair.Full, error) {
	newKP, err := keypair.Random()
	if err != nil {
		return nil, errors.Wrap(err, "failed to random a new key pair")
	}

	err = lo.Friendbot.GetFreeLumens(newKP.Address())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a stellar account")
	}

	dbkey := fmt.Sprintf("%s", newKP.Address())

	encryptedSeed, nonce, err := crypto.Encrypt([]byte(newKP.Seed()), passphrase)
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

	err = lo.DB.Save([]byte(dbkey), stellarAccountBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save stellar account")
	}

	return newKP, nil
}

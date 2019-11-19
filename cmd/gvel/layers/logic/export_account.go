package logic

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/velo-protocol/DRSv1/cmd/gvel/entity"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/crypto"
	"github.com/velo-protocol/DRSv1/libs/convert"
	"strings"
)

func (lo *logic) ExportAccount(input *entity.ExportAccountInput) (*entity.ExportAccountOutput, error) {
	accountBytes, err := lo.DB.Get([]byte(input.PublicKey))
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return nil, errors.Wrapf(err, "failed to get account from db")
	}
	if accountBytes == nil {
		return nil, errors.Errorf("account %s does not exist", input.PublicKey)
	}

	var account entity.StellarAccount
	err = json.Unmarshal(accountBytes, &account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal account")
	}

	seedBytes, err := crypto.Decrypt(account.EncryptedSeed, input.Passphrase)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decrypt the seed of %s with given passphrase", input.PublicKey)
	}

	keyPair, err := vconvert.SecretKeyToKeyPair(string(seedBytes))
	if err != nil {
		return nil, errors.Wrap(err, "failed to derive keypair from seed")
	}

	return &entity.ExportAccountOutput{
		ExportedKeyPair: keyPair,
	}, nil
}

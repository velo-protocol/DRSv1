package logic

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/crypto"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/parser"
	"gitlab.com/velo-labs/cen/libs/convert"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
)

func (lo *logic) RebalanceReserve(input *entity.RebalanceInput) (*entity.RebalanceOutput, error) {
	defaultAccount := lo.AppConfig.GetDefaultAccount()
	accountBytes, err := lo.DB.Get([]byte(defaultAccount))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account from db")
	}

	var account entity.StellarAccount
	err = json.Unmarshal(accountBytes, &account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal account")
	}

	seedBytes, err := crypto.Decrypt(account.EncryptedSeed, input.Passphrase)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decrypt the seed of %s with given passphrase", defaultAccount)
	}

	keyPair, err := vconvert.SecretKeyToKeyPair(string(seedBytes))
	if err != nil {
		return nil, errors.Wrap(err, "failed to derive keypair from seed")
	}

	result, err := lo.Velo.Client(keyPair).RebalanceReserve(context.Background(), vtxnbuild.RebalanceReserve{})
	if err != nil {
		err = parser.ParseHorizonError(err, lo.AppConfig.GetHorizonUrl(), lo.AppConfig.GetNetworkPassphrase())
		return nil, errors.Wrap(err, "failed to rebalance reserve")
	}

	return &entity.RebalanceOutput{
		TxResult: result,
	}, nil
}

package logic

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/crypto"
	"gitlab.com/velo-labs/cen/libs/convert"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
)

func (lo *logic) MintCredit(input *entity.MintCreditInput) (*entity.MintCreditOutput, error) {
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

	result, err := lo.Velo.Client(keyPair).MintCredit(context.Background(), vtxnbuild.MintCredit{
		AssetCodeToBeIssued: input.AssetToBeMinted,
		CollateralAssetCode: input.CollateralAssetCode,
		CollateralAmount:    input.CollateralAmount,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to mint credit")
	}

	return &entity.MintCreditOutput{
		AssetToBeMinted:     input.AssetToBeMinted,
		CollateralAssetCode: input.CollateralAssetCode,
		CollateralAmount:    input.CollateralAmount,
		SourceAddress:       defaultAccount,
		TxResult:            result,
	}, nil
}
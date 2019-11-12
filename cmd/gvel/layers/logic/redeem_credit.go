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

func (lo *logic) RedeemCredit(input *entity.RedeemCreditInput) (*entity.RedeemCreditOutput, error) {
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

	result, err := lo.Velo.Client(keyPair).RedeemCredit(context.Background(), vtxnbuild.RedeemCredit{
		AssetCode: input.AssetCode,
		Issuer:    input.AssetIssuer,
		Amount:    input.Amount,
	})
	if err != nil {
		err = parser.ParseHorizonError(err, lo.AppConfig.GetHorizonUrl(), lo.AppConfig.GetNetworkPassphrase())
		return nil, errors.Wrap(err, "failed to redeem stable credit")
	}

	return &entity.RedeemCreditOutput{
		AssetCode:   input.AssetCode,
		AssetIssuer: input.AssetIssuer,
		Amount:      input.Amount,
		TxResult:    result.HorizonResult,
	}, nil
}

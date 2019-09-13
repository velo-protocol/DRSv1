package stellar

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"gitlab.com/velo-labs/cen/node/app/environments"
	"gitlab.com/velo-labs/cen/node/app/utils"
)

func (repo *repo) BuildMintTx(
	drsAccount *horizon.Account,
	amount string,
	assetName string,
	issuerAddress string,
	distributorAddress string,
) (string, error) {
	var txOps []txnbuild.Operation

	drsKP, err := utils.KpFromSeedString(env.DrsPrivateKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to derived KP from seed key")
	}

	issueAssetOp := txnbuild.Payment{
		Destination: distributorAddress,
		Amount:      amount,
		Asset: txnbuild.CreditAsset{
			Code:   assetName,
			Issuer: issuerAddress,
		},
		SourceAccount: &horizon.Account{
			AccountID: issuerAddress,
		},
	}
	txOps = append(txOps, &issueAssetOp)

	mintTx := txnbuild.Transaction{
		SourceAccount: drsAccount,
		Operations:    txOps,
		Network:       env.NetworkPassphrase,
		Timebounds:    txnbuild.NewTimeout(300),
	}

	mintTxB64, err := mintTx.BuildSignEncode(drsKP)
	if err != nil {
		return "", errors.Wrap(err, "failed to build and sign mintTx")
	}

	return mintTxB64, nil
}

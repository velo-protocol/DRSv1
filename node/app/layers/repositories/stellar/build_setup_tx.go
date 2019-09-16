package stellar

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/environments"
	"gitlab.com/velo-labs/cen/node/app/utils"
)

func (repo *repo) BuildSetupTx(
	drsAccount *horizon.Account,
	peggedValue string,
	peggedCurrency string,
	assetName string,
	creditOwnerAddress string,
) (setupTxB64 string, issuerAddress string, distributorAddress string, err error) {
	var txOps []txnbuild.Operation

	drsKP, err := utils.KpFromSeedString(env.DrsPrivateKey)
	if err != nil {
		return "", "", "", errors.Wrap(err, "failed to derived KP from seed key")
	}

	issuerKP, err := keypair.Random()
	if err != nil {
		return "", "", "", errors.Wrap(err, "failed to create issuer KP")
	}

	distributorKP, err := keypair.Random()
	if err != nil {
		return "", "", "", errors.Wrap(err, "failed to create distributor KP")
	}

	createIssuerOp := txnbuild.CreateAccount{
		Destination: issuerKP.Address(),
		Amount:      "3",
	}
	txOps = append(txOps, &createIssuerOp)

	createDistributorOp := txnbuild.CreateAccount{
		Destination: distributorKP.Address(),
		Amount:      "2",
	}
	txOps = append(txOps, &createDistributorOp)

	storePeggedValueOp := txnbuild.ManageData{
		SourceAccount: &horizon.Account{
			AccountID: issuerKP.Address(),
		},
		Name:  "peggedValue",
		Value: []byte(peggedValue),
	}
	txOps = append(txOps, &storePeggedValueOp)

	storePeggedCurrencyOp := txnbuild.ManageData{
		SourceAccount: &horizon.Account{
			AccountID: issuerKP.Address(),
		},
		Name:  "peggedCurrency",
		Value: []byte(peggedCurrency),
	}
	txOps = append(txOps, &storePeggedCurrencyOp)

	storeAssetNameOp := txnbuild.ManageData{
		SourceAccount: &horizon.Account{
			AccountID: issuerKP.Address(),
		},
		Name:  "assetName",
		Value: []byte(assetName),
	}
	txOps = append(txOps, &storeAssetNameOp)

	distTrustOp := txnbuild.ChangeTrust{
		Limit: constants.MaxTrustlineLimit,
		Line: txnbuild.CreditAsset{
			Code:   assetName,
			Issuer: issuerKP.Address(),
		},
		SourceAccount: &horizon.Account{
			AccountID: distributorKP.Address(),
		},
	}
	txOps = append(txOps, &distTrustOp)

	zeroThreshold := txnbuild.Threshold(0)
	oneThreshold := txnbuild.Threshold(1)
	twoThreshold := txnbuild.Threshold(2)

	setIssuerWeightOp := txnbuild.SetOptions{
		MasterWeight:    &zeroThreshold,
		LowThreshold:    &twoThreshold,
		MediumThreshold: &twoThreshold,
		HighThreshold:   &twoThreshold,
		SourceAccount: &horizon.Account{
			AccountID: issuerKP.Address(),
		},
	}
	txOps = append(txOps, &setIssuerWeightOp)

	addDrsIssuerSignerOp := txnbuild.SetOptions{
		Signer: &txnbuild.Signer{
			Address: env.DrsPublicKey,
			Weight:  oneThreshold,
		},
		SourceAccount: &horizon.Account{
			AccountID: issuerKP.Address(),
		},
	}
	txOps = append(txOps, &addDrsIssuerSignerOp)

	addCreditOwnerIssuerSignerOp := txnbuild.SetOptions{
		Signer: &txnbuild.Signer{
			Address: creditOwnerAddress,
			Weight:  oneThreshold,
		},
		SourceAccount: &horizon.Account{
			AccountID: issuerKP.Address(),
		},
	}
	txOps = append(txOps, &addCreditOwnerIssuerSignerOp)

	setDistributorWeightOp := txnbuild.SetOptions{
		MasterWeight:    &zeroThreshold,
		LowThreshold:    &oneThreshold,
		MediumThreshold: &oneThreshold,
		HighThreshold:   &oneThreshold,
		SourceAccount: &horizon.Account{
			AccountID: distributorKP.Address(),
		},
	}
	txOps = append(txOps, &setDistributorWeightOp)

	addCreditOwnerDistributorSignerOp := txnbuild.SetOptions{
		Signer: &txnbuild.Signer{
			Address: creditOwnerAddress,
			Weight:  oneThreshold,
		},
		SourceAccount: &horizon.Account{
			AccountID: distributorKP.Address(),
		},
	}
	txOps = append(txOps, &addCreditOwnerDistributorSignerOp)

	setupTx := txnbuild.Transaction{
		SourceAccount: drsAccount,
		Operations:    txOps,
		Network:       env.NetworkPassphrase,
		Timebounds:    txnbuild.NewTimeout(300),
	}

	setupTxB64, err = setupTx.BuildSignEncode(drsKP, distributorKP, issuerKP)
	if err != nil {
		return "", "", "", errors.Wrap(err, "failed to build and sign mintTx")
	}

	return setupTxB64, issuerKP.Address(), distributorKP.Address(), nil
}

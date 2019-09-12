package stellar_drsops

import (
	"github.com/interstellar/starlight/errors"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"gitlab.com/velo-labs/cen/app/constants"
	env "gitlab.com/velo-labs/cen/app/environments"
	"gitlab.com/velo-labs/cen/app/modules/stellar"
	"gitlab.com/velo-labs/cen/app/services/operation"
	"gitlab.com/velo-labs/cen/app/utils"
)

type ops struct {
	StellarRepository stellar.Repository
}

func (o *ops) Setup(
	peggedValue string,
	peggedCurrency string,
	assetName string,
	creditOwnerAddress string,
) (setupTxB64 string, issuerAddress string, distributorAddress string, err error) {
	var txops []txnbuild.Operation

	drsKP, err := utils.KpFromSeedString(env.DrsPrivkey)
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
	txops = append(txops, &createIssuerOp)

	createDistributorOp := txnbuild.CreateAccount{
		Destination: distributorKP.Address(),
		Amount:      "2",
	}
	txops = append(txops, &createDistributorOp)

	storePeggedValueOp := txnbuild.ManageData{
		SourceAccount: &horizon.Account{
			AccountID: issuerKP.Address(),
		},
		Name:  "peggedValue",
		Value: []byte(peggedValue),
	}
	txops = append(txops, &storePeggedValueOp)

	storePeggedCurrencyOp := txnbuild.ManageData{
		SourceAccount: &horizon.Account{
			AccountID: issuerKP.Address(),
		},
		Name:  "peggedCurrency",
		Value: []byte(peggedCurrency),
	}
	txops = append(txops, &storePeggedCurrencyOp)

	storeAssetNameOp := txnbuild.ManageData{
		SourceAccount: &horizon.Account{
			AccountID: issuerKP.Address(),
		},
		Name:  "assetName",
		Value: []byte(assetName),
	}
	txops = append(txops, &storeAssetNameOp)

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
	txops = append(txops, &distTrustOp)

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
	txops = append(txops, &setIssuerWeightOp)

	addDrsIssuerSignerOp := txnbuild.SetOptions{
		Signer: &txnbuild.Signer{
			Address: env.DrsAddress,
			Weight:  oneThreshold,
		},
		SourceAccount: &horizon.Account{
			AccountID: issuerKP.Address(),
		},
	}
	txops = append(txops, &addDrsIssuerSignerOp)

	addCreditOwnerIssuerSignerOp := txnbuild.SetOptions{
		Signer: &txnbuild.Signer{
			Address: creditOwnerAddress,
			Weight:  oneThreshold,
		},
		SourceAccount: &horizon.Account{
			AccountID: issuerKP.Address(),
		},
	}
	txops = append(txops, &addCreditOwnerIssuerSignerOp)

	setDistributorWeightOp := txnbuild.SetOptions{
		MasterWeight:    &zeroThreshold,
		LowThreshold:    &oneThreshold,
		MediumThreshold: &oneThreshold,
		HighThreshold:   &oneThreshold,
		SourceAccount: &horizon.Account{
			AccountID: distributorKP.Address(),
		},
	}
	txops = append(txops, &setDistributorWeightOp)

	addCreditOwnerDistributorSignerOp := txnbuild.SetOptions{
		Signer: &txnbuild.Signer{
			Address: creditOwnerAddress,
			Weight:  oneThreshold,
		},
		SourceAccount: &horizon.Account{
			AccountID: distributorKP.Address(),
		},
	}
	txops = append(txops, &addCreditOwnerDistributorSignerOp)

	drsAccount, err := o.StellarRepository.LoadAccount(env.DrsAddress)
	if err != nil {
		return "", "", "", errors.Wrap(err, "failed to load the drs account")
	}

	setupTx := txnbuild.Transaction{
		SourceAccount: drsAccount,
		Operations:    txops,
		Network:       env.NetworkPassphrase,
		Timebounds:    txnbuild.NewTimeout(300),
	}

	setupTxB64, err = setupTx.BuildSignEncode(drsKP, distributorKP, issuerKP)
	if err != nil {
		return "", "", "", errors.Wrap(err, "failed to build and sign mintTx")
	}

	return setupTxB64, issuerKP.Address(), distributorKP.Address(), nil
}

func (o *ops) Mint(
	amount string,
	assetName string,
	issuerAddress string,
	distributorAddress string,
) (string, error) {
	var txops []txnbuild.Operation

	drsKP, err := utils.KpFromSeedString(env.DrsPrivkey)
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
	txops = append(txops, &issueAssetOp)

	drsAccount, err := o.StellarRepository.LoadAccount(env.DrsAddress)
	if err != nil {
		return "", errors.Wrap(err, "failed to load the drs account")
	}

	mintTx := txnbuild.Transaction{
		SourceAccount: drsAccount,
		Operations:    txops,
		Network:       env.NetworkPassphrase,
		Timebounds:    txnbuild.NewTimeout(300),
	}

	mintTxB64, err := mintTx.BuildSignEncode(drsKP)
	if err != nil {
		return "", errors.Wrap(err, "failed to build and sign mintTx")
	}

	return mintTxB64, nil
}

func NewDrsOps(stellarRepository stellar.Repository) operation.Interface {
	return &ops{
		StellarRepository: stellarRepository,
	}
}

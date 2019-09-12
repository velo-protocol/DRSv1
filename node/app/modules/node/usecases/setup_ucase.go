package usecases

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
	"gitlab.com/velo-labs/cen/node/app/entities"
	env "gitlab.com/velo-labs/cen/node/app/environments"
)

func (uc *usecase) Setup(
	issuerCreationTx string,
	peggedValue string,
	peggedCurrency string,
	assetName string,
) (*entities.Credit, error) {
	var txe xdr.TransactionEnvelope

	err := xdr.SafeUnmarshalBase64(issuerCreationTx, &txe)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal the issuer creation tx")
	}

	if len(txe.Tx.Operations) != 1 {
		return nil, errors.New("issuer creation tx must contains only one operation")
	}

	if txe.Tx.Operations[0].Body.PaymentOp.Amount != 30000000 {
		return nil, errors.New("issuer creation tx must send 3 XLM to the DRS address")
	}

	if txe.Tx.Operations[0].Body.PaymentOp.Destination.Address() != env.DrsAddress {
		return nil, errors.New("issuer creation tx must send 3 XLM to the DRS address as a destination")
	}

	issuerCreationTxSuccess, err := uc.StellarRepository.SubmitTransaction(issuerCreationTx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to submit issuer creation tx")
	}

	setupTxB64, issuerAddress, distributorAddress, err := uc.Drsops.Setup(peggedValue, peggedCurrency, assetName, txe.Tx.SourceAccount.Address())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a setup tx")
	}

	setupTxSuccess, err := uc.StellarRepository.SubmitTransaction(setupTxB64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to submit setup txe to stellar")
	}

	creditEntity := &entities.Credit{
		CreditOwnerAddress:   txe.Tx.SourceAccount.Address(),
		IssuerCreationTxHash: issuerCreationTxSuccess.Hash,
		SetupTxHash:          setupTxSuccess.Hash,
		IssuerAddress:        issuerAddress,
		DistributorAddress:   distributorAddress,
		AssetName:            assetName,
		PeggedCurrency:       peggedCurrency,
		PeggedValue:          peggedValue,
	}

	err = uc.NodeRepository.SaveCredit(*creditEntity)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save credit")
	}

	return creditEntity, nil
}

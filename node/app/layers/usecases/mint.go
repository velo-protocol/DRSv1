package usecases

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	env "gitlab.com/velo-labs/cen/node/app/environments"
)

func (useCase *useCase) Mint(
	postCollateralXdr string,
	assetName string,
	mintAmount string,
) (*string, error) {
	var txe xdr.TransactionEnvelope

	err := xdr.SafeUnmarshalBase64(postCollateralXdr, &txe)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal the post collateral xdr")
	}

	if len(txe.Tx.Operations) != 1 {
		return nil, errors.Wrap(err, "post collateral xdr must contain only one operation")
	}

	var code, issuer string
	err = txe.Tx.Operations[0].Body.PaymentOp.Asset.Extract(nil, &code, &issuer)
	if err != nil {
		return nil, errors.Wrap(err, "unable to extract asset from paymentOp")
	}

	if code != constants.VeloAbv && issuer != env.VeloIssuerAddress {
		return nil, errors.Wrap(err, "paymentOp not posting VELO to the drs")
	}

	if txe.Tx.Operations[0].Body.PaymentOp.Destination.Address() != env.DrsAddress {
		return nil, errors.Wrap(err, "paymentOp not posting VELO to the drs")
	}

	//TODO: load issuerAddress and distributorAddress from leveldb

	return nil, nil
}

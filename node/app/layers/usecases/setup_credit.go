package usecases

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/amount"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"gitlab.com/velo-labs/cen/libs/convert"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/environments"
	"gitlab.com/velo-labs/cen/node/app/errors"
	"strings"
)

func (useCase *useCase) SetupCredit(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*string, nerrors.NodeError) {
	if err := veloTx.VeloOp.Validate(); err != nil {
		return nil, nerrors.ErrInvalidArgument{Message: err.Error()}
	}

	txSenderPublicKey := veloTx.TxEnvelope().VeloTx.SourceAccount.Address()
	txSenderKeyPair, err := vconvert.PublicKeyToKeyPair(txSenderPublicKey)
	if err != nil {
		return nil, nerrors.ErrInvalidArgument{Message: err.Error()}
	}
	if veloTx.TxEnvelope().Signatures == nil {
		return nil, nerrors.ErrUnAuthenticated{Message: constants.ErrSignatureNotFound}
	}
	if txSenderKeyPair.Hint() != veloTx.TxEnvelope().Signatures[0].Hint {
		return nil, nerrors.ErrUnAuthenticated{Message: constants.ErrSignatureNotMatchSourceAccount}
	}

	// get tx sender account
	txSenderAccount, err := useCase.StellarRepo.GetAccount(veloTx.SourceAccount.GetAccountID())
	if err != nil {
		return nil, nerrors.ErrNotFound{
			Message: errors.Wrap(err, "fail to get tx sender account").Error(),
		}
	}

	drsAccountData, err := useCase.StellarRepo.GetDrsAccountData()
	if err != nil {
		return nil, nerrors.ErrInternal{
			Message: errors.Wrap(err, "fail to get data of drs account").Error(),
		}
	}

	// validate trusted partner role
	trustedPartnerList, err := useCase.StellarRepo.GetAccountData(drsAccountData.TrustedPartnerListAddress)
	if err != nil {
		return nil, nerrors.ErrInternal{
			Message: errors.Wrap(err, "fail to get data of trusted partner list account").Error(),
		}
	}

	trustedPartnerMetaEncoded, ok := trustedPartnerList[txSenderKeyPair.Address()]
	if !ok {
		return nil, nerrors.ErrPermissionDenied{
			Message: fmt.Sprintf(constants.ErrFormatSignerNotHavePermission, constants.VeloOpSetupCredit),
		}
	}

	trustedPartnerMetaAddress, err := base64.StdEncoding.DecodeString(trustedPartnerMetaEncoded)
	if err != nil {
		return nil, nerrors.ErrInternal{
			Message: errors.Wrapf(err, `fail to decode data "%s`, trustedPartnerMetaEncoded).Error(),
		}
	}

	// get trusted partner meta
	trustedPartnerMeta, err := useCase.StellarRepo.GetAccountData(string(trustedPartnerMetaAddress))
	if err != nil {
		return nil, nerrors.ErrInternal{Message: err.Error()}
	}

	for key, _ := range trustedPartnerMeta {
		assetDetail := strings.Split(key, "_")
		if assetDetail[0] == veloTx.TxEnvelope().VeloTx.VeloOp.Body.SetupCreditOp.AssetCode {
			return nil, nerrors.ErrInternal{Message: "the issuing and distribution account for asset code to specified already"}
		}
	}

	signedTx, err := buildSetupTx(txSenderAccount, veloTx.TxEnvelope().VeloTx.VeloOp.Body.SetupCreditOp, &txnbuild.SimpleAccount{AccountID: string(trustedPartnerMetaAddress)})
	if err != nil {
		return nil, nerrors.ErrInternal{Message: err.Error()}
	}

	return &signedTx, nil
}

func buildSetupTx(trustedPartnerAccount *horizon.Account, setupCreditOp *vxdr.SetupCreditOp, trustedPartnerMetaAddress *txnbuild.SimpleAccount) (setupTxB64 string, err error) {
	drsKp, err := vconvert.SecretKeyToKeyPair(env.DrsSecretKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to derived KP from seed key")
	}

	issuerKp, err := keypair.Random()
	if err != nil {
		return "", errors.Wrap(err, "failed to create issuer KP")
	}

	distributorKp, err := keypair.Random()
	if err != nil {
		return "", errors.Wrap(err, "failed to create distributor KP")
	}

	tx := txnbuild.Transaction{
		SourceAccount: trustedPartnerAccount,
		Operations: []txnbuild.Operation{
			// Trusted Party must pay tx fee to DRS
			&txnbuild.Payment{
				Destination:   drsKp.Address(),
				Amount:        "6",
				Asset:         txnbuild.NativeAsset{},
				SourceAccount: trustedPartnerAccount,
			},
			// Create issuer & distributor account
			&txnbuild.CreateAccount{
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: drsKp.Address(),
				},
				Destination: issuerKp.Address(),
				Amount:      "3.5",
			},
			&txnbuild.CreateAccount{
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: drsKp.Address(),
				},
				Destination: distributorKp.Address(),
				Amount:      "2",
			},
			// Add metadata to issuer
			&txnbuild.ManageData{
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: issuerKp.Address(),
				},
				Name:  "peggedValue",
				Value: []byte(amount.String(setupCreditOp.PeggedValue)),
			},
			&txnbuild.ManageData{
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: issuerKp.Address(),
				},
				Name:  "peggedCurrency",
				Value: []byte(setupCreditOp.PeggedCurrency),
			},
			&txnbuild.ManageData{
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: issuerKp.Address(),
				},
				Name:  "assetCode",
				Value: []byte(setupCreditOp.AssetCode),
			},
			// Create trust line between distributor ans issuer account
			&txnbuild.ChangeTrust{
				Limit: constants.MaxTrustlineLimit,
				Line: txnbuild.CreditAsset{
					Code:   setupCreditOp.AssetCode,
					Issuer: issuerKp.Address(),
				},
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: distributorKp.Address(),
				},
			},
			// Add signer to issuer account
			&txnbuild.SetOptions{
				Signer: &txnbuild.Signer{
					Address: trustedPartnerAccount.GetAccountID(),
					Weight:  txnbuild.Threshold(1),
				},
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: issuerKp.Address(),
				},
			},
			&txnbuild.SetOptions{
				Signer: &txnbuild.Signer{
					Address: env.DrsPublicKey,
					Weight:  txnbuild.Threshold(1),
				},
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: issuerKp.Address(),
				},
			},
			// Add signer to distributor account
			&txnbuild.SetOptions{
				Signer: &txnbuild.Signer{
					Address: trustedPartnerAccount.GetAccountID(),
					Weight:  txnbuild.Threshold(1),
				},
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: distributorKp.Address(),
				},
			},
			// Set threshold for both account
			&txnbuild.SetOptions{
				MasterWeight:    txnbuild.NewThreshold(0),
				LowThreshold:    txnbuild.NewThreshold(2),
				MediumThreshold: txnbuild.NewThreshold(2),
				HighThreshold:   txnbuild.NewThreshold(2),
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: issuerKp.Address(),
				},
			},
			&txnbuild.SetOptions{
				MasterWeight: txnbuild.NewThreshold(0),
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: distributorKp.Address(),
				},
			},
			// Add meta data to trusted partner
			&txnbuild.Payment{
				Destination: trustedPartnerMetaAddress.AccountID,
				Amount:      "0.5",
				Asset:       txnbuild.NativeAsset{},
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: drsKp.Address(),
				},
			},
			&txnbuild.ManageData{
				Name:          setupCreditOp.AssetCode + "_" + issuerKp.Address(),
				Value:         []byte(distributorKp.Address()),
				SourceAccount: trustedPartnerMetaAddress,
			},
		},
		Network:    env.NetworkPassphrase,
		Timebounds: txnbuild.NewTimeout(env.StellarTxTimeBoundInMinutes * 60), // seconds
	}

	signedTxXdr, err := tx.BuildSignEncode(drsKp, distributorKp, issuerKp)
	if err != nil {
		return "", errors.Wrap(err, "failed to build and sign tx")
	}

	return signedTxXdr, nil
}

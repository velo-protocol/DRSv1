package usecases

import (
	"context"
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"
	"github.com/stellar/go/amount"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"gitlab.com/velo-labs/cen/libs/convert"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/environments"
	"gitlab.com/velo-labs/cen/node/app/errors"
	"gitlab.com/velo-labs/cen/node/app/utils"
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

	trustedPartnerEntity, err := useCase.WhitelistRepo.FindOneWhitelist(entities.WhiteListFilter{
		StellarPublicKey: pointer.ToString(txSenderPublicKey),
		RoleCode:         pointer.ToString(string(vxdr.RoleTrustedPartner)),
	})
	if err != nil {
		return nil, nerrors.ErrInternal{Message: err.Error()}
	}
	if trustedPartnerEntity == nil {
		return nil, nerrors.ErrPermissionDenied{
			Message: fmt.Sprintf(constants.ErrFormatSignerNotHavePermission, constants.VeloOpSetupCredit),
		}
	}

	trustedPartnerAccount, err := useCase.StellarRepo.LoadAccount(trustedPartnerEntity.StellarPublicKey)
	if err != nil {
		return nil, nerrors.ErrNotFound{Message: err.Error()}
	}

	signedTx, err := buildSetupTx(trustedPartnerAccount, veloTx.TxEnvelope().VeloTx.VeloOp.Body.SetupCreditOp)
	if err != nil {
		return nil, nerrors.ErrInternal{Message: err.Error()}
	}

	return &signedTx, nil
}

func buildSetupTx(trustedPartnerAccount *horizon.Account, setupCreditOp *vxdr.SetupCreditOp) (setupTxB64 string, err error) {
	drsKp, err := utils.KpFromSeedString(env.DrsSecretKey)
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
				Amount:        "5.5",
				Asset:         txnbuild.NativeAsset{},
				SourceAccount: trustedPartnerAccount,
			},
			// Create issuer & distributor account
			&txnbuild.CreateAccount{
				SourceAccount: &horizon.Account{
					AccountID: drsKp.Address(),
				},
				Destination: issuerKp.Address(),
				Amount:      "3.5",
			},
			&txnbuild.CreateAccount{
				SourceAccount: &horizon.Account{
					AccountID: drsKp.Address(),
				},
				Destination: distributorKp.Address(),
				Amount:      "2",
			},
			// Add metadata to issuer
			&txnbuild.ManageData{
				SourceAccount: &horizon.Account{
					AccountID: issuerKp.Address(),
				},
				Name:  "peggedValue",
				Value: []byte(amount.String(setupCreditOp.PeggedValue)),
			},
			&txnbuild.ManageData{
				SourceAccount: &horizon.Account{
					AccountID: issuerKp.Address(),
				},
				Name:  "peggedCurrency",
				Value: []byte(setupCreditOp.PeggedCurrency),
			},
			&txnbuild.ManageData{
				SourceAccount: &horizon.Account{
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
				SourceAccount: &horizon.Account{
					AccountID: distributorKp.Address(),
				},
			},
			// Add signer to issuer account
			&txnbuild.SetOptions{
				Signer: &txnbuild.Signer{
					Address: trustedPartnerAccount.GetAccountID(),
					Weight:  txnbuild.Threshold(1),
				},
				SourceAccount: &horizon.Account{
					AccountID: issuerKp.Address(),
				},
			},
			&txnbuild.SetOptions{
				Signer: &txnbuild.Signer{
					Address: env.DrsPublicKey,
					Weight:  txnbuild.Threshold(1),
				},
				SourceAccount: &horizon.Account{
					AccountID: issuerKp.Address(),
				},
			},
			// Add signer to distributor account
			&txnbuild.SetOptions{
				Signer: &txnbuild.Signer{
					Address: trustedPartnerAccount.GetAccountID(),
					Weight:  txnbuild.Threshold(1),
				},
				SourceAccount: &horizon.Account{
					AccountID: distributorKp.Address(),
				},
			},
			// Set threshold for both account
			&txnbuild.SetOptions{
				MasterWeight:    txnbuild.NewThreshold(0),
				LowThreshold:    txnbuild.NewThreshold(2),
				MediumThreshold: txnbuild.NewThreshold(2),
				HighThreshold:   txnbuild.NewThreshold(2),
				SourceAccount: &horizon.Account{
					AccountID: issuerKp.Address(),
				},
			},
			&txnbuild.SetOptions{
				MasterWeight: txnbuild.NewThreshold(0),
				SourceAccount: &horizon.Account{
					AccountID: distributorKp.Address(),
				},
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

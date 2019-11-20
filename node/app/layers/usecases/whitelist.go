package usecases

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"github.com/velo-protocol/DRSv1/libs/convert"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"github.com/velo-protocol/DRSv1/libs/xdr"
	"github.com/velo-protocol/DRSv1/node/app/constants"
	"github.com/velo-protocol/DRSv1/node/app/entities"
	"github.com/velo-protocol/DRSv1/node/app/environments"
	"github.com/velo-protocol/DRSv1/node/app/errors"
)

func (useCase *useCase) Whitelist(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*entities.WhitelistOutput, nerrors.NodeError) {
	if err := veloTx.VeloOp.Validate(); err != nil {
		return nil, nerrors.ErrInvalidArgument{Message: err.Error()}
	}

	txSenderPublicKey := veloTx.TxEnvelope().VeloTx.SourceAccount.Address()
	whitelistOp := veloTx.TxEnvelope().VeloTx.VeloOp.Body.WhitelistOp

	// additional parameter validation
	if whitelistOp.Role == vxdr.RolePriceFeeder && whitelistOp.Currency == "" {
		return nil, nerrors.ErrInvalidArgument{
			Message: constants.ErrPriceFeederCurrencyMustNotBlank,
		}
	} else if whitelistOp.Role != vxdr.RolePriceFeeder && whitelistOp.Currency != "" {
		return nil, nerrors.ErrInvalidArgument{
			Message: constants.ErrCurrencyMustBeBlank,
		}
	}

	// validate tx signature
	txSenderKeyPair, err := vconvert.PublicKeyToKeyPair(txSenderPublicKey)
	if err != nil {
		return nil, nerrors.ErrUnAuthenticated{Message: err.Error()}
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
			Message: errors.Wrap(err, constants.ErrGetSenderAccount).Error(),
		}
	}

	// get drs account
	drsAccountData, err := useCase.StellarRepo.GetDrsAccountData()
	if err != nil {
		return nil, nerrors.ErrInternal{
			Message: errors.Wrap(err, constants.ErrGetDrsAccountData).Error(),
		}
	}

	// get lists of whitelisted account of each role
	accounts, err := useCase.StellarRepo.GetAccounts(
		drsAccountData.RegulatorListAddress,
		drsAccountData.TrustedPartnerListAddress,
		drsAccountData.PriceFeederListAddress,
	)
	if err != nil {
		return nil, nerrors.ErrInternal{
			Message: errors.Wrap(err, constants.ErrGetRoleListAccount).Error(),
		}
	}
	var (
		regulatorListData      = accounts[0].Data
		trustedPartnerListData = accounts[1].Data
		priceFeederListData    = accounts[2].Data
	)

	// validate tx sender role, in which must be regulator
	if _, ok := regulatorListData[txSenderKeyPair.Address()]; !ok {
		return nil, nerrors.ErrPermissionDenied{
			Message: fmt.Sprintf(constants.ErrFormatSignerNotHavePermission, constants.VeloOpWhitelist),
		}
	}

	// prepare tx
	drsKp, err := vconvert.SecretKeyToKeyPair(env.DrsSecretKey)
	if err != nil {
		return nil, nerrors.ErrInternal{
			Message: errors.Wrap(err, constants.ErrDerivedKeyPairFromSeed).Error(),
		}
	}

	switch whitelistOp.Role {
	case vxdr.RoleRegulator:
		// duplication check
		if _, ok := regulatorListData[whitelistOp.Address.Address()]; ok {
			return nil, nerrors.ErrAlreadyExists{
				Message: fmt.Sprintf(constants.ErrWhitelistAlreadyWhitelisted, whitelistOp.Address.Address(), vxdr.RoleMap[whitelistOp.Role]),
			}
		}

		tx := txnbuild.Transaction{
			SourceAccount: txSenderAccount,
			Operations: []txnbuild.Operation{
				// Regulator must pay tx fee to DRS
				&txnbuild.Payment{
					Destination:   drsKp.Address(),
					Amount:        "0.5",
					Asset:         txnbuild.NativeAsset{},
					SourceAccount: txSenderAccount,
				},
				// DRS pay to RegulatorList for account reserve
				&txnbuild.Payment{
					Destination: drsAccountData.RegulatorListAddress,
					Amount:      "0.5",
					Asset:       txnbuild.NativeAsset{},
					SourceAccount: &txnbuild.SimpleAccount{
						AccountID: drsKp.Address(),
					},
				},
				// Add the new address to the RegulatorList account
				&txnbuild.ManageData{
					Name:  whitelistOp.Address.Address(),
					Value: []byte("true"),
					SourceAccount: &txnbuild.SimpleAccount{
						AccountID: drsAccountData.RegulatorListAddress,
					},
				},
			},
			Network:    env.NetworkPassphrase,
			Timebounds: txnbuild.NewTimeout(env.StellarTxTimeBoundInMinutes),
		}
		signedTxXdr, err := tx.BuildSignEncode(drsKp)
		if err != nil {
			return nil, nerrors.ErrInternal{
				Message: errors.Wrap(err, constants.ErrBuildAndSignTransaction).Error(),
			}
		}
		return &entities.WhitelistOutput{
			SignedStellarTxXdr: signedTxXdr,
			Address:            whitelistOp.Address.Address(),
			Role:               string(vxdr.RoleRegulator),
		}, nil

	case vxdr.RoleTrustedPartner:
		// duplication check
		if _, ok := trustedPartnerListData[whitelistOp.Address.Address()]; ok {
			return nil, nerrors.ErrAlreadyExists{
				Message: fmt.Sprintf(constants.ErrWhitelistAlreadyWhitelisted, whitelistOp.Address.Address(), vxdr.RoleMap[whitelistOp.Role]),
			}
		}

		trustedPartnerMetaKp, err := keypair.Random()
		if err != nil {
			return nil, nerrors.ErrInternal{
				Message: errors.Wrap(err, constants.ErrCreateTrustedPartnerMetaKeyPair).Error(),
			}
		}

		tx := txnbuild.Transaction{
			SourceAccount: txSenderAccount,
			Operations: []txnbuild.Operation{
				// Regulator must pay tx fee to DRS
				&txnbuild.Payment{
					Destination:   drsKp.Address(),
					Amount:        "2",
					Asset:         txnbuild.NativeAsset{},
					SourceAccount: txSenderAccount,
				},
				// DRS pay to TrustedPartnerList for account reserve
				&txnbuild.Payment{
					Destination: drsAccountData.TrustedPartnerListAddress,
					Amount:      "0.5",
					Asset:       txnbuild.NativeAsset{},
					SourceAccount: &txnbuild.SimpleAccount{
						AccountID: drsKp.Address(),
					},
				},
				// Add the new address to the TrustedPartnerList account
				&txnbuild.ManageData{
					Name:  whitelistOp.Address.Address(),
					Value: []byte(trustedPartnerMetaKp.Address()),
					SourceAccount: &txnbuild.SimpleAccount{
						AccountID: drsAccountData.TrustedPartnerListAddress,
					},
				},
				// DRS create a TrustedPartnerMeta account
				&txnbuild.CreateAccount{
					SourceAccount: &txnbuild.SimpleAccount{
						AccountID: drsKp.Address(),
					},
					Destination: trustedPartnerMetaKp.Address(),
					Amount:      "1.5",
				},
				// Add signer and drop master key
				&txnbuild.SetOptions{
					MasterWeight: txnbuild.NewThreshold(0),
					Signer: &txnbuild.Signer{
						Address: drsKp.Address(),
						Weight:  txnbuild.Threshold(1),
					},
					SourceAccount: &txnbuild.SimpleAccount{
						AccountID: trustedPartnerMetaKp.Address(),
					},
				},
			},
			Network:    env.NetworkPassphrase,
			Timebounds: txnbuild.NewTimeout(env.StellarTxTimeBoundInMinutes),
		}

		signedTxXdr, err := tx.BuildSignEncode(drsKp, trustedPartnerMetaKp)
		if err != nil {
			return nil, nerrors.ErrInternal{
				Message: errors.Wrap(err, constants.ErrBuildAndSignTransaction).Error(),
			}
		}
		return &entities.WhitelistOutput{
			SignedStellarTxXdr:        signedTxXdr,
			Address:                   whitelistOp.Address.Address(),
			Role:                      string(vxdr.RoleTrustedPartner),
			TrustedPartnerMetaAddress: trustedPartnerMetaKp.Address(),
		}, nil
	case vxdr.RolePriceFeeder:
		// duplication check
		if _, ok := priceFeederListData[whitelistOp.Address.Address()]; ok {
			return nil, nerrors.ErrAlreadyExists{
				Message: fmt.Sprintf(constants.ErrWhitelistAlreadyWhitelisted, whitelistOp.Address.Address(), vxdr.RoleMap[whitelistOp.Role]),
			}
		}

		tx := txnbuild.Transaction{
			SourceAccount: txSenderAccount,
			Operations: []txnbuild.Operation{
				// Regulator must pay tx fee to DRS
				&txnbuild.Payment{
					Destination:   drsKp.Address(),
					Amount:        "1.5",
					Asset:         txnbuild.NativeAsset{},
					SourceAccount: txSenderAccount,
				},
				// DRS pay to PriceFeederList for account reserve
				&txnbuild.Payment{
					Destination: drsAccountData.PriceFeederListAddress,
					Amount:      "0.5",
					Asset:       txnbuild.NativeAsset{},
					SourceAccount: &txnbuild.SimpleAccount{
						AccountID: drsKp.Address(),
					},
				},
				// DRS pay to PriceAccount for account reserve
				&txnbuild.Payment{
					Destination: drsAccountData.VeloPriceAddress(whitelistOp.Currency),
					Amount:      "1",
					Asset:       txnbuild.NativeAsset{},
					SourceAccount: &txnbuild.SimpleAccount{
						AccountID: drsKp.Address(),
					},
				},
				// Add the new address to the PriceFeederList account
				&txnbuild.ManageData{
					Name:  whitelistOp.Address.Address(),
					Value: []byte(whitelistOp.Currency),
					SourceAccount: &txnbuild.SimpleAccount{
						AccountID: drsAccountData.PriceFeederListAddress,
					},
				},
				// Add signer to PriceAccount
				&txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: whitelistOp.Address.Address(),
						Weight:  txnbuild.Threshold(1),
					},
					SourceAccount: &txnbuild.SimpleAccount{
						AccountID: drsAccountData.VeloPriceAddress(whitelistOp.Currency),
					},
				},
			},
			Network:    env.NetworkPassphrase,
			Timebounds: txnbuild.NewTimeout(env.StellarTxTimeBoundInMinutes),
		}
		signedTxXdr, err := tx.BuildSignEncode(drsKp)
		if err != nil {
			return nil, nerrors.ErrInternal{
				Message: errors.Wrap(err, constants.ErrBuildAndSignTransaction).Error(),
			}
		}

		return &entities.WhitelistOutput{
			SignedStellarTxXdr: signedTxXdr,
			Address:            whitelistOp.Address.Address(),
			Role:               string(vxdr.RolePriceFeeder),
			Currency:           string(whitelistOp.Currency),
		}, nil
	default:
		return nil, nerrors.ErrInternal{
			Message: constants.ErrUnknownRoleSpecified,
		}
	}

}

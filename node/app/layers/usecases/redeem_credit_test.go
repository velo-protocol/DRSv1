package usecases_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/amount"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/environments"
	"gitlab.com/velo-labs/cen/node/app/errors"
	"testing"
)

func TestUseCase_RedeemCredit(t *testing.T) {

	var (
		redeemAmount = decimal.NewFromFloat(1000)
		peggedValue  = decimal.NewFromFloat(15000000)
		medianPrice  = decimal.NewFromFloat(23000000)

		vThbIssuerAddress = "GAN6D232HXTF4OHL7J36SAJD3M22H26B2O4QFVRO32OEM523KTMB6Q72"
		vThbAsset         = "vTHB"

		trustedPartnerAddress     = publicKey2
		trustedPartnerMetaAddress = publicKey3
		peggedCurrency            = "THB"

		getMockVeloTx = func() *vtxnbuild.VeloTx {
			return &vtxnbuild.VeloTx{
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: publicKey1,
				},
				VeloOp: &vtxnbuild.RedeemCredit{
					AssetCode: vThbAsset,
					Issuer:    vThbIssuerAddress,
					Amount:    redeemAmount.String(),
				},
			}
		}
	)

	t.Run("success", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccount(vThbIssuerAddress).
			Return(&horizon.Account{
				AccountID: vThbIssuerAddress,
				Sequence:  "1",
				Signers: []horizon.Signer{{
					Key:    env.DrsPublicKey,
					Weight: 1,
				}, {
					Key:    trustedPartnerAddress,
					Weight: 1,
				}, {
					Key:    vThbIssuerAddress,
					Weight: 0,
				}},
				Data: map[string]string{
					"peggedValue":    base64.StdEncoding.EncodeToString([]byte(peggedValue.String())),
					"peggedCurrency": base64.StdEncoding.EncodeToString([]byte(peggedCurrency)),
				},
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccount(trustedPartnerAddress).
			Return(&horizon.Account{
				AccountID: trustedPartnerAddress,
				Sequence:  "1",
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountDecodedDataByKey(drsAccountDataEnity.TrustedPartnerListAddress, trustedPartnerAddress).
			Return(trustedPartnerMetaAddress, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountDecodedDataByKey(trustedPartnerMetaAddress, fmt.Sprintf("%s_%s", vThbAsset, vThbIssuerAddress)).
			Return(fmt.Sprintf("%s_%s", vThbAsset, vThbIssuerAddress), nil)
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.New(medianPrice.IntPart(), -7), nil)

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		output, nErr := helper.useCase.RedeemCredit(context.Background(), veloTx)
		assert.NoError(t, nErr)
		assert.NotEmpty(t, output.SignedStellarTxXdr)

		stellarTx, err := txnbuild.TransactionFromXDR(output.SignedStellarTxXdr)
		assert.NoError(t, err)

		txEnv := stellarTx.TxEnvelope()
		assert.Equal(t, publicKey1, txEnv.Tx.SourceAccount.Address())

		assert.Equal(t, xdr.OperationTypePayment, txEnv.Tx.Operations[0].Body.Type)
		assert.Equal(t, publicKey1, txEnv.Tx.Operations[0].SourceAccount.Address())
		assert.Equal(t, vThbIssuerAddress, txEnv.Tx.Operations[0].Body.PaymentOp.Asset.AlphaNum4.Issuer.Address())
		assert.Equal(t, vThbAsset, func() string {
			bytes, _ := txEnv.Tx.Operations[0].Body.PaymentOp.Asset.AlphaNum4.AssetCode.MarshalBinary()
			return string(bytes)
		}())
		assert.Equal(t, vThbIssuerAddress, txEnv.Tx.Operations[0].Body.PaymentOp.Destination.Address())
		assert.Equal(t, redeemAmount.StringFixed(7), amount.String(txEnv.Tx.Operations[0].Body.PaymentOp.Amount))

		assert.Equal(t, xdr.OperationTypePayment, txEnv.Tx.Operations[1].Body.Type)
		assert.Equal(t, env.DrsPublicKey, txEnv.Tx.Operations[1].SourceAccount.Address())
		assert.Equal(t, env.VeloIssuerPublicKey, txEnv.Tx.Operations[1].Body.PaymentOp.Asset.AlphaNum4.Issuer.Address())
		assert.Equal(t, "VELO", func() string {
			bytes, _ := txEnv.Tx.Operations[1].Body.PaymentOp.Asset.AlphaNum4.AssetCode.MarshalBinary()
			return string(bytes)
		}())
		assert.Equal(t, publicKey1, txEnv.Tx.Operations[1].Body.PaymentOp.Destination.Address())
		assert.Equal(t, "652.1739130", amount.String(txEnv.Tx.Operations[1].Body.PaymentOp.Amount))
	})

	t.Run("Error - velo op validation fail", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.RedeemCredit{
				AssetCode: vThbAsset,
				Issuer:    vThbIssuerAddress,
			},
		}

		_, err := helper.useCase.RedeemCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.IsType(t, nerrors.ErrInvalidArgument{}, err)
	})

	t.Run("Error - VeloTx missing signer", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign()

		_, err := helper.useCase.RedeemCredit(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrSignatureNotFound)
		assert.IsType(t, nerrors.ErrUnAuthenticated{}, err)
	})

	t.Run("Error - VeloTx wrong signer", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp2)

		_, err := helper.useCase.RedeemCredit(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrSignatureNotMatchSourceAccount)
		assert.IsType(t, nerrors.ErrUnAuthenticated{}, err)
	})

	t.Run("Error - fail to get tx sender account", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(nil, errors.New("stellar return error"))

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		_, err := helper.useCase.RedeemCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrGetSenderAccount)
		assert.IsType(t, nerrors.ErrNotFound{}, err)
	})

	t.Run("Error - fail to get issuer account", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccount(vThbIssuerAddress).
			Return(nil, errors.New("some error has occurred"))

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		_, err := helper.useCase.RedeemCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrGetIssuerAccount)
		assert.IsType(t, nerrors.ErrNotFound{}, err)
	})

	t.Run("Error - signer count must be 3", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccount(vThbIssuerAddress).
			Return(&horizon.Account{
				AccountID: vThbIssuerAddress,
				Sequence:  "1",
				Signers:   []horizon.Signer{},
			}, nil)

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		_, err := helper.useCase.RedeemCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf(constants.ErrInvalidIssuerAccount, "signer count must be 3"))
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
	})

	t.Run("Error - invalid pegged value format, cannot parse base64", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccount(vThbIssuerAddress).
			Return(&horizon.Account{
				AccountID: vThbIssuerAddress,
				Sequence:  "1",
				Signers: []horizon.Signer{{
					Key:    env.DrsPublicKey,
					Weight: 1,
				}, {
					Key:    trustedPartnerAddress,
					Weight: 1,
				}, {
					Key:    vThbIssuerAddress,
					Weight: 0,
				}},
				Data: map[string]string{
					"peggedValue":    "BAD_VALUE",
					"peggedCurrency": "BAD_VALUE",
				},
			}, nil)

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		_, err := helper.useCase.RedeemCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf(constants.ErrInvalidIssuerAccount, "invalid pegged value format"))
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
	})

	t.Run("Error - invalid pegged value format", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccount(vThbIssuerAddress).
			Return(&horizon.Account{
				AccountID: vThbIssuerAddress,
				Sequence:  "1",
				Signers: []horizon.Signer{{
					Key:    env.DrsPublicKey,
					Weight: 1,
				}, {
					Key:    trustedPartnerAddress,
					Weight: 1,
				}, {
					Key:    vThbIssuerAddress,
					Weight: 0,
				}},
				Data: map[string]string{
					"peggedValue":    base64.StdEncoding.EncodeToString([]byte("BAD_VALUE")),
					"peggedCurrency": "BAD_VALUE",
				},
			}, nil)

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		_, err := helper.useCase.RedeemCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf(constants.ErrInvalidIssuerAccount, "invalid pegged value format"))
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
	})

	t.Run("Error - invalid pegged currency format", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccount(vThbIssuerAddress).
			Return(&horizon.Account{
				AccountID: vThbIssuerAddress,
				Sequence:  "1",
				Signers: []horizon.Signer{{
					Key:    env.DrsPublicKey,
					Weight: 1,
				}, {
					Key:    trustedPartnerAddress,
					Weight: 1,
				}, {
					Key:    vThbIssuerAddress,
					Weight: 0,
				}},
				Data: map[string]string{
					"peggedValue":    base64.StdEncoding.EncodeToString([]byte(peggedValue.String())),
					"peggedCurrency": "BAD_VALUE",
				},
			}, nil)

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		_, err := helper.useCase.RedeemCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf(constants.ErrInvalidIssuerAccount, "invalid pegged currency format"))
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
	})

	t.Run("Error - no drs account as a signer in issuer account", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccount(vThbIssuerAddress).
			Return(&horizon.Account{
				AccountID: vThbIssuerAddress,
				Sequence:  "1",
				Signers: []horizon.Signer{{
					Key:    "UNKNOWN_ADDRESS",
					Weight: 1,
				}, {
					Key:    trustedPartnerAddress,
					Weight: 1,
				}, {
					Key:    vThbIssuerAddress,
					Weight: 0,
				}},
				Data: map[string]string{
					"peggedValue":    base64.StdEncoding.EncodeToString([]byte(peggedValue.String())),
					"peggedCurrency": base64.StdEncoding.EncodeToString([]byte(peggedCurrency)),
				},
			}, nil)
		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		_, err := helper.useCase.RedeemCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf(constants.ErrInvalidIssuerAccount, "no drs as signer"))
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
	})

	t.Run("Error - fail to get trusted partner account", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccount(vThbIssuerAddress).
			Return(&horizon.Account{
				AccountID: vThbIssuerAddress,
				Sequence:  "1",
				Signers: []horizon.Signer{{
					Key:    env.DrsPublicKey,
					Weight: 1,
				}, {
					Key:    trustedPartnerAddress,
					Weight: 1,
				}, {
					Key:    vThbIssuerAddress,
					Weight: 0,
				}},
				Data: map[string]string{
					"peggedValue":    base64.StdEncoding.EncodeToString([]byte(peggedValue.String())),
					"peggedCurrency": base64.StdEncoding.EncodeToString([]byte(peggedCurrency)),
				},
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccount(trustedPartnerAddress).
			Return(nil, errors.New("some error has occurred"))

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		_, err := helper.useCase.RedeemCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrGetTrustedPartnerAccountDetail)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
	})

	t.Run("Error - fail to get drs account", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccount(vThbIssuerAddress).
			Return(&horizon.Account{
				AccountID: vThbIssuerAddress,
				Sequence:  "1",
				Signers: []horizon.Signer{{
					Key:    env.DrsPublicKey,
					Weight: 1,
				}, {
					Key:    trustedPartnerAddress,
					Weight: 1,
				}, {
					Key:    vThbIssuerAddress,
					Weight: 0,
				}},
				Data: map[string]string{
					"peggedValue":    base64.StdEncoding.EncodeToString([]byte(peggedValue.String())),
					"peggedCurrency": base64.StdEncoding.EncodeToString([]byte(peggedCurrency)),
				},
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccount(trustedPartnerAddress).
			Return(&horizon.Account{
				AccountID: trustedPartnerAddress,
				Sequence:  "1",
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(nil, errors.New("some error has occurred"))

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		_, err := helper.useCase.RedeemCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrGetDrsAccountData)
		assert.IsType(t, nerrors.ErrInternal{}, err)
	})

	t.Run("Error - fail to verify trusted partner account", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccount(vThbIssuerAddress).
			Return(&horizon.Account{
				AccountID: vThbIssuerAddress,
				Sequence:  "1",
				Signers: []horizon.Signer{{
					Key:    env.DrsPublicKey,
					Weight: 1,
				}, {
					Key:    trustedPartnerAddress,
					Weight: 1,
				}, {
					Key:    vThbIssuerAddress,
					Weight: 0,
				}},
				Data: map[string]string{
					"peggedValue":    base64.StdEncoding.EncodeToString([]byte(peggedValue.String())),
					"peggedCurrency": base64.StdEncoding.EncodeToString([]byte(peggedCurrency)),
				},
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccount(trustedPartnerAddress).
			Return(&horizon.Account{
				AccountID: trustedPartnerAddress,
				Sequence:  "1",
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountDecodedDataByKey(drsAccountDataEnity.TrustedPartnerListAddress, trustedPartnerAddress).
			Return("", errors.New("some error has occurred"))

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		_, err := helper.useCase.RedeemCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrVerifyTrustedPartnerAccount)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
	})

	t.Run("Error - fail to verify asset code", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccount(vThbIssuerAddress).
			Return(&horizon.Account{
				AccountID: vThbIssuerAddress,
				Sequence:  "1",
				Signers: []horizon.Signer{{
					Key:    env.DrsPublicKey,
					Weight: 1,
				}, {
					Key:    trustedPartnerAddress,
					Weight: 1,
				}, {
					Key:    vThbIssuerAddress,
					Weight: 0,
				}},
				Data: map[string]string{
					"peggedValue":    base64.StdEncoding.EncodeToString([]byte(peggedValue.String())),
					"peggedCurrency": base64.StdEncoding.EncodeToString([]byte(peggedCurrency)),
				},
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccount(trustedPartnerAddress).
			Return(&horizon.Account{
				AccountID: trustedPartnerAddress,
				Sequence:  "1",
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountDecodedDataByKey(drsAccountDataEnity.TrustedPartnerListAddress, trustedPartnerAddress).
			Return(trustedPartnerMetaAddress, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountDecodedDataByKey(trustedPartnerMetaAddress, fmt.Sprintf("%s_%s", vThbAsset, vThbIssuerAddress)).
			Return("", errors.New("some error has occurred"))

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		_, err := helper.useCase.RedeemCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrVerifyAssetCode)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
	})

	t.Run("Error - fail to get price of pegged currency", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccount(vThbIssuerAddress).
			Return(&horizon.Account{
				AccountID: vThbIssuerAddress,
				Sequence:  "1",
				Signers: []horizon.Signer{{
					Key:    env.DrsPublicKey,
					Weight: 1,
				}, {
					Key:    trustedPartnerAddress,
					Weight: 1,
				}, {
					Key:    vThbIssuerAddress,
					Weight: 0,
				}},
				Data: map[string]string{
					"peggedValue":    base64.StdEncoding.EncodeToString([]byte(peggedValue.String())),
					"peggedCurrency": base64.StdEncoding.EncodeToString([]byte(peggedCurrency)),
				},
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccount(trustedPartnerAddress).
			Return(&horizon.Account{
				AccountID: trustedPartnerAddress,
				Sequence:  "1",
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountDecodedDataByKey(drsAccountDataEnity.TrustedPartnerListAddress, trustedPartnerAddress).
			Return(trustedPartnerMetaAddress, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountDecodedDataByKey(trustedPartnerMetaAddress, fmt.Sprintf("%s_%s", vThbAsset, vThbIssuerAddress)).
			Return(fmt.Sprintf("%s_%s", vThbAsset, vThbIssuerAddress), nil)
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.Zero, errors.New("some error has occurred"))

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		_, err := helper.useCase.RedeemCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrGetPriceOfPeggedCurrency)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
	})

	t.Run("Error - median price must be greater than zero", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{
				AccountID: publicKey1,
				Sequence:  "1",
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccount(vThbIssuerAddress).
			Return(&horizon.Account{
				AccountID: vThbIssuerAddress,
				Sequence:  "1",
				Signers: []horizon.Signer{{
					Key:    env.DrsPublicKey,
					Weight: 1,
				}, {
					Key:    trustedPartnerAddress,
					Weight: 1,
				}, {
					Key:    vThbIssuerAddress,
					Weight: 0,
				}},
				Data: map[string]string{
					"peggedValue":    base64.StdEncoding.EncodeToString([]byte(peggedValue.String())),
					"peggedCurrency": base64.StdEncoding.EncodeToString([]byte(peggedCurrency)),
				},
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccount(trustedPartnerAddress).
			Return(&horizon.Account{
				AccountID: trustedPartnerAddress,
				Sequence:  "1",
			}, nil)
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountDecodedDataByKey(drsAccountDataEnity.TrustedPartnerListAddress, trustedPartnerAddress).
			Return(trustedPartnerMetaAddress, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccountDecodedDataByKey(trustedPartnerMetaAddress, fmt.Sprintf("%s_%s", vThbAsset, vThbIssuerAddress)).
			Return(fmt.Sprintf("%s_%s", vThbAsset, vThbIssuerAddress), nil)
		helper.mockStellarRepo.EXPECT().
			GetMedianPriceFromPriceAccount(drsAccountDataEnity.PriceThbVeloAddress).
			Return(decimal.NewFromFloat(-1.0), nil)

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		_, err := helper.useCase.RedeemCredit(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrMedianPriceMustBeGreaterThanZero)
		assert.IsType(t, nerrors.ErrPrecondition{}, err)
	})
}

package usecases_test

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/assert"
	vtxnbuild "gitlab.com/velo-labs/cen/libs/txnbuild"
	vxdr "gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	nerrors "gitlab.com/velo-labs/cen/node/app/errors"
	"testing"
)

func TestUseCase_CreateWhitelist(t *testing.T) {
	t.Run("Error - whitelist op validation fail", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.Whitelist{
				Address: publicKey2,
				Role:    "BAD_ROLE",
			},
		}

		_, err := helper.useCase.CreateWhitelist(context.Background(), veloTx)
		assert.IsType(t, nerrors.ErrInvalidArgument{}, err)
	})
	t.Run("Error - currency must not be blank for price feeder role", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.Whitelist{
				Address: publicKey2,
				Role:    string(vxdr.RolePriceFeeder),
			},
		}
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		_, err := helper.useCase.CreateWhitelist(context.Background(), veloTx)

		assert.EqualError(t, err, constants.ErrPriceFeederCurrencyMustNotBlank)
		assert.IsType(t, nerrors.ErrInvalidArgument{}, err)
	})
	t.Run("Error - currency must not be blank for price feeder role", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.Whitelist{
				Address:  publicKey2,
				Role:     string(vxdr.RoleRegulator),
				Currency: string(vxdr.CurrencyTHB),
			},
		}
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		_, err := helper.useCase.CreateWhitelist(context.Background(), veloTx)

		assert.EqualError(t, err, constants.ErrCurrencyMustBeBlank)
		assert.IsType(t, nerrors.ErrInvalidArgument{}, err)
	})
	t.Run("Error - signature not found", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.Whitelist{
				Address: publicKey2,
				Role:    string(vxdr.RoleRegulator),
			},
		}
		_ = veloTx.Build()

		_, err := helper.useCase.CreateWhitelist(context.Background(), veloTx)

		assert.EqualError(t, err, constants.ErrSignatureNotFound)
		assert.IsType(t, nerrors.ErrUnAuthenticated{}, err)
	})
	t.Run("Error - invalid signatures", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.Whitelist{
				Address: publicKey2,
				Role:    string(vxdr.RoleRegulator),
			},
		}
		_ = veloTx.Build()
		_ = veloTx.Sign(kp2)

		_, err := helper.useCase.CreateWhitelist(context.Background(), veloTx)

		assert.EqualError(t, err, constants.ErrSignatureNotMatchSourceAccount)
		assert.IsType(t, nerrors.ErrUnAuthenticated{}, err)
	})
	t.Run("Error - tx sender account not found", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.Whitelist{
				Address: publicKey2,
				Role:    string(vxdr.RoleRegulator),
			},
		}
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(nil, errors.New("some error has occurred"))

		_, err := helper.useCase.CreateWhitelist(context.Background(), veloTx)

		assert.IsType(t, nerrors.ErrNotFound{}, err)
	})
	t.Run("Error - fail to get drs account data", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.Whitelist{
				Address: publicKey2,
				Role:    string(vxdr.RoleRegulator),
			},
		}
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{AccountID: publicKey1, Sequence: "1"}, nil)
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(nil, errors.New("some error has occurred"))

		_, err := helper.useCase.CreateWhitelist(context.Background(), veloTx)

		assert.Contains(t, err.Error(), constants.ErrGetDrsAccount)
		assert.IsType(t, nerrors.ErrInternal{}, err)
	})
	t.Run("Error - fail to get role list accounts", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.Whitelist{
				Address: publicKey2,
				Role:    string(vxdr.RoleRegulator),
			},
		}
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{AccountID: publicKey1, Sequence: "1"}, nil)
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccounts(drsAccountDataEnity.RegulatorListAddress, drsAccountDataEnity.TrustedPartnerListAddress, drsAccountDataEnity.PriceFeederListAddress).
			Return(nil, errors.New("some error has occurred"))

		_, err := helper.useCase.CreateWhitelist(context.Background(), veloTx)

		assert.Contains(t, err.Error(), constants.ErrGetRoleListAccount)
		assert.IsType(t, nerrors.ErrInternal{}, err)
	})
	t.Run("Error - tx sender role validation fail", func(t *testing.T) {
		helper := initTest(t)
		defer helper.mockController.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.Whitelist{
				Address: publicKey2,
				Role:    string(vxdr.RoleRegulator),
			},
		}
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		helper.mockStellarRepo.EXPECT().
			GetAccount(publicKey1).
			Return(&horizon.Account{AccountID: publicKey1, Sequence: "1"}, nil)
		helper.mockStellarRepo.EXPECT().
			GetDrsAccountData().
			Return(&drsAccountDataEnity, nil)
		helper.mockStellarRepo.EXPECT().
			GetAccounts(drsAccountDataEnity.RegulatorListAddress, drsAccountDataEnity.TrustedPartnerListAddress, drsAccountDataEnity.PriceFeederListAddress).
			Return([]horizon.Account{
				{
					AccountID: drsAccountDataEnity.RegulatorListAddress,
					Data:      map[string]string{},
				},
				{
					AccountID: drsAccountDataEnity.TrustedPartnerListAddress,
					Data:      map[string]string{},
				},
				{
					AccountID: drsAccountDataEnity.PriceFeederListAddress,
					Data:      map[string]string{},
				},
			}, nil)

		_, err := helper.useCase.CreateWhitelist(context.Background(), veloTx)

		assert.EqualError(t, err, fmt.Sprintf(constants.ErrFormatSignerNotHavePermission, constants.VeloOpWhitelist))
		assert.IsType(t, nerrors.ErrPermissionDenied{}, err)
	})

	t.Run("When role == REGULATOR", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			helper := initTest(t)
			defer helper.mockController.Finish()

			veloTx := &vtxnbuild.VeloTx{
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: publicKey1,
				},
				VeloOp: &vtxnbuild.Whitelist{
					Address: publicKey2,
					Role:    string(vxdr.RoleRegulator),
				},
			}
			_ = veloTx.Build()
			_ = veloTx.Sign(kp1)

			helper.mockStellarRepo.EXPECT().
				GetAccount(publicKey1).
				Return(&horizon.Account{AccountID: publicKey1, Sequence: "1"}, nil)
			helper.mockStellarRepo.EXPECT().
				GetDrsAccountData().
				Return(&drsAccountDataEnity, nil)
			helper.mockStellarRepo.EXPECT().
				GetAccounts(drsAccountDataEnity.RegulatorListAddress, drsAccountDataEnity.TrustedPartnerListAddress, drsAccountDataEnity.PriceFeederListAddress).
				Return([]horizon.Account{
					{
						AccountID: drsAccountDataEnity.RegulatorListAddress,
						Data: map[string]string{
							publicKey1: base64.StdEncoding.EncodeToString([]byte("true")),
						},
					},
					{
						AccountID: drsAccountDataEnity.TrustedPartnerListAddress,
						Data:      map[string]string{},
					},
					{
						AccountID: drsAccountDataEnity.PriceFeederListAddress,
						Data:      map[string]string{},
					},
				}, nil)

			signedTxXdr, err := helper.useCase.CreateWhitelist(context.Background(), veloTx)

			assert.NoError(t, err)
			assert.NotNil(t, signedTxXdr)
		})
		t.Run("Error - public key 2 has already been whitelisted as a REGULATOR", func(t *testing.T) {
			helper := initTest(t)
			defer helper.mockController.Finish()

			veloTx := &vtxnbuild.VeloTx{
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: publicKey1,
				},
				VeloOp: &vtxnbuild.Whitelist{
					Address: publicKey2,
					Role:    string(vxdr.RoleRegulator),
				},
			}
			_ = veloTx.Build()
			_ = veloTx.Sign(kp1)

			helper.mockStellarRepo.EXPECT().
				GetAccount(publicKey1).
				Return(&horizon.Account{AccountID: publicKey1, Sequence: "1"}, nil)
			helper.mockStellarRepo.EXPECT().
				GetDrsAccountData().
				Return(&drsAccountDataEnity, nil)
			helper.mockStellarRepo.EXPECT().
				GetAccounts(drsAccountDataEnity.RegulatorListAddress, drsAccountDataEnity.TrustedPartnerListAddress, drsAccountDataEnity.PriceFeederListAddress).
				Return([]horizon.Account{
					{
						AccountID: drsAccountDataEnity.RegulatorListAddress,
						Data: map[string]string{
							publicKey1: base64.StdEncoding.EncodeToString([]byte("true")),
							publicKey2: base64.StdEncoding.EncodeToString([]byte("true")),
						},
					},
					{
						AccountID: drsAccountDataEnity.TrustedPartnerListAddress,
						Data:      map[string]string{},
					},
					{
						AccountID: drsAccountDataEnity.PriceFeederListAddress,
						Data:      map[string]string{},
					},
				}, nil)

			_, err := helper.useCase.CreateWhitelist(context.Background(), veloTx)

			assert.EqualError(t, err, fmt.Sprintf(constants.ErrWhitelistAlreadyWhitelisted, publicKey2, vxdr.RoleMap[vxdr.RoleRegulator]))
			assert.IsType(t, nerrors.ErrAlreadyExists{}, err)
		})
	})
	t.Run("When role == TRUSTED_PARTNER", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			helper := initTest(t)
			defer helper.mockController.Finish()

			veloTx := &vtxnbuild.VeloTx{
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: publicKey1,
				},
				VeloOp: &vtxnbuild.Whitelist{
					Address: publicKey2,
					Role:    string(vxdr.RoleTrustedPartner),
				},
			}
			_ = veloTx.Build()
			_ = veloTx.Sign(kp1)

			helper.mockStellarRepo.EXPECT().
				GetAccount(publicKey1).
				Return(&horizon.Account{AccountID: publicKey1, Sequence: "1"}, nil)
			helper.mockStellarRepo.EXPECT().
				GetDrsAccountData().
				Return(&drsAccountDataEnity, nil)
			helper.mockStellarRepo.EXPECT().
				GetAccounts(drsAccountDataEnity.RegulatorListAddress, drsAccountDataEnity.TrustedPartnerListAddress, drsAccountDataEnity.PriceFeederListAddress).
				Return([]horizon.Account{
					{
						AccountID: drsAccountDataEnity.RegulatorListAddress,
						Data: map[string]string{
							publicKey1: base64.StdEncoding.EncodeToString([]byte("true")),
						},
					},
					{
						AccountID: drsAccountDataEnity.TrustedPartnerListAddress,
						Data:      map[string]string{},
					},
					{
						AccountID: drsAccountDataEnity.PriceFeederListAddress,
						Data:      map[string]string{},
					},
				}, nil)

			signedTxXdr, err := helper.useCase.CreateWhitelist(context.Background(), veloTx)

			assert.NoError(t, err)
			assert.NotNil(t, signedTxXdr)
		})
		t.Run("Error - public key 2 has already been whitelisted as a TRUSTED_PARTNER", func(t *testing.T) {
			helper := initTest(t)
			defer helper.mockController.Finish()

			veloTx := &vtxnbuild.VeloTx{
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: publicKey1,
				},
				VeloOp: &vtxnbuild.Whitelist{
					Address: publicKey2,
					Role:    string(vxdr.RoleTrustedPartner),
				},
			}
			_ = veloTx.Build()
			_ = veloTx.Sign(kp1)

			helper.mockStellarRepo.EXPECT().
				GetAccount(publicKey1).
				Return(&horizon.Account{AccountID: publicKey1, Sequence: "1"}, nil)
			helper.mockStellarRepo.EXPECT().
				GetDrsAccountData().
				Return(&drsAccountDataEnity, nil)
			helper.mockStellarRepo.EXPECT().
				GetAccounts(drsAccountDataEnity.RegulatorListAddress, drsAccountDataEnity.TrustedPartnerListAddress, drsAccountDataEnity.PriceFeederListAddress).
				Return([]horizon.Account{
					{
						AccountID: drsAccountDataEnity.RegulatorListAddress,
						Data: map[string]string{
							publicKey1: base64.StdEncoding.EncodeToString([]byte("true")),
						},
					},
					{
						AccountID: drsAccountDataEnity.TrustedPartnerListAddress,
						Data: map[string]string{
							publicKey2: base64.StdEncoding.EncodeToString([]byte("PUBLIC_KEY_2_META_ADDRESS")),
						},
					},
					{
						AccountID: drsAccountDataEnity.PriceFeederListAddress,
						Data:      map[string]string{},
					},
				}, nil)

			_, err := helper.useCase.CreateWhitelist(context.Background(), veloTx)

			assert.EqualError(t, err, fmt.Sprintf(constants.ErrWhitelistAlreadyWhitelisted, publicKey2, vxdr.RoleMap[vxdr.RoleTrustedPartner]))
			assert.IsType(t, nerrors.ErrAlreadyExists{}, err)
		})
	})
	t.Run("When role == PRICE_FEEDER", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			helper := initTest(t)
			defer helper.mockController.Finish()

			veloTx := &vtxnbuild.VeloTx{
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: publicKey1,
				},
				VeloOp: &vtxnbuild.Whitelist{
					Address:  publicKey2,
					Role:     string(vxdr.RolePriceFeeder),
					Currency: string(vxdr.CurrencyTHB),
				},
			}
			_ = veloTx.Build()
			_ = veloTx.Sign(kp1)

			helper.mockStellarRepo.EXPECT().
				GetAccount(publicKey1).
				Return(&horizon.Account{AccountID: publicKey1, Sequence: "1"}, nil)
			helper.mockStellarRepo.EXPECT().
				GetDrsAccountData().
				Return(&drsAccountDataEnity, nil)
			helper.mockStellarRepo.EXPECT().
				GetAccounts(drsAccountDataEnity.RegulatorListAddress, drsAccountDataEnity.TrustedPartnerListAddress, drsAccountDataEnity.PriceFeederListAddress).
				Return([]horizon.Account{
					{
						AccountID: drsAccountDataEnity.RegulatorListAddress,
						Data: map[string]string{
							publicKey1: base64.StdEncoding.EncodeToString([]byte("true")),
						},
					},
					{
						AccountID: drsAccountDataEnity.TrustedPartnerListAddress,
						Data:      map[string]string{},
					},
					{
						AccountID: drsAccountDataEnity.PriceFeederListAddress,
						Data:      map[string]string{},
					},
				}, nil)

			signedTxXdr, err := helper.useCase.CreateWhitelist(context.Background(), veloTx)

			assert.NoError(t, err)
			assert.NotNil(t, signedTxXdr)
		})
		t.Run("Error - public key 2 has already been whitelisted as a PRICE_FEEDER", func(t *testing.T) {
			helper := initTest(t)
			defer helper.mockController.Finish()

			veloTx := &vtxnbuild.VeloTx{
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: publicKey1,
				},
				VeloOp: &vtxnbuild.Whitelist{
					Address:  publicKey2,
					Role:     string(vxdr.RolePriceFeeder),
					Currency: string(vxdr.CurrencyTHB),
				},
			}
			_ = veloTx.Build()
			_ = veloTx.Sign(kp1)

			helper.mockStellarRepo.EXPECT().
				GetAccount(publicKey1).
				Return(&horizon.Account{AccountID: publicKey1, Sequence: "1"}, nil)
			helper.mockStellarRepo.EXPECT().
				GetDrsAccountData().
				Return(&drsAccountDataEnity, nil)
			helper.mockStellarRepo.EXPECT().
				GetAccounts(drsAccountDataEnity.RegulatorListAddress, drsAccountDataEnity.TrustedPartnerListAddress, drsAccountDataEnity.PriceFeederListAddress).
				Return([]horizon.Account{
					{
						AccountID: drsAccountDataEnity.RegulatorListAddress,
						Data: map[string]string{
							publicKey1: base64.StdEncoding.EncodeToString([]byte("true")),
						},
					},
					{
						AccountID: drsAccountDataEnity.TrustedPartnerListAddress,
						Data:      map[string]string{},
					},
					{
						AccountID: drsAccountDataEnity.PriceFeederListAddress,
						Data: map[string]string{
							publicKey2: base64.StdEncoding.EncodeToString([]byte("PUBLIC_KEY_2_META_ADDRESS")),
						},
					},
				}, nil)

			_, err := helper.useCase.CreateWhitelist(context.Background(), veloTx)

			assert.EqualError(t, err, fmt.Sprintf(constants.ErrWhitelistAlreadyWhitelisted, publicKey2, vxdr.RoleMap[vxdr.RolePriceFeeder]))
			assert.IsType(t, nerrors.ErrAlreadyExists{}, err)
		})
	})

}

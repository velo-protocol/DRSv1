package usecases_test

import (
	"context"
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/libs/convert"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	nerrors "gitlab.com/velo-labs/cen/node/app/errors"
	"testing"
)

func TestUseCase_UpdatePrice(t *testing.T) {
	var (
		kp1, _ = vconvert.SecretKeyToKeyPair(secretKey1)
		kp2, _ = vconvert.SecretKeyToKeyPair(secretKey2)

		getMockVeloTx = func() *vtxnbuild.VeloTx {
			return &vtxnbuild.VeloTx{
				SourceAccount: &txnbuild.SimpleAccount{
					AccountID: publicKey1,
				},
				VeloOp: &vtxnbuild.PriceUpdate{
					Asset:                       "VELO",
					Currency:                    "THB",
					PriceInCurrencyPerAssetUnit: "1",
				},
			}
		}
	)

	t.Run("Happy", func(t *testing.T) {
		useCase, testHelper, mockCtrl := initUseCaseTest(t)
		defer mockCtrl.Finish()

		veloTx := getMockVeloTx()
		err := veloTx.Build()
		assert.NoError(t, err)
		_ = veloTx.Sign(kp1)

		testHelper.MockWhiteListRepo.EXPECT().
			FindOneWhitelist(entities.WhiteListFilter{
				StellarPublicKey: &publicKey1,
				RoleCode:         pointer.ToString(string(vxdr.RolePriceFeeder)),
			}).
			Return(&entities.WhiteList{
				StellarPublicKey: publicKey1,
				RoleCode:         string(vxdr.RolePriceFeeder),
			}, nil)

		createPriceEntry := &entities.CreatePriceEntry{
			FeederPublicKey:             publicKey1,
			Asset:                       veloTx.TxEnvelope().VeloTx.VeloOp.Body.PriceUpdateOp.Asset,
			PriceInCurrencyPerAssetUnit: decimal.New(int64(veloTx.TxEnvelope().VeloTx.VeloOp.Body.PriceUpdateOp.PriceInCurrencyPerAssetUnit), -7),
			Currency:                    string(veloTx.TxEnvelope().VeloTx.VeloOp.Body.PriceUpdateOp.Currency),
		}

		testHelper.MockPriceRepo.EXPECT().CreatePriceEntry(createPriceEntry).Return(createPriceEntry, nil)

		err = useCase.UpdatePrice(context.Background(), veloTx)
		assert.NoError(t, err)
	})

	t.Run("Error - invalid argument empty asset", func(t *testing.T) {
		useCase, _, mockCtrl := initUseCaseTest(t)
		defer mockCtrl.Finish()

		veloTx := &vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.PriceUpdate{
				Asset:                       "",
				Currency:                    "THB",
				PriceInCurrencyPerAssetUnit: "1",
			},
		}

		err := useCase.UpdatePrice(context.Background(), veloTx)
		assert.Error(t, err)
	})

	t.Run("Error - VeloTx missing signer", func(t *testing.T) {
		useCase, _, mockCtrl := initUseCaseTest(t)
		defer mockCtrl.Finish()

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign()

		err := useCase.UpdatePrice(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrSignatureNotFound)
	})

	t.Run("Error - VeloTx wrong signer", func(t *testing.T) {
		useCase, _, mockCtrl := initUseCaseTest(t)
		defer mockCtrl.Finish()

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp2)

		err := useCase.UpdatePrice(context.Background(), veloTx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrSignatureNotMatchSourceAccount)
	})

	t.Run("Error - can't query on whitelist table", func(t *testing.T) {
		useCase, testHelper, mockCtrl := initUseCaseTest(t)
		defer mockCtrl.Finish()

		testHelper.MockWhiteListRepo.EXPECT().
			FindOneWhitelist(entities.WhiteListFilter{
				StellarPublicKey: &publicKey1,
				RoleCode:         pointer.ToString(string(vxdr.RolePriceFeeder)),
			}).
			Return(nil, errors.New(constants.ErrToGetDataFromDatabase))

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		err := useCase.UpdatePrice(context.Background(), veloTx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), constants.ErrToGetDataFromDatabase)
	})

	t.Run("Error - this user has no permission", func(t *testing.T) {
		useCase, testHelper, mockCtrl := initUseCaseTest(t)
		defer mockCtrl.Finish()

		testHelper.MockWhiteListRepo.EXPECT().
			FindOneWhitelist(entities.WhiteListFilter{
				StellarPublicKey: &publicKey1,
				RoleCode:         pointer.ToString(string(vxdr.RolePriceFeeder)),
			}).
			Return(nil, nil)

		veloTx := getMockVeloTx()
		_ = veloTx.Build()
		_ = veloTx.Sign(kp1)

		err := useCase.UpdatePrice(context.Background(), veloTx)
		assert.Contains(t, err.Error(), fmt.Sprintf(constants.ErrFormatSignerNotHavePermission, constants.VeloOpPriceFeeder))
		assert.IsType(t, nerrors.ErrPermissionDenied{}, err)
	})

	t.Run("Error - fail to update price", func(t *testing.T) {
		useCase, testHelper, mockCtrl := initUseCaseTest(t)
		defer mockCtrl.Finish()

		veloTx := getMockVeloTx()
		err := veloTx.Build()
		assert.NoError(t, err)
		_ = veloTx.Sign(kp1)

		testHelper.MockWhiteListRepo.EXPECT().
			FindOneWhitelist(entities.WhiteListFilter{
				StellarPublicKey: &publicKey1,
				RoleCode:         pointer.ToString(string(vxdr.RolePriceFeeder)),
			}).
			Return(&entities.WhiteList{
				StellarPublicKey: publicKey1,
				RoleCode:         string(vxdr.RolePriceFeeder),
			}, nil)

		createPriceEntry := &entities.CreatePriceEntry{
			FeederPublicKey:             publicKey1,
			Asset:                       veloTx.TxEnvelope().VeloTx.VeloOp.Body.PriceUpdateOp.Asset,
			PriceInCurrencyPerAssetUnit: decimal.New(int64(veloTx.TxEnvelope().VeloTx.VeloOp.Body.PriceUpdateOp.PriceInCurrencyPerAssetUnit), -7),
			Currency:                    string(veloTx.TxEnvelope().VeloTx.VeloOp.Body.PriceUpdateOp.Currency),
		}

		testHelper.MockPriceRepo.EXPECT().CreatePriceEntry(createPriceEntry).Return(nil, errors.New("some error has occurred"))

		err = useCase.UpdatePrice(context.Background(), veloTx)
		assert.Error(t, err)
		assert.IsType(t, nerrors.ErrInternal{}, err)
	})
}

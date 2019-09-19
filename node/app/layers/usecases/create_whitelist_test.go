package usecases

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/assert"
	vconvert "gitlab.com/velo-labs/cen/libs/convert"
	vtxnbuild "gitlab.com/velo-labs/cen/libs/txnbuild"
	vxdr "gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/layers/mocks"
	"testing"
)

func TestUseCase_CreateWhiteList(t *testing.T) {
	const (
		publicKey1 = "GBVI3QZYXCWQBSGZ4TNJOHDZ5KZYGZOVSE46TVAYJYTMNCGW2PWLWO73"
		secretKey1 = "SBR25NMQRKQ4RLGNV5XB3MMQB4ADVYSMPGVBODQVJE7KPTDR6KGK3XMX"
		publicKey2 = "GC2ROYZQH5FTVEPQZF7CAB32SCJC7DWVKILDUAT5BCU5O7HEI7HFUB25"
		secretKey2 = "SCHQI345PYWHM2APNR4MN433HNCBS7VDUROOZKTYHZUBBTHI2YHNCJ4G"
	)

	var (
		kp1, _ = vconvert.SecretKeyToKeyPair(secretKey1)
		kp2, _ = vconvert.SecretKeyToKeyPair(secretKey2)

		newMockWhiteListRepo = func() (*mocks.MockWhiteListRepo, func()) {
			ctrl := gomock.NewController(t)
			mockedWhiteListRepo := mocks.NewMockWhiteListRepo(ctrl)
			return mockedWhiteListRepo, ctrl.Finish
		}
	)

	stellarPublicAddress := publicKey1
	roleCode := string(vxdr.RoleRegulator)

	t.Run("Success", func(t *testing.T) {
		mockedWhiteListRepo, finish := newMockWhiteListRepo()
		defer finish()

		findWhiteListEntity := entities.WhiteList{
			ID: "e13d778c-d2c8-452b-8ead-368d43447fcd",
			StellarPublicKey: publicKey1,
			RoleCode: string(vxdr.RoleRegulator),
		}

		roleEntity := entities.Role{
			ID: 1,
			Name: "Price feeder",
			Code: "PRICE_FEEDER",
		}

		createWhitelistEntity := entities.WhiteList{
			StellarPublicKey: publicKey2,
			RoleCode: string(vxdr.RolePriceFeeder),
		}

		filter := entities.WhiteListFilter{
			StellarPublicKey: &stellarPublicAddress,
			RoleCode: &roleCode,
		}

		mockedWhiteListRepo.EXPECT().FindOneWhitelist(filter).Return(&findWhiteListEntity, nil)

		mockedWhiteListRepo.EXPECT().FindOneRole(string(vxdr.RolePriceFeeder)).Return(&roleEntity, nil)

		mockedWhiteListRepo.EXPECT().CreateWhitelist(&createWhitelistEntity).Return(&createWhitelistEntity, nil)

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.WhiteList{
				Address: publicKey2,
				Role:    string(vxdr.RolePriceFeeder),
			},
		}).BuildSignEncode(kp1, kp2)

		veloTx, _ := vtxnbuild.TransactionFromXDR(veloTxB64)
		envelope := veloTx.TxEnvelope()

		useCase := Init(nil, mockedWhiteListRepo)
		err := useCase.CreateWhiteList(context.Background(), envelope)

		assert.NoError(t, err)
	})

	t.Run("Error - invalid signatures", func(t *testing.T) {
		mockedWhiteListRepo, finish := newMockWhiteListRepo()
		defer finish()

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.WhiteList{
				Address: publicKey2,
				Role:    string(vxdr.RolePriceFeeder),
			},
		}).BuildSignEncode(kp2)

		veloTx, _ := vtxnbuild.TransactionFromXDR(veloTxB64)
		envelope := veloTx.TxEnvelope()

		useCase := Init(nil, mockedWhiteListRepo)
		err := useCase.CreateWhiteList(context.Background(), envelope)

		assert.EqualError(t, err, "can't create white list: bad signature")
	})

	t.Run("Error - can't query on whitelist table", func(t *testing.T) {
		mockedWhiteListRepo, finish := newMockWhiteListRepo()
		defer finish()

		filter := entities.WhiteListFilter{
			StellarPublicKey: &stellarPublicAddress,
			RoleCode: &roleCode,
		}

		mockedWhiteListRepo.EXPECT().FindOneWhitelist(filter).Return(nil, constants.ErrToGetDataFromDatabase)

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.WhiteList{
				Address: publicKey2,
				Role:    string(vxdr.RolePriceFeeder),
			},
		}).BuildSignEncode(kp1, kp2)

		veloTx, _ := vtxnbuild.TransactionFromXDR(veloTxB64)
		envelope := veloTx.TxEnvelope()

		useCase := Init(nil, mockedWhiteListRepo)
		err := useCase.CreateWhiteList(context.Background(), envelope)

		assert.EqualError(t, err, "can't create white list: can't get data from database")
	})

	t.Run("Error - pass query on whitelist table and can't query on role table", func(t *testing.T) {
		mockedWhiteListRepo, finish := newMockWhiteListRepo()
		defer finish()

		findWhiteListEntity := entities.WhiteList{
			ID: "e13d778c-d2c8-452b-8ead-368d43447fcd",
			StellarPublicKey: publicKey1,
			RoleCode: string(vxdr.RoleRegulator),
		}

		filter := entities.WhiteListFilter{
			StellarPublicKey: &stellarPublicAddress,
			RoleCode: &roleCode,
		}

		mockedWhiteListRepo.EXPECT().FindOneWhitelist(filter).Return(&findWhiteListEntity, nil)
		mockedWhiteListRepo.EXPECT().FindOneRole(string(vxdr.RolePriceFeeder)).Return(nil, constants.ErrToGetDataFromDatabase)

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.WhiteList{
				Address: publicKey2,
				Role:    string(vxdr.RolePriceFeeder),
			},
		}).BuildSignEncode(kp1, kp2)

		veloTx, _ := vtxnbuild.TransactionFromXDR(veloTxB64)
		envelope := veloTx.TxEnvelope()

		useCase := Init(nil, mockedWhiteListRepo)
		err := useCase.CreateWhiteList(context.Background(), envelope)

		assert.EqualError(t, err, "can't create white list: can't get data from database")
	})

	t.Run("Error - source account don't have regulator role", func(t *testing.T) {
		mockedWhiteListRepo, finish := newMockWhiteListRepo()
		defer finish()

		filter := entities.WhiteListFilter{
			StellarPublicKey: &stellarPublicAddress,
			RoleCode: &roleCode,
		}

		mockedWhiteListRepo.EXPECT().FindOneWhitelist(filter).Return(nil, nil)

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.WhiteList{
				Address: publicKey2,
				Role:    string(vxdr.RolePriceFeeder),
			},
		}).BuildSignEncode(kp1, kp2)

		veloTx, _ := vtxnbuild.TransactionFromXDR(veloTxB64)
		envelope := veloTx.TxEnvelope()

		useCase := Init(nil, mockedWhiteListRepo)
		err := useCase.CreateWhiteList(context.Background(), envelope)

		assert.EqualError(t, err, "can't create white list: unauthorized to perform an action")
	})

	t.Run("Error - send whitelist to save but fill invalid role", func(t *testing.T) {
		mockedWhiteListRepo, finish := newMockWhiteListRepo()
		defer finish()

		findWhiteListEntity := entities.WhiteList{
			ID: "e13d778c-d2c8-452b-8ead-368d43447fcd",
			StellarPublicKey: publicKey1,
			RoleCode: string(vxdr.RoleRegulator),
		}

		filter := entities.WhiteListFilter{
			StellarPublicKey: &stellarPublicAddress,
			RoleCode: &roleCode,
		}

		mockedWhiteListRepo.EXPECT().FindOneWhitelist(filter).Return(&findWhiteListEntity, nil)
		mockedWhiteListRepo.EXPECT().FindOneRole(string(vxdr.RolePriceFeeder)).Return(nil, constants.ErrToGetDataFromDatabase)

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.WhiteList{
				Address: publicKey2,
				Role:    string(vxdr.RolePriceFeeder),
			},
		}).BuildSignEncode(kp1, kp2)

		veloTx, _ := vtxnbuild.TransactionFromXDR(veloTxB64)
		envelope := veloTx.TxEnvelope()

		useCase := Init(nil, mockedWhiteListRepo)
		err := useCase.CreateWhiteList(context.Background(), envelope)

		assert.EqualError(t, err, "can't create white list: can't get data from database")
	})

	t.Run("Error - can't save whitelist table", func(t *testing.T) {
		mockedWhiteListRepo, finish := newMockWhiteListRepo()
		defer finish()

		findWhiteListEntity := entities.WhiteList{
			ID: "e13d778c-d2c8-452b-8ead-368d43447fcd",
			StellarPublicKey: publicKey1,
			RoleCode: string(vxdr.RoleRegulator),
		}

		roleEntity := entities.Role{
			ID: 1,
			Name: "Price feeder",
			Code: "PRICE_FEEDER",
		}

		createWhitelistEntity := entities.WhiteList{
			StellarPublicKey: publicKey2,
			RoleCode: string(vxdr.RolePriceFeeder),
		}

		filter := entities.WhiteListFilter{
			StellarPublicKey: &stellarPublicAddress,
			RoleCode: &roleCode,
		}

		mockedWhiteListRepo.EXPECT().FindOneWhitelist(filter).Return(&findWhiteListEntity, nil)
		mockedWhiteListRepo.EXPECT().FindOneRole(string(vxdr.RolePriceFeeder)).Return(&roleEntity, nil)
		mockedWhiteListRepo.EXPECT().CreateWhitelist(&createWhitelistEntity).Return(nil, constants.ErrToSaveDatabase)

		veloTxB64, _ := (&vtxnbuild.VeloTx{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: publicKey1,
			},
			VeloOp: &vtxnbuild.WhiteList{
				Address: publicKey2,
				Role:    string(vxdr.RolePriceFeeder),
			},
		}).BuildSignEncode(kp1, kp2)

		veloTx, _ := vtxnbuild.TransactionFromXDR(veloTxB64)
		envelope := veloTx.TxEnvelope()

		useCase := Init(nil, mockedWhiteListRepo)
		err := useCase.CreateWhiteList(context.Background(), envelope)

		assert.EqualError(t, err, "can't create white list: can't save to database")
	})

}
package usecases_test

import (
	"testing"
)

func TestUseCase_CreateWhiteList(t *testing.T) {
	//const (
	//	publicKey1 = "GBVI3QZYXCWQBSGZ4TNJOHDZ5KZYGZOVSE46TVAYJYTMNCGW2PWLWO73"
	//	secretKey1 = "SBR25NMQRKQ4RLGNV5XB3MMQB4ADVYSMPGVBODQVJE7KPTDR6KGK3XMX"
	//	publicKey2 = "GC2ROYZQH5FTVEPQZF7CAB32SCJC7DWVKILDUAT5BCU5O7HEI7HFUB25"
	//	secretKey2 = "SCHQI345PYWHM2APNR4MN433HNCBS7VDUROOZKTYHZUBBTHI2YHNCJ4G"
	//)
	//
	//var (
	//	kp1, _ = vconvert.SecretKeyToKeyPair(secretKey1)
	//	kp2, _ = vconvert.SecretKeyToKeyPair(secretKey2)
	//
	//	newMockWhiteListRepo = func() (*mocks.MockWhiteListRepo, func()) {
	//		ctrl := gomock.NewController(t)
	//		mockedWhiteListRepo := mocks.NewMockWhiteListRepo(ctrl)
	//		return mockedWhiteListRepo, ctrl.Finish
	//	}
	//)
	//
	//stellarPublicAddress := publicKey1
	//roleCode := string(vxdr.RoleRegulator)
	//
	//t.Run("Success", func(t *testing.T) {
	//	mockedWhiteListRepo, finish := newMockWhiteListRepo()
	//	defer finish()
	//
	//	findWhiteListEntity := entities.WhiteList{
	//		ID:               "e13d778c-d2c8-452b-8ead-368d43447fcd",
	//		StellarPublicKey: publicKey1,
	//		RoleCode:         string(vxdr.RoleRegulator),
	//	}
	//
	//	roleEntity := entities.Role{
	//		ID:   1,
	//		Name: "Price feeder",
	//		Code: "PRICE_FEEDER",
	//	}
	//
	//	createWhitelistEntity := entities.WhiteList{
	//		StellarPublicKey: publicKey2,
	//		RoleCode:         string(vxdr.RolePriceFeeder),
	//	}
	//
	//	filter := entities.WhiteListFilter{
	//		StellarPublicKey: &stellarPublicAddress,
	//		RoleCode:         &roleCode,
	//	}
	//
	//	mockedWhiteListRepo.EXPECT().FindOneWhitelist(filter).Return(&findWhiteListEntity, nil)
	//
	//	mockedWhiteListRepo.EXPECT().FindOneRole(string(vxdr.RolePriceFeeder)).Return(&roleEntity, nil)
	//
	//	mockedWhiteListRepo.EXPECT().CreateWhitelist(&createWhitelistEntity).Return(&createWhitelistEntity, nil)
	//
	//	veloTx := &vtxnbuild.VeloTx{
	//		SourceAccount: &txnbuild.SimpleAccount{
	//			AccountID: publicKey1,
	//		},
	//		VeloOp: &vtxnbuild.WhiteList{
	//			Address: publicKey2,
	//			Role:    string(vxdr.RolePriceFeeder),
	//		},
	//	}
	//	_ = veloTx.Build()
	//	_ = veloTx.Sign(kp1)
	//
	//	useCase := usecases.Init(nil)
	//	err := useCase.CreateWhiteList(context.Background(), veloTx)
	//
	//	assert.Nil(t, err)
	//})
	//
	//t.Run("Error - invalid signatures", func(t *testing.T) {
	//	mockedWhiteListRepo, finish := newMockWhiteListRepo()
	//	defer finish()
	//
	//	veloTx := &vtxnbuild.VeloTx{
	//		SourceAccount: &txnbuild.SimpleAccount{
	//			AccountID: publicKey1,
	//		},
	//		VeloOp: &vtxnbuild.WhiteList{
	//			Address: publicKey2,
	//			Role:    string(vxdr.RolePriceFeeder),
	//		},
	//	}
	//	_ = veloTx.Build()
	//	_ = veloTx.Sign(kp2)
	//
	//	useCase := usecases.Init(nil)
	//	err := useCase.CreateWhiteList(context.Background(), veloTx)
	//
	//	assert.Equal(t, err.Error(), constants.ErrSignatureNotMatchSourceAccount)
	//})
	//
	//t.Run("Error - can't query on whitelist table", func(t *testing.T) {
	//	mockedWhiteListRepo, finish := newMockWhiteListRepo()
	//	defer finish()
	//
	//	filter := entities.WhiteListFilter{
	//		StellarPublicKey: &stellarPublicAddress,
	//		RoleCode:         &roleCode,
	//	}
	//
	//	mockedWhiteListRepo.EXPECT().FindOneWhitelist(filter).Return(nil, errors.New(constants.ErrToGetDataFromDatabase))
	//
	//	veloTx := &vtxnbuild.VeloTx{
	//		SourceAccount: &txnbuild.SimpleAccount{
	//			AccountID: publicKey1,
	//		},
	//		VeloOp: &vtxnbuild.WhiteList{
	//			Address: publicKey2,
	//			Role:    string(vxdr.RolePriceFeeder),
	//		},
	//	}
	//	_ = veloTx.Build()
	//	_ = veloTx.Sign(kp1)
	//
	//	useCase := usecases.Init(nil)
	//	err := useCase.CreateWhiteList(context.Background(), veloTx)
	//
	//	assert.Equal(t, err.Error(), constants.ErrToGetDataFromDatabase)
	//})
	//
	//t.Run("Error - pass query on whitelist table and can't query on role table", func(t *testing.T) {
	//	mockedWhiteListRepo, finish := newMockWhiteListRepo()
	//	defer finish()
	//
	//	findWhiteListEntity := entities.WhiteList{
	//		ID:               "e13d778c-d2c8-452b-8ead-368d43447fcd",
	//		StellarPublicKey: publicKey1,
	//		RoleCode:         string(vxdr.RoleRegulator),
	//	}
	//
	//	filter := entities.WhiteListFilter{
	//		StellarPublicKey: &stellarPublicAddress,
	//		RoleCode:         &roleCode,
	//	}
	//
	//	mockedWhiteListRepo.EXPECT().FindOneWhitelist(filter).Return(&findWhiteListEntity, nil)
	//	mockedWhiteListRepo.EXPECT().FindOneRole(string(vxdr.RolePriceFeeder)).Return(nil, errors.New(constants.ErrToGetDataFromDatabase))
	//
	//	veloTx := &vtxnbuild.VeloTx{
	//		SourceAccount: &txnbuild.SimpleAccount{
	//			AccountID: publicKey1,
	//		},
	//		VeloOp: &vtxnbuild.WhiteList{
	//			Address: publicKey2,
	//			Role:    string(vxdr.RolePriceFeeder),
	//		},
	//	}
	//	_ = veloTx.Build()
	//	_ = veloTx.Sign(kp1)
	//
	//	useCase := usecases.Init(nil)
	//	err := useCase.CreateWhiteList(context.Background(), veloTx)
	//
	//	assert.Equal(t, err.Error(), constants.ErrToGetDataFromDatabase)
	//})
	//
	//t.Run("Error - pass query on whitelist table and empty role on role table", func(t *testing.T) {
	//	mockedWhiteListRepo, finish := newMockWhiteListRepo()
	//	defer finish()
	//
	//	findWhiteListEntity := entities.WhiteList{
	//		ID:               "e13d778c-d2c8-452b-8ead-368d43447fcd",
	//		StellarPublicKey: publicKey1,
	//		RoleCode:         string(vxdr.RoleRegulator),
	//	}
	//
	//	filter := entities.WhiteListFilter{
	//		StellarPublicKey: &stellarPublicAddress,
	//		RoleCode:         &roleCode,
	//	}
	//
	//	mockedWhiteListRepo.EXPECT().FindOneWhitelist(filter).Return(&findWhiteListEntity, nil)
	//	mockedWhiteListRepo.EXPECT().FindOneRole(string(vxdr.RolePriceFeeder)).Return(nil, nil)
	//
	//	veloTx := &vtxnbuild.VeloTx{
	//		SourceAccount: &txnbuild.SimpleAccount{
	//			AccountID: publicKey1,
	//		},
	//		VeloOp: &vtxnbuild.WhiteList{
	//			Address: publicKey2,
	//			Role:    string(vxdr.RolePriceFeeder),
	//		},
	//	}
	//	_ = veloTx.Build()
	//	_ = veloTx.Sign(kp1)
	//
	//	useCase := usecases.Init(nil)
	//	err := useCase.CreateWhiteList(context.Background(), veloTx)
	//
	//	assert.Equal(t, err.Error(), constants.ErrRoleNotFound)
	//})
	//
	//t.Run("Error - source account don't have regulator role", func(t *testing.T) {
	//	mockedWhiteListRepo, finish := newMockWhiteListRepo()
	//	defer finish()
	//
	//	filter := entities.WhiteListFilter{
	//		StellarPublicKey: &stellarPublicAddress,
	//		RoleCode:         &roleCode,
	//	}
	//
	//	mockedWhiteListRepo.EXPECT().FindOneWhitelist(filter).Return(nil, nil)
	//
	//	veloTx := &vtxnbuild.VeloTx{
	//		SourceAccount: &txnbuild.SimpleAccount{
	//			AccountID: publicKey1,
	//		},
	//		VeloOp: &vtxnbuild.WhiteList{
	//			Address: publicKey2,
	//			Role:    string(vxdr.RolePriceFeeder),
	//		},
	//	}
	//	_ = veloTx.Build()
	//	_ = veloTx.Sign(kp1)
	//
	//	useCase := usecases.Init(nil)
	//	err := useCase.CreateWhiteList(context.Background(), veloTx)
	//
	//	assert.Equal(t, err.Error(), fmt.Sprintf(constants.ErrFormatSignerNotHavePermission, constants.VeloOpWhiteList))
	//})
	//
	//t.Run("Error - send whitelist to save but fill invalid role", func(t *testing.T) {
	//	mockedWhiteListRepo, finish := newMockWhiteListRepo()
	//	defer finish()
	//
	//	findWhiteListEntity := entities.WhiteList{
	//		ID:               "e13d778c-d2c8-452b-8ead-368d43447fcd",
	//		StellarPublicKey: publicKey1,
	//		RoleCode:         string(vxdr.RoleRegulator),
	//	}
	//
	//	filter := entities.WhiteListFilter{
	//		StellarPublicKey: &stellarPublicAddress,
	//		RoleCode:         &roleCode,
	//	}
	//
	//	mockedWhiteListRepo.EXPECT().FindOneWhitelist(filter).Return(&findWhiteListEntity, nil)
	//	mockedWhiteListRepo.EXPECT().FindOneRole(string(vxdr.RolePriceFeeder)).Return(nil, errors.New(constants.ErrToGetDataFromDatabase))
	//
	//	veloTx := &vtxnbuild.VeloTx{
	//		SourceAccount: &txnbuild.SimpleAccount{
	//			AccountID: publicKey1,
	//		},
	//		VeloOp: &vtxnbuild.WhiteList{
	//			Address: publicKey2,
	//			Role:    string(vxdr.RolePriceFeeder),
	//		},
	//	}
	//	_ = veloTx.Build()
	//	_ = veloTx.Sign(kp1)
	//
	//	useCase := usecases.Init(nil)
	//	err := useCase.CreateWhiteList(context.Background(), veloTx)
	//
	//	assert.Equal(t, err.Error(), constants.ErrToGetDataFromDatabase)
	//})
	//
	//t.Run("Error - can't save whitelist table", func(t *testing.T) {
	//	mockedWhiteListRepo, finish := newMockWhiteListRepo()
	//	defer finish()
	//
	//	findWhiteListEntity := entities.WhiteList{
	//		ID:               "e13d778c-d2c8-452b-8ead-368d43447fcd",
	//		StellarPublicKey: publicKey1,
	//		RoleCode:         string(vxdr.RoleRegulator),
	//	}
	//
	//	roleEntity := entities.Role{
	//		ID:   1,
	//		Name: "Price feeder",
	//		Code: "PRICE_FEEDER",
	//	}
	//
	//	createWhitelistEntity := entities.WhiteList{
	//		StellarPublicKey: publicKey2,
	//		RoleCode:         string(vxdr.RolePriceFeeder),
	//	}
	//
	//	filter := entities.WhiteListFilter{
	//		StellarPublicKey: &stellarPublicAddress,
	//		RoleCode:         &roleCode,
	//	}
	//
	//	mockedWhiteListRepo.EXPECT().FindOneWhitelist(filter).Return(&findWhiteListEntity, nil)
	//	mockedWhiteListRepo.EXPECT().FindOneRole(string(vxdr.RolePriceFeeder)).Return(&roleEntity, nil)
	//	mockedWhiteListRepo.EXPECT().CreateWhitelist(&createWhitelistEntity).Return(nil, errors.New(constants.ErrToSaveDatabase))
	//
	//	veloTx := &vtxnbuild.VeloTx{
	//		SourceAccount: &txnbuild.SimpleAccount{
	//			AccountID: publicKey1,
	//		},
	//		VeloOp: &vtxnbuild.WhiteList{
	//			Address: publicKey2,
	//			Role:    string(vxdr.RolePriceFeeder),
	//		},
	//	}
	//	_ = veloTx.Build()
	//	_ = veloTx.Sign(kp1)
	//
	//	useCase := usecases.Init(nil)
	//	err := useCase.CreateWhiteList(context.Background(), veloTx)
	//
	//	assert.Equal(t, err.Error(), constants.ErrToSaveDatabase)
	//})
	//
	//t.Run("Error - can't save whitelist table, cause: already exits", func(t *testing.T) {
	//	mockedWhiteListRepo, finish := newMockWhiteListRepo()
	//	defer finish()
	//
	//	findWhiteListEntity := entities.WhiteList{
	//		ID:               "e13d778c-d2c8-452b-8ead-368d43447fcd",
	//		StellarPublicKey: publicKey1,
	//		RoleCode:         string(vxdr.RoleRegulator),
	//	}
	//
	//	roleEntity := entities.Role{
	//		ID:   1,
	//		Name: "Price feeder",
	//		Code: "PRICE_FEEDER",
	//	}
	//
	//	createWhitelistEntity := entities.WhiteList{
	//		StellarPublicKey: publicKey2,
	//		RoleCode:         string(vxdr.RolePriceFeeder),
	//	}
	//
	//	filter := entities.WhiteListFilter{
	//		StellarPublicKey: &stellarPublicAddress,
	//		RoleCode:         &roleCode,
	//	}
	//
	//	mockedWhiteListRepo.EXPECT().FindOneWhitelist(filter).Return(&findWhiteListEntity, nil)
	//	mockedWhiteListRepo.EXPECT().FindOneRole(string(vxdr.RolePriceFeeder)).Return(&roleEntity, nil)
	//	mockedWhiteListRepo.EXPECT().CreateWhitelist(&createWhitelistEntity).Return(nil, errors.New("duplicate key value violates unique constraint"))
	//
	//	veloTx := &vtxnbuild.VeloTx{
	//		SourceAccount: &txnbuild.SimpleAccount{
	//			AccountID: publicKey1,
	//		},
	//		VeloOp: &vtxnbuild.WhiteList{
	//			Address: publicKey2,
	//			Role:    string(vxdr.RolePriceFeeder),
	//		},
	//	}
	//	_ = veloTx.Build()
	//	_ = veloTx.Sign(kp1)
	//
	//	useCase := usecases.Init(nil)
	//	err := useCase.CreateWhiteList(context.Background(), veloTx)
	//
	//	assert.Contains(t, err.Error(), fmt.Sprintf(constants.ErrWhiteListAlreadyWhiteListed, publicKey1, vxdr.RoleMap[vxdr.RolePriceFeeder]))
	//})

}

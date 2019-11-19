package usecases_test

import (
	"github.com/golang/mock/gomock"
	"github.com/velo-protocol/DRSv1/node/app/layers/mocks"
	"github.com/velo-protocol/DRSv1/node/app/layers/usecases"
	"github.com/velo-protocol/DRSv1/node/app/testhelpers"
	"testing"
)

type helper struct {
	useCase         usecases.UseCase
	mockSubUseCase  *mocks.MockSubUseCase
	mockStellarRepo *mocks.MockStellarRepo
	mockController  *gomock.Controller
}

var (
	publicKey1 = testhelpers.PublicKey1
	secretKey1 = testhelpers.SecretKey1

	publicKey2 = testhelpers.PublicKey2
	secretKey2 = testhelpers.SecretKey2

	publicKey3 = testhelpers.PublicKey3
	secretKey3 = testhelpers.SecretKey3

	publicKey4 = testhelpers.PublicKey4
	secretKey4 = testhelpers.SecretKey4

	publicKey5 = testhelpers.PublicKey5
	secretKey5 = testhelpers.SecretKey5

	drsPublicKey = testhelpers.DrsPublicKey
	drsSecretKey = testhelpers.DrsSecretKey

	kp1 = testhelpers.Kp1
	kp2 = testhelpers.Kp2
	kp3 = testhelpers.Kp3

	drsAccountDataEnity = testhelpers.DrsAccountDataEntity
)

func initTest(t *testing.T) helper {
	testhelpers.InitEnv()

	mockCtrl := gomock.NewController(t)
	mockStellarRepo := mocks.NewMockStellarRepo(mockCtrl)
	mockSubUseCase := mocks.NewMockSubUseCase(mockCtrl)

	return helper{
		useCase:         usecases.Init(mockStellarRepo, mockSubUseCase),
		mockSubUseCase:  mockSubUseCase,
		mockStellarRepo: mockStellarRepo,
		mockController:  mockCtrl,
	}
}

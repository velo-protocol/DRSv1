package usecases_test

import (
	"github.com/golang/mock/gomock"
	"gitlab.com/velo-labs/cen/node/app/layers/mocks"
	"gitlab.com/velo-labs/cen/node/app/layers/usecases"
	"gitlab.com/velo-labs/cen/node/app/testhelpers"
	"testing"
)

type helper struct {
	useCase         usecases.UseCase
	mockStellarRepo *mocks.MockStellarRepo
	mockController  *gomock.Controller
}

var (
	publicKey1 = testhelpers.PublicKey1
	secretKey1 = testhelpers.SecretKey1

	publicKey2 = testhelpers.PublicKey2
	secretKey2 = testhelpers.SecretKey2

	publicKey3 = testhelpers.PublicKey3
	secretKey3 = testhelpers.PublicKey3

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

	return helper{
		useCase:         usecases.Init(mockStellarRepo),
		mockStellarRepo: mockStellarRepo,
		mockController:  mockCtrl,
	}
}

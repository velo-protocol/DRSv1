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

	drsPublicKey = "GCQCXIDTFMIL4VOAXWUQNRAMC46TTJDHZ3DDJVD32ND7B4OKANIUKB5N"
	drsSecretKey = "SDE374OE44ZU73KAUFYPNMQEUGCDIJLTIIUZ4W2MKWBPPAK36ID26ECU"

	kp1 = testhelpers.Kp1
	kp2 = testhelpers.Kp2
	kp3 = testhelpers.Kp3
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

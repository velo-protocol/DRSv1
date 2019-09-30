package usecases_test

import (
	"github.com/golang/mock/gomock"
	"gitlab.com/velo-labs/cen/node/app/layers/mocks"
	"gitlab.com/velo-labs/cen/node/app/layers/usecases"
	"gitlab.com/velo-labs/cen/node/app/testhelpers"
	"testing"
)

type testHelper struct {
	MockStellarRepo *mocks.MockStellarRepo
}

var (
	publicKey1 = testhelpers.PublicKey1
	secretKey1 = testhelpers.SecretKey1

	publicKey2 = testhelpers.PublicKey2
	secretKey2 = testhelpers.SecretKey2

	drsPublicKey = "GCQCXIDTFMIL4VOAXWUQNRAMC46TTJDHZ3DDJVD32ND7B4OKANIUKB5N"
	drsSecretKey = "SDE374OE44ZU73KAUFYPNMQEUGCDIJLTIIUZ4W2MKWBPPAK36ID26ECU"
)

func initUseCaseTest(t *testing.T) (usecases.UseCase, *testHelper, *gomock.Controller) {
	testhelpers.InitEnv()

	testHelper := new(testHelper)

	mockCtrl := gomock.NewController(t)

	testHelper.MockStellarRepo = mocks.NewMockStellarRepo(mockCtrl)

	useCase := usecases.Init(testHelper.MockStellarRepo)

	return useCase, testHelper, mockCtrl
}

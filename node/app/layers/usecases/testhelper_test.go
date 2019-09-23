package usecases_test

import (
	"github.com/golang/mock/gomock"
	"gitlab.com/velo-labs/cen/node/app/layers/mocks"
	"gitlab.com/velo-labs/cen/node/app/layers/usecases"
	"gitlab.com/velo-labs/cen/node/app/testhelpers"
	"testing"
)

type testHelper struct {
	MockWhiteListRepo *mocks.MockWhiteListRepo
	MockStellarRepo   *mocks.MockStellarRepo
}

var (
	publicKey1 = "GDU5DE2ZKAZ4BPIQ34ZXSJAIVGFF75SPYTCEHWR7PFNBG5SNDNTMGZFB"
	secretKey1 = "SCOPQ6WGCOO6SDUXOF6TUEOW2EE3DYT335JEBTTRQ4OGIBWZKBBPBSPY"

	publicKey2 = "GCA4XQDOYKPA57ZSWOAJMR2ACWK3MHGUYRX5NYEVSPHJTTFRVZCWZOUU"
	secretKey2 = "SCHHUC4WW7DST4TTNTOZ744F546DNKFXGMRFQPPR2MLZGFRHNGW3SXEI"

	drsPublicKey = "GCQCXIDTFMIL4VOAXWUQNRAMC46TTJDHZ3DDJVD32ND7B4OKANIUKB5N"
	drsSecretKey = "SDE374OE44ZU73KAUFYPNMQEUGCDIJLTIIUZ4W2MKWBPPAK36ID26ECU"
)

func initUseCaseTest(t *testing.T) (usecases.UseCase, *testHelper, *gomock.Controller) {
	testhelpers.InitEnv()

	testHelper := new(testHelper)

	mockCtrl := gomock.NewController(t)

	testHelper.MockWhiteListRepo = mocks.NewMockWhiteListRepo(mockCtrl)
	testHelper.MockStellarRepo = mocks.NewMockStellarRepo(mockCtrl)

	useCase := usecases.Init(testHelper.MockStellarRepo, testHelper.MockWhiteListRepo)

	return useCase, testHelper, mockCtrl
}

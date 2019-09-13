package usecases_test

import (
	"github.com/stellar/go/protocols/horizon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	_nodeRepositoryMocked "gitlab.com/velo-labs/cen/node/app/modules/node/mocks"
	"gitlab.com/velo-labs/cen/node/app/modules/node/usecases"
	_stellarRepositoryMocked "gitlab.com/velo-labs/cen/node/app/modules/stellar/mocks"
	_operationServiceMocked "gitlab.com/velo-labs/cen/node/app/services/operation/mocks"
	test_helpers "gitlab.com/velo-labs/cen/node/app/test-helpers"
	"testing"
)

func TestUsecase_Setup(t *testing.T) {
	test_helpers.InitEnv()

	t.Run("happy", func(t *testing.T) {
		mockedDrsOps := new(_operationServiceMocked.Interface)
		mockedStellarRepository := new(_stellarRepositoryMocked.Repository)
		mockedNodeRepository := new(_nodeRepositoryMocked.Repository)

		mockedDrsOps.On(
			"Setup",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(test_helpers.GetIssuerCreationTxB64(), test_helpers.GetRandStellarAccount().Address(), test_helpers.GetRandStellarAccount().Address(), nil)

		mockedStellarRepository.On("SubmitTransaction", mock.AnythingOfType("string")).
			Return(&horizon.TransactionSuccess{
				Hash: "fake-hash",
			}, nil)

		mockedNodeRepository.On("SaveCredit", mock.AnythingOfType("entities.Credit")).Return(nil)

		uc := usecases.NewNodeUseCase(mockedDrsOps, mockedNodeRepository, mockedStellarRepository)

		setupResult, err := uc.Setup(test_helpers.GetSetupTxB64(), "1", "USD", "vUSD")

		assert.NoError(t, err)
		assert.NotEmpty(t, setupResult)
		assert.Equal(t, "fake-hash", setupResult.IssuerCreationTxHash)
	})
}

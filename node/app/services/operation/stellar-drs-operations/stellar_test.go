package stellar_drsops_test

import (
	"github.com/stellar/go/xdr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	_stellarRepositoryMocks "gitlab.com/velo-labs/cen/app/modules/stellar/mocks"
	drsops "gitlab.com/velo-labs/cen/app/services/operation/stellar-drs-operations"
	test_helpers "gitlab.com/velo-labs/cen/app/test-helpers"
	"testing"
)

func TestOps_Mint(t *testing.T) {
	test_helpers.InitEnv()
	mockedStellarAccount := test_helpers.GetStellarAccount()

	t.Run("happy", func(t *testing.T) {
		mockedStellarRepository := new(_stellarRepositoryMocks.Repository)

		mockedStellarRepository.On("LoadAccount", mock.AnythingOfType("string")).
			Return(&mockedStellarAccount, nil)

		sv := drsops.NewDrsOps(mockedStellarRepository)
		setupTxB64, _, _, err := sv.Setup("1", "USD", "vUSD", "GAA456EKXSBHCCEGLA5DABYFOS3CQJJ566RYFMOKKJW4TV5FZ7NVL3AD")
		assert.NoError(t, err)

		var txe xdr.TransactionEnvelope

		err = xdr.SafeUnmarshalBase64(setupTxB64, &txe)

		assert.NoError(t, err)
		assert.NotEmpty(t, setupTxB64)
		assert.Equal(t, 11, len(txe.Tx.Operations))
	})
}

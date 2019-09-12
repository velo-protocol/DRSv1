package repositories_test

import (
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/testutil"
	"gitlab.com/velo-labs/cen/node/app/modules/node/repositories"
	test_helpers "gitlab.com/velo-labs/cen/node/app/test-helpers"
	"testing"
)

func TestRepository_SaveCredit(t *testing.T) {
	gomega.RegisterTestingT(t)

	t.Run("happy", func(t *testing.T) {
		mockedLevelDB, err := leveldb.Open(testutil.NewStorage(), &opt.Options{DisableLargeBatchTransaction: true})
		assert.NoError(t, err)

		nodeRepository := repositories.NewNodeRepository(mockedLevelDB)

		err = nodeRepository.SaveCredit(test_helpers.GetCreditEntity())

		assert.NoError(t, err)
	})
}

package whitelist_test

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"testing"
)

func TestRepo_CommitTx(t *testing.T) {

	t.Run("Happy", func(t *testing.T) {
		sqlMock, repo := initRepoTest()

		sqlMock.ExpectBegin()
		sqlMock.ExpectCommit()

		dbTx := repo.BeginTx()
		assert.NotNil(t, dbTx)
		assert.IsType(t, &gorm.DB{}, dbTx)

		err := repo.CommitTx(dbTx)

		assert.NoError(t, err)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("Error", func(t *testing.T) {
		sqlMock, repo := initRepoTest()

		sqlMock.ExpectBegin()
		sqlMock.ExpectCommit().
			WillReturnError(errors.New(constants.ErrToCommitTransaction))

		dbTx := repo.BeginTx()
		assert.NotNil(t, dbTx)
		assert.IsType(t, &gorm.DB{}, dbTx)

		err := repo.CommitTx(dbTx)

		assert.Error(t, err)
		assert.EqualError(t, errors.New(constants.ErrToCommitTransaction), err.Error())
		assert.NoError(t, sqlMock.ExpectationsWereMet())

	})
}

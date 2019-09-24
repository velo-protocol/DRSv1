package whitelist_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"regexp"
	"testing"
)

var expectedSQLCommand = fmt.Sprintf(`INSERT INTO "%s"`, constants.WhiteListTable)

func TestRepo_CreateWhitelistTx(t *testing.T) {

	t.Run("Happy", func(t *testing.T) {
		sqlMock, repo := initRepoTest()

		expectedResult := &entities.WhiteList{
			ID:               "8008e30d-33a9-4e6f-a544-283b52788f2a",
			StellarPublicKey: "GA54PTUVTYZNT4RO5BJDHW5EX5YCMIZQOW62ZZCQYUDUXBAINS3L6PU3",
			RoleCode:         string(vxdr.RoleTrustedPartner),
		}

		sqlMock.ExpectBegin()
		sqlMock.ExpectExec(regexp.QuoteMeta(expectedSQLCommand)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		dbTx := repo.BeginTx()
		result, err := repo.CreateWhitelistTx(dbTx, expectedResult)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult.ID, result.ID)
		assert.Equal(t, expectedResult.StellarPublicKey, result.StellarPublicKey)
		assert.Equal(t, expectedResult.RoleCode, result.RoleCode)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("Error - cannot save to database", func(t *testing.T) {
		sqlMock, repo := initRepoTest()

		sqlMock.ExpectBegin()
		sqlMock.ExpectExec(regexp.QuoteMeta(expectedSQLCommand)).
			WillReturnError(errors.New(constants.ErrToSaveDatabase))

		dbTx := repo.BeginTx()
		result, err := repo.CreateWhitelistTx(dbTx, &entities.WhiteList{})

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, errors.New(constants.ErrToSaveDatabase), err.Error())
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

}

func TestRepo_CreateWhitelist(t *testing.T) {

	t.Run("Happy", func(t *testing.T) {
		sqlMock, repo := initRepoTest()

		expectedResult := &entities.WhiteList{
			ID:               "8008e30d-33a9-4e6f-a544-283b52788f2a",
			StellarPublicKey: "GA54PTUVTYZNT4RO5BJDHW5EX5YCMIZQOW62ZZCQYUDUXBAINS3L6PU3",
			RoleCode:         string(vxdr.RoleTrustedPartner),
		}

		sqlMock.ExpectBegin()
		sqlMock.ExpectExec(regexp.QuoteMeta(expectedSQLCommand)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		sqlMock.ExpectCommit()

		result, err := repo.CreateWhitelist(expectedResult)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult.ID, result.ID)
		assert.Equal(t, expectedResult.StellarPublicKey, result.StellarPublicKey)
		assert.Equal(t, expectedResult.RoleCode, result.RoleCode)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("Error - cannot save to database", func(t *testing.T) {
		sqlMock, repo := initRepoTest()

		sqlMock.ExpectBegin()
		sqlMock.ExpectExec(regexp.QuoteMeta(expectedSQLCommand)).
			WillReturnError(errors.New(constants.ErrToSaveDatabase))

		result, err := repo.CreateWhitelist(&entities.WhiteList{})

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, errors.New(constants.ErrToSaveDatabase), err.Error())
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

}

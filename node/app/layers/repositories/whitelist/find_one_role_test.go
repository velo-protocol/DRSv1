package whitelist_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	vxdr "gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"regexp"
	"testing"
)

func TestRepo_FindOneRole(t *testing.T) {

	expectedSQLCommand := fmt.Sprintf(`SELECT * FROM "%s"`, constants.RoleTable)

	t.Run("Happy", func(t *testing.T) {
		sqlMock, repo := initRepoTest()

		expectedResult := entities.Role{
			ID:   1,
			Name: "KYC checker",
			Code: string(vxdr.RoleRegulator),
		}

		rows := sqlmock.NewRows([]string{"id", "name", "code"}).
			AddRow(expectedResult.ID, expectedResult.Name, expectedResult.Code)

		sqlMock.ExpectQuery(regexp.QuoteMeta(expectedSQLCommand)).
			WithArgs(string(vxdr.RoleTrustedPartner)).
			WillReturnRows(rows)

		result, err := repo.FindOneRole(string(vxdr.RoleTrustedPartner))

		assert.Nil(t, err)
		assert.Equal(t, expectedResult.ID, result.ID)
		assert.Equal(t, expectedResult.Name, result.Name)
		assert.Equal(t, expectedResult.Code, result.Code)

		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("Happy - Record not found ", func(t *testing.T) {
		sqlMock, repo := initRepoTest()

		sqlMock.ExpectQuery(regexp.QuoteMeta(expectedSQLCommand)).
			WithArgs(string(vxdr.RoleTrustedPartner)).
			WillReturnRows(&sqlmock.Rows{})

		result, err := repo.FindOneRole(string(vxdr.RoleTrustedPartner))

		assert.Nil(t, result)
		assert.Nil(t, err)

		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("Error ", func(t *testing.T) {
		sqlMock, repo := initRepoTest()

		sqlMock.ExpectQuery(regexp.QuoteMeta(expectedSQLCommand)).
			WithArgs(string(vxdr.RoleTrustedPartner)).
			WillReturnError(errors.New(constants.ErrToGetDataFromDatabase))

		result, err := repo.FindOneRole(string(vxdr.RoleTrustedPartner))

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, errors.New(constants.ErrToGetDataFromDatabase), err.Error())

		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})
}

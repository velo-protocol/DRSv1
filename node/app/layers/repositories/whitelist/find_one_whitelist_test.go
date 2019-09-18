package whitelist_test

import (
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"regexp"
	"testing"
)

func TestRepo_FindOneWhitelist(t *testing.T) {

	expectedSqlCommand := fmt.Sprintf(`SELECT * FROM "%s"`, constants.WhiteListTable)

	t.Run("Happy", func(t *testing.T) {
		sqlMock, repo := initRepoTest()

		expectedResult := entities.WhiteList{
			ID:               "08f1655c-f650-4b4f-80e5-1395e64c0809",
			StellarPublicKey: "GDP3LU4CM3L2PQRNAEMKZPJ5BEE7I3XYO4VUQQISFGBTT6URPW4VFCIK",
			RoleCode:         string(vxdr.RoleRegulator),
		}

		filter := entities.WhiteListFilter{
			StellarPublicKey: pointer.ToString("GDP3LU4CM3L2PQRNAEMKZPJ5BEE7I3XYO4VUQQISFGBTT6URPW4VFCIK"),
		}

		rows := sqlmock.NewRows([]string{"id", "stellar_public_key", "role_code"}).
			AddRow(expectedResult.ID, expectedResult.StellarPublicKey, expectedResult.RoleCode)

		sqlMock.ExpectQuery(regexp.QuoteMeta(expectedSqlCommand)).
			WithArgs(filter.StellarPublicKey).
			WillReturnRows(rows)

		result, err := repo.FindOneWhitelist(filter)

		assert.Nil(t, err)
		assert.Equal(t, expectedResult.ID, result.ID)
		assert.Equal(t, expectedResult.StellarPublicKey, result.StellarPublicKey)
		assert.Equal(t, expectedResult.RoleCode, result.RoleCode)

		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("Happy - Record not found ", func(t *testing.T) {
		sqlMock, repo := initRepoTest()

		filter := entities.WhiteListFilter{
			StellarPublicKey: pointer.ToString("GDP3LU4CM3L2PQRNAEMKZPJ5BEE7I3XYO4VUQQISFGBTT6URPW4VFCIK"),
		}

		sqlMock.ExpectQuery(regexp.QuoteMeta(expectedSqlCommand)).
			WithArgs(filter.StellarPublicKey).
			WillReturnRows(&sqlmock.Rows{})

		result, err := repo.FindOneWhitelist(filter)

		assert.Nil(t, result)
		assert.Nil(t, err)

		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("Error ", func(t *testing.T) {
		sqlMock, repo := initRepoTest()

		sqlMock.ExpectQuery(regexp.QuoteMeta(expectedSqlCommand)).
			WillReturnError(constants.ErrToGetDataFromDatabase)

		result, err := repo.FindOneWhitelist(entities.WhiteListFilter{})

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, constants.ErrToGetDataFromDatabase, err.Error())

		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})
}

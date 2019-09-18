package whitelist_test

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepo_BeginTx(t *testing.T) {

	t.Run("Happy", func(t *testing.T) {
		sqlMock, repo := initRepoTest()

		sqlMock.ExpectBegin()

		dbTx := repo.BeginTx()

		assert.NotNil(t, dbTx)
		assert.IsType(t, &gorm.DB{}, dbTx)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

}

package price_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"regexp"
	"testing"
)

func TestRepo_CreatePriceEntryTx(t *testing.T) {
	expectedSQLCommand := fmt.Sprintf(`INSERT INTO "%s"`, constants.PriceEntryTable)

	t.Run("Happy", func(t *testing.T) {
		sqlMock, repo, mockDb := initRepoTest()

		priceEntryEntity := &entities.CreatePriceEntry{
			FeederPublicKey:             "GA54PTUVTYZNT4RO5BJDHW5EX5YCMIZQOW62ZZCQYUDUXBAINS3L6PU3",
			Asset:                       "VELO",
			PriceInCurrencyPerAssetUnit: decimal.NewFromFloat(2),
			Currency:                    "THB",
		}

		sqlMock.ExpectBegin()
		sqlMock.ExpectExec(regexp.QuoteMeta(expectedSQLCommand)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		sqlMock.ExpectCommit()

		result, err := repo.CreatePriceEntryTx(mockDb.Begin(), priceEntryEntity)

		assert.NoError(t, err)
		assert.Equal(t, priceEntryEntity.FeederPublicKey, result.FeederPublicKey)
		assert.Equal(t, priceEntryEntity.Asset, result.Asset)
		assert.Equal(t, priceEntryEntity.PriceInCurrencyPerAssetUnit, result.PriceInCurrencyPerAssetUnit)
		assert.Equal(t, priceEntryEntity.Currency, result.Currency)
		//assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("Error", func(t *testing.T) {
		sqlMock, repo, mockDb := initRepoTest()

		sqlMock.ExpectBegin()
		sqlMock.ExpectExec(regexp.QuoteMeta(expectedSQLCommand)).
			WillReturnError(errors.New(constants.ErrToSaveDatabase))

		result, err := repo.CreatePriceEntryTx(mockDb.Begin(), &entities.CreatePriceEntry{})

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, errors.New(constants.ErrToSaveDatabase), err.Error())
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})
}

func TestRepo_CreatePriceEntry(t *testing.T) {
	expectedSQLCommand := fmt.Sprintf(`INSERT INTO "%s"`, constants.PriceEntryTable)

	t.Run("Happy", func(t *testing.T) {
		sqlMock, repo, _ := initRepoTest()

		priceEntryEntity := &entities.CreatePriceEntry{
			FeederPublicKey:             "GA54PTUVTYZNT4RO5BJDHW5EX5YCMIZQOW62ZZCQYUDUXBAINS3L6PU3",
			Asset:                       "VELO",
			PriceInCurrencyPerAssetUnit: decimal.NewFromFloat(2),
			Currency:                    "THB",
		}

		sqlMock.ExpectBegin()
		sqlMock.ExpectExec(regexp.QuoteMeta(expectedSQLCommand)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		sqlMock.ExpectCommit()

		result, err := repo.CreatePriceEntry(priceEntryEntity)

		assert.NoError(t, err)
		assert.Equal(t, priceEntryEntity.FeederPublicKey, result.FeederPublicKey)
		assert.Equal(t, priceEntryEntity.Asset, result.Asset)
		assert.Equal(t, priceEntryEntity.PriceInCurrencyPerAssetUnit, result.PriceInCurrencyPerAssetUnit)
		assert.Equal(t, priceEntryEntity.Currency, result.Currency)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("Error", func(t *testing.T) {
		sqlMock, repo, _ := initRepoTest()

		sqlMock.ExpectBegin()
		sqlMock.ExpectExec(regexp.QuoteMeta(expectedSQLCommand)).
			WillReturnError(errors.New(constants.ErrToSaveDatabase))

		result, err := repo.CreatePriceEntry(&entities.CreatePriceEntry{})

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, errors.New(constants.ErrToSaveDatabase), err.Error())
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})
}

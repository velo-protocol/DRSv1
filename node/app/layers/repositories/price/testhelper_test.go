package price_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"gitlab.com/velo-labs/cen/node/app/layers/repositories/price"
)

func initRepoTest() (sqlmock.Sqlmock, price.Repo, *gorm.DB) {
	db, sqlMock, _ := sqlmock.New()

	gormMock, _ := gorm.Open("", db)

	_ = gormMock.LogMode(true)

	repo := price.InitRepo(gormMock)

	return sqlMock, repo, gormMock
}

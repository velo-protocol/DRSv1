package whitelist_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"gitlab.com/velo-labs/cen/node/app/layers/repositories/whitelist"
)

func initRepoTest() (sqlmock.Sqlmock, whitelist.Repo) {
	db, sqlMock, _ := sqlmock.New()

	gormMock, _ := gorm.Open("", db)

	_ = gormMock.LogMode(true)

	repo := whitelist.InitRepo(gormMock)

	return sqlMock, repo
}

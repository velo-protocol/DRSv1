package logic

import (
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/repositories/database"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/repositories/stellar"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/repositories/velo"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/config"
)

type logic struct {
	DB        database.Repository
	Stellar   stellar.Repository
	Velo      velo.Repository
	AppConfig config.Configuration
}

func NewLogic(db database.Repository, st stellar.Repository, velo velo.Repository, config config.Configuration) Logic {
	return &logic{
		DB:        db,
		Stellar:   st,
		Velo:      velo,
		AppConfig: config,
	}
}

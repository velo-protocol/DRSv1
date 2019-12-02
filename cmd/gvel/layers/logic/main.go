package logic

import (
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/repositories/database"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/repositories/stellar"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/repositories/velo"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/config"
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

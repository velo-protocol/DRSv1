package logic

import (
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/repositories/database"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/repositories/friendbot"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/repositories/velo"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/config"
)

type logic struct {
	DB        database.Repository
	FriendBot friendbot.Repository
	Velo      velo.Repository
	AppConfig config.Configuration
}

func NewLogic(db database.Repository, fb friendbot.Repository, velo velo.Repository, config config.Configuration) Logic {
	return &logic{
		DB:        db,
		FriendBot: fb,
		Velo:      velo,
		AppConfig: config,
	}
}

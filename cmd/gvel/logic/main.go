package logic

import (
	"gitlab.com/velo-labs/cen/cmd/gvel/db"
	"gitlab.com/velo-labs/cen/cmd/gvel/friendbot"
)

type logic struct {
	DB        db.DB
	Friendbot friendbot.Repository
}

func NewLogic(db db.DB, fb friendbot.Repository) Logic {
	return &logic{
		DB:        db,
		Friendbot: fb,
	}
}

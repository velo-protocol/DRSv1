package logic

import (
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/repositories/database"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/repositories/friendbot"
)

type logic struct {
	DB        database.Repository
	FriendBot friendbot.Repository
}

func NewLogic(db database.Repository, fb friendbot.Repository) Logic {
	return &logic{
		DB:        db,
		FriendBot: fb,
	}
}

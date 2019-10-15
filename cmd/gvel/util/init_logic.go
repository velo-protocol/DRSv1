package util

import (
	"github.com/pkg/errors"
	"gitlab.com/velo-labs/cen/cmd/gvel/db/leveldb"
	_friendbot "gitlab.com/velo-labs/cen/cmd/gvel/friendbot"
	"gitlab.com/velo-labs/cen/cmd/gvel/logic"
)

func InitLogicWithoutDB() logic.Logic {
	return logic.NewLogic(nil, nil)
}

func InitLogic(dbPath string, friendbotURL string) (logic.Logic, error) {
	leveldbConn, err := leveldb.NewLevelDB(dbPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init logic")
	}

	friendbot := _friendbot.NewFriendbot(friendbotURL)

	lo := logic.NewLogic(leveldbConn, friendbot)

	return lo, err
}

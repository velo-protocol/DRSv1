package logic

import (
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/config"
)

func (lo *logic) Init(configFilePath string) error {
	err := config.InitConfigFile(configFilePath)
	if err != nil {
		return err
	}

	_, err = leveldb.OpenFile(viper.GetString("accountDbPath"), nil)
	if err != nil {
		return err
	}

	//viper.AutomaticEnv()

	return nil
}

package logic

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"os"
	"path"
)

func (lo *logic) Init(configFilePath string) error {
	err := setConfigFile(configFilePath)
	if err != nil {
		return err
	}

	_, err = leveldb.OpenFile(path.Join(configFilePath, "/account"), nil)
	if err != nil {
		return err
	}

	viper.AutomaticEnv()

	return nil
}

func setConfigFile(configPath string) error {
	viper.SetConfigFile(path.Join(configPath, "/config.json"))
	err := viper.ReadInConfig()
	if err == nil {
		log.Println("config file is existed")
		return nil
	}

	err = os.MkdirAll(configPath, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "failed to create a config folder")
	}

	err = os.MkdirAll(path.Join(configPath, "/account"), os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "failed to create a db folder")
	}

	viper.SetDefault("configPath", path.Join(configPath, "/config.json"))
	viper.SetDefault("accountPath", path.Join(configPath, "/account"))

	err = viper.WriteConfig()
	if err != nil {
		return errors.Wrap(err, "failed to write a config to the disk")
	}

	return nil
}

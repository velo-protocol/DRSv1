package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gitlab.com/velo-labs/cen/cmd/gvel/constants"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
	"os"
	"path"
)

func InitConfigFile(configFilePath string) error {
	return setupConfigFile(configFilePath)
}

func Exists() bool {
	return viper.GetBool("initialized")
}

func Load() {
	_ = load(constants.DefaultConfigFilePath)
}

func load(configPath string) error {
	viper.SetConfigType("json")
	viper.SetConfigFile(path.Join(configPath, "/config.json"))
	return viper.ReadInConfig()
}

func setupConfigFile(configPath string) error {
	_ = load(configPath)

	if Exists() {
		console.Logger.Error("config file found")
		return nil
	}

	err := os.MkdirAll(configPath, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "failed to create a config folder")
	}

	err = os.MkdirAll(path.Join(configPath, "/db/account"), os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "failed to create a db folder")
	}

	//viper.SetDefault("configPath", path.Join(configPath, "/config.json"))

	// Set default config
	viper.SetDefault("initialized", true) // a flag to check for config file existence
	viper.SetDefault("accountDbPath", path.Join(configPath, "/db/account"))
	viper.SetDefault("friendBotUrl", constants.DefaultFriendBotUrl)

	err = viper.WriteConfig()
	if err != nil {
		return errors.Wrap(err, "failed to write a config to the disk")
	}

	return nil
}

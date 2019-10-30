package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gitlab.com/velo-labs/cen/cmd/gvel/constants"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
	"os"
	"path"
)

type configuration struct {
	viper *viper.Viper
}

func NewConfiguration() *configuration {
	return &configuration{
		viper: viper.GetViper(),
	}
}

func (configuration *configuration) LoadDefault() {
	_ = configuration.Load(constants.DefaultConfigFilePath)
}

func (configuration *configuration) Load(configPath string) error {
	configuration.viper.SetConfigType("json")
	configuration.viper.SetConfigFile(path.Join(configPath, "/config.json"))
	return configuration.viper.ReadInConfig()
}

func (configuration *configuration) InitConfigFile(configFilePath string) error {
	_ = configuration.Load(configFilePath)

	if configuration.Exists() {
		console.Logger.Error("config file found")
		return nil
	}

	err := os.MkdirAll(configFilePath, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "failed to create a config folder")
	}

	err = os.MkdirAll(path.Join(configFilePath, "/db/account"), os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "failed to create a db folder")
	}

	// Set default config
	configuration.viper.SetDefault("initialized", true) // a flag to check for config file existence
	configuration.viper.SetDefault("accountDbPath", path.Join(configFilePath, "/db/account"))
	configuration.viper.SetDefault("defaultAccount", "")

	configuration.viper.SetDefault("friendBotUrl", constants.DefaultFriendBotUrl)
	configuration.viper.SetDefault("horizonUrl", constants.DefaultHorizonUrl)
	configuration.viper.SetDefault("veloNodeUrl", constants.DefaultVeloNodeUrl)
	configuration.viper.SetDefault("networkPassphrase", constants.DefaultNetworkPassphrase)

	err = configuration.viper.WriteConfig()
	if err != nil {
		return errors.Wrap(err, "failed to write a config to the disk")
	}

	return nil
}

func (configuration *configuration) Exists() bool {
	return configuration.viper.GetBool("initialized")
}

func (configuration *configuration) GetDefaultAccount() string {
	return configuration.viper.GetString("defaultAccount")
}

func (configuration *configuration) SetDefaultAccount(account string) error {
	configuration.viper.Set("defaultAccount", account)
	return configuration.viper.WriteConfig()
}

func (configuration *configuration) GetAccountDbPath() string {
	return configuration.viper.GetString("accountDbPath")
}

func (configuration *configuration) GetFriendBotUrl() string {
	return configuration.viper.GetString("friendBotUrl")
}

func (configuration *configuration) GetHorizonUrl() string {
	return configuration.viper.GetString("horizonBotUrl")
}

func (configuration *configuration) GetVeloNodeUrl() string {
	return configuration.viper.GetString("veloNodeUrl")
}

func (configuration *configuration) GetNetworkPassphrase() string {
	return configuration.viper.GetString("networkPassphrase")
}

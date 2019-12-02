package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/velo-protocol/DRSv1/cmd/gvel/constants"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/console"
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
	configuration.viper.SetDefault(constants.ConfInitialized, true) // a flag to check for config file existence
	configuration.viper.SetDefault(constants.ConfAccountDbPath, path.Join(configFilePath, "/db/account"))
	configuration.viper.SetDefault(constants.ConfDefaultAccount, "")

	configuration.viper.SetDefault(constants.ConfHorizonUrl, constants.DefaultHorizonUrl)
	configuration.viper.SetDefault(constants.ConfVeloNodeUrl, constants.DefaultVeloNodeUrl)
	configuration.viper.SetDefault(constants.ConfNetworkPassphrase, constants.DefaultNetworkPassphrase)
	configuration.viper.SetDefault(constants.ConfIsTestNet, true)

	err = configuration.viper.WriteConfig()
	if err != nil {
		return errors.Wrap(err, "failed to write a config to the disk")
	}

	return nil
}

func (configuration *configuration) Exists() bool {
	return configuration.viper.GetBool(constants.ConfInitialized)
}

func (configuration *configuration) GetDefaultAccount() string {
	return configuration.viper.GetString(constants.ConfDefaultAccount)
}

func (configuration *configuration) SetDefaultAccount(account string) error {
	configuration.viper.Set(constants.ConfDefaultAccount, account)
	return configuration.viper.WriteConfig()
}

func (configuration *configuration) GetAccountDbPath() string {
	return configuration.viper.GetString(constants.ConfAccountDbPath)
}

func (configuration *configuration) GetHorizonUrl() string {
	return configuration.viper.GetString(constants.ConfHorizonUrl)
}

func (configuration *configuration) GetVeloNodeUrl() string {
	return configuration.viper.GetString(constants.ConfVeloNodeUrl)
}

func (configuration *configuration) GetNetworkPassphrase() string {
	return configuration.viper.GetString(constants.ConfNetworkPassphrase)
}

func (configuration *configuration) GetIsTestNet() bool {
	return configuration.viper.GetBool(constants.ConfIsTestNet)
}

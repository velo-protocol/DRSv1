package init

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	_default "gitlab.com/velo-labs/cen/cmd/gvel/default"
	"log"
	"os"
	"path"
)

func NewInitCmd() *cobra.Command {
	cmd := cobra.Command{
		Use: "init",
		Short: "Use init command for initializing all configurations",
		Run: initRunner,
	}

	return &cmd
}

func initRunner(cmd *cobra.Command, args []string) {
	err := setConfigFile(_default.DefaultConfigFilePath)
	if err != nil {
		panic(err)
	}

	_, err = leveldb.OpenFile(path.Join(_default.DefaultConfigFilePath, "/db"), nil)
	if err != nil {
		panic(err)
	}

	viper.AutomaticEnv()

	log.Printf("gvel had been initialized\n")
	log.Printf("using config file at: %s\n", _default.DefaultConfigFilePath)
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

	err = os.MkdirAll(path.Join(configPath, "/db"), os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "failed to create a db folder")
	}

	viper.SetDefault("configPath", path.Join(configPath, "/config.json"))
	viper.SetDefault("databasePath", path.Join(configPath,"/db"))

	err = viper.WriteConfig()
	if err != nil {
		return errors.Wrap(err, "failed to write a config to the disk")
	}

	return nil
}

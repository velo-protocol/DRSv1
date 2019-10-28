package config

import (
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/viper"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
	"os"
)

type helper struct {
	config     *configuration
	loggerHook *test.Hook
	done       func()
}

func initTest() *helper {
	logger, hook := test.NewNullLogger()
	console.Logger = logger

	return &helper{
		config:     NewConfiguration(),
		loggerHook: hook,
		done: func() {
			hook.Reset()
			viper.Reset()
			_ = os.RemoveAll("./.gvel")
		},
	}

}

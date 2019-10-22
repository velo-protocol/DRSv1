package console

import "github.com/sirupsen/logrus"

var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()
	Logger.Formatter = &logrus.TextFormatter{
		ForceColors:               false,
		DisableColors:             false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          true,
		FullTimestamp:             false,
		TimestampFormat:           "",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    true,
		QuoteEmptyFields:          false,
		FieldMap:                  nil,
	}
}

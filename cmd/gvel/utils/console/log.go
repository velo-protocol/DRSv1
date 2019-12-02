package console

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger *logrus.Logger
var TableLogger *logrus.Logger

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

	TableLogger = logrus.New()
	TableLogger.Out = os.Stdout
	TableLogger.Formatter = &tableFormatter{}
}

type tableFormatter struct{}

func (tableFormatter *tableFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message), nil
}

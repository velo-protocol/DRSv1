package constant

import (
	"os"
	"path"
)

var (
	DefaultConfigFilePath     = path.Join(os.Getenv("HOME"), "/.gvel")
	DefaultGvelConfigPath     = path.Join(DefaultConfigFilePath, "/config.json")
	DefaultGevelAccountDbPath = path.Join(DefaultConfigFilePath, "/account")
)

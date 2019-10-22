package constants

import (
	"os"
	"path"
)

var (
	DefaultConfigFilePath    = path.Join(os.Getenv("HOME"), "/.gvel")
	DefaultGvelConfigPath    = path.Join(DefaultConfigFilePath, "/config.json")
	DefaultGvelDbAccountPath = path.Join(DefaultConfigFilePath, "/db/account")
)

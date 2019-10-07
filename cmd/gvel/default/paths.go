package _default

import (
	"os"
	"path"
)

var (
	DefaultConfigFilePath = path.Join(os.Getenv("HOME"), "/.velo")
)

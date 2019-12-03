package constants

import (
	"github.com/stellar/go/network"
	"os"
	"path"
)

var DefaultConfigFilePath = path.Join(os.Getenv("HOME"), "/.gvel")

const (
	DefaultHorizonUrl        = "https://horizon-testnet.stellar.org"
	DefaultVeloNodeUrl       = "testnet-drsv1-0.velo.org:8080"
	DefaultNetworkPassphrase = network.TestNetworkPassphrase
)

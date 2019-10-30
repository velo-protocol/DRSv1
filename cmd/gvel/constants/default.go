package constants

import (
	"github.com/stellar/go/network"
	"os"
	"path"
)

var DefaultConfigFilePath = path.Join(os.Getenv("HOME"), "/.gvel")

const (
	DefaultFriendBotUrl      = "https://horizon-testnet.stellar.org/friendbot?addr=%s"
	DefaultHorizonUrl        = "https://horizon-testnet.stellar.org"
	DefaultVeloNodeUrl       = "dev-velo-cen-node-01.velo.org:8080"
	DefaultNetworkPassphrase = network.TestNetworkPassphrase
)

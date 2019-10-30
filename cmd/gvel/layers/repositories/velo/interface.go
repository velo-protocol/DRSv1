package velo

import (
	"github.com/stellar/go/keypair"
	"gitlab.com/velo-labs/cen/libs/client"
)

type Repository interface {
	Client(keyPair *keypair.Full) vclient.ClientInterface
}

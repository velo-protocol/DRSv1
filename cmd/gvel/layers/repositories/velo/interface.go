package velo

import (
	"github.com/stellar/go/keypair"
	"github.com/velo-protocol/DRSv1/libs/client"
)

type Repository interface {
	Client(keyPair *keypair.Full) vclient.ClientInterface
}

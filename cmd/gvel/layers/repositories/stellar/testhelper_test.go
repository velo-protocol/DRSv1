package stellar_test

import (
	"github.com/stellar/go/clients/horizonclient"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/repositories/stellar"
)

var (
	FakeFriendBotUrl   = "https://fake-friendbot.stellar.org/friendbot?addr=%s"
	FakeHorizonUrl     = "https://fake-friendbot.stellar.org"
	FakeStellarAddress = "GCQ3BC5FLZ7T6OWOTIM3YLWEJ6FLPO525KTRO3JMPEXI3A4SRSJQG7KG"
)

type helper struct {
	repo                stellar.Repository
	mockedHorizonClient *horizonclient.MockClient
}

func initTest() *helper {
	h := helper{}
	h.mockedHorizonClient = &horizonclient.MockClient{}
	h.repo = stellar.NewStellarWithClientInterface(FakeHorizonUrl, h.mockedHorizonClient)
	return &h
}

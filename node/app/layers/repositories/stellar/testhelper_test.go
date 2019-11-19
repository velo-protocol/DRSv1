package stellar_test

import (
	"github.com/stellar/go/clients/horizonclient"
	"github.com/velo-protocol/DRSv1/node/app/layers/repositories/stellar"
	"github.com/velo-protocol/DRSv1/node/app/testhelpers"
)

type helper struct {
	repo                stellar.Repo
	mockedHorizonClient *horizonclient.MockClient
}

func initTest() helper {
	testhelpers.InitEnv()

	h := helper{}
	h.mockedHorizonClient = &horizonclient.MockClient{}
	h.repo = stellar.Init(h.mockedHorizonClient)
	return h
}

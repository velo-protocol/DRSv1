package stellar_test

import (
	"github.com/stellar/go/clients/horizonclient"
	"gitlab.com/velo-labs/cen/node/app/layers/repositories/stellar"
	"gitlab.com/velo-labs/cen/node/app/testhelpers"
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

package stellar

import (
	"github.com/stellar/go/clients/horizonclient"
)

type repo struct {
	HorizonClient horizonclient.ClientInterface
}

func Init(horizonClient horizonclient.ClientInterface) Repo {
	return &repo{
		HorizonClient: horizonClient,
	}
}

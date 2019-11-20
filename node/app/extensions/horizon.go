package extensions

import (
	"github.com/stellar/go/clients/horizonclient"
	"github.com/velo-protocol/DRSv1/node/app/environments"
	"net/http"
)

func GetHorizonClient() *horizonclient.Client {
	horizonClient := horizonclient.Client{
		HorizonURL: env.HorizonURL,
		HTTP:       http.DefaultClient,
	}

	return &horizonClient
}

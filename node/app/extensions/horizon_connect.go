package extensions

import (
	"github.com/stellar/go/clients/horizonclient"
	env "gitlab.com/velo-labs/cen/node/app/environments"
	"net/http"
)

func ConnectHorizon() *horizonclient.Client {
	horizonClient := horizonclient.Client{
		HorizonURL: env.HorizonUrl,
		HTTP:       http.DefaultClient,
	}

	return &horizonClient
}

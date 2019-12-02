package stellar

import (
	"fmt"
	"github.com/stellar/go/clients/horizonclient"
)

type stellar struct {
	FriendBotURL  string
	HorizonURL    string
	HorizonClient horizonclient.ClientInterface
}

func NewStellar(horizonClient *horizonclient.Client) Repository {
	return &stellar{
		HorizonURL:    horizonClient.HorizonURL,
		HorizonClient: horizonClient,
		FriendBotURL:  fmt.Sprintf("%s/friendbot?addr=%%s", horizonClient.HorizonURL),
	}
}

func NewStellarWithClientInterface(horizonUrl string, horizonClient horizonclient.ClientInterface) Repository {
	return &stellar{
		HorizonURL:    horizonUrl,
		HorizonClient: horizonClient,
		FriendBotURL:  fmt.Sprintf("%s/friendbot?addr=%%s", horizonUrl),
	}
}

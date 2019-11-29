package velo

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/console"
	"github.com/velo-protocol/DRSv1/libs/client"
	"google.golang.org/grpc"
	"net/http"
)

type velo struct {
	veloClient        vclient.ClientInterface
	veloNodeUrl       string
	horizonUrl        string
	networkPassphrase string
}

func NewVelo(veloNodeUrl string, horizonUrl string, networkPassphrase string) Repository {
	return &velo{
		veloClient:        nil,
		veloNodeUrl:       veloNodeUrl,
		horizonUrl:        horizonUrl,
		networkPassphrase: networkPassphrase,
	}
}

func (velo *velo) Client(keyPair *keypair.Full) vclient.ClientInterface {
	if velo.veloClient != nil {
		velo.veloClient.SetKeyPair(keyPair)
		return velo.veloClient
	}

	grpcConn, err := grpc.Dial(velo.veloNodeUrl, grpc.WithInsecure())
	if err != nil {
		console.ExitWithError(console.ExitBadConnection, errors.Wrap(err, "cannot connect to Velo Node"))
	}

	velo.veloClient = vclient.NewClient(&horizonclient.Client{
		HorizonURL: velo.horizonUrl,
		HTTP:       http.DefaultClient,
	}, velo.networkPassphrase, grpcConn, keyPair)
	return velo.veloClient
}

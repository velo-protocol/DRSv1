package vclient

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	cenGrpc "gitlab.com/velo-labs/cen/grpc"
	"gitlab.com/velo-labs/cen/libs/convert"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client struct {
	horizonClient     horizonclient.ClientInterface
	networkPassphrase string
	keyPair           *keypair.Full

	veloNodeClient cenGrpc.VeloNodeClient
	grpcConnection *grpc.ClientConn
}

type ClientInterface interface {
	Close() error

	Whitelist(ctx context.Context, veloOp vtxnbuild.Whitelist) (*horizon.TransactionSuccess, error)
	SetupCredit(ctx context.Context, veloOp vtxnbuild.SetupCredit) (*horizon.TransactionSuccess, error)
	PriceUpdate(ctx context.Context, veloOp vtxnbuild.PriceUpdate) (*horizon.TransactionSuccess, error)
	MintCredit(ctx context.Context, veloOp vtxnbuild.MintCredit) (*horizon.TransactionSuccess, error)
	RedeemCredit(ctx context.Context, veloOp vtxnbuild.RedeemCredit) (*horizon.TransactionSuccess, error)

	GetExchangeRate(ctx context.Context, request *cenGrpc.GetExchangeRateRequest) (*cenGrpc.GetExchangeRateRequest, error)
}

func NewDefaultPublicClient(veloNodeUrl string, stellarAccountSecretKey string) (*Client, error) {
	grpcConn, err := grpc.Dial(veloNodeUrl, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to Velo Node")
	}

	return NewPublicClient(grpcConn, stellarAccountSecretKey)
}

func NewDefaultTestNetClient(veloNodeUrl string, stellarAccountSecretKey string) (*Client, error) {
	grpcConn, err := grpc.Dial(veloNodeUrl, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to Velo Node")
	}

	return NewTestNetClient(grpcConn, stellarAccountSecretKey)
}

func NewPublicClient(grpcConn *grpc.ClientConn, stellarAccountSecretKey string) (*Client, error) {
	keyPair, err := vconvert.SecretKeyToKeyPair(stellarAccountSecretKey)
	if err != nil {
		return nil, err
	}

	return newClient(horizonclient.DefaultPublicNetClient, network.PublicNetworkPassphrase, cenGrpc.NewVeloNodeClient(grpcConn), grpcConn, keyPair), nil
}

func NewTestNetClient(grpcConn *grpc.ClientConn, stellarAccountSecretKey string) (*Client, error) {
	keyPair, err := vconvert.SecretKeyToKeyPair(stellarAccountSecretKey)
	if err != nil {
		return nil, err
	}

	return newClient(horizonclient.DefaultTestNetClient, network.TestNetworkPassphrase, cenGrpc.NewVeloNodeClient(grpcConn), grpcConn, keyPair), nil
}

func NewClient(horizonClient horizonclient.ClientInterface, networkPassphrase string, grpcConn *grpc.ClientConn, keyPair *keypair.Full) *Client {
	return newClient(horizonClient, networkPassphrase, cenGrpc.NewVeloNodeClient(grpcConn), grpcConn, keyPair)
}

func newClient(horizonClient horizonclient.ClientInterface, networkPassphrase string, veloNodeClient cenGrpc.VeloNodeClient, grpcConn *grpc.ClientConn, keyPair *keypair.Full) *Client {
	return &Client{
		horizonClient:     horizonClient,
		networkPassphrase: networkPassphrase,
		keyPair:           keyPair,
		veloNodeClient:    veloNodeClient,
		grpcConnection:    grpcConn,
	}
}

func (client *Client) Close() (err error) {
	return client.grpcConnection.Close()
}

func (client *Client) Whitelist(ctx context.Context, veloOp vtxnbuild.Whitelist) (*horizon.TransactionSuccess, error) {
	return client.executeVeloTx(ctx, &veloOp)
}

func (client *Client) SetupCredit(ctx context.Context, veloOp vtxnbuild.SetupCredit) (*horizon.TransactionSuccess, error) {
	return client.executeVeloTx(ctx, &veloOp)
}

func (client *Client) PriceUpdate(ctx context.Context, veloOp vtxnbuild.PriceUpdate) (*horizon.TransactionSuccess, error) {
	return client.executeVeloTx(ctx, &veloOp)
}

func (client *Client) MintCredit(ctx context.Context, veloOp vtxnbuild.MintCredit) (*horizon.TransactionSuccess, error) {
	return client.executeVeloTx(ctx, &veloOp)
}

func (client *Client) RedeemCredit(ctx context.Context, veloOp vtxnbuild.RedeemCredit) (*horizon.TransactionSuccess, error) {
	return client.executeVeloTx(ctx, &veloOp)
}

func (client *Client) executeVeloTx(ctx context.Context, veloOp vtxnbuild.VeloOp) (*horizon.TransactionSuccess, error) {
	veloTx := vtxnbuild.VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: client.keyPair.Address(),
		},
		VeloOp: veloOp,
	}

	signedVeloTxB64, err := veloTx.BuildSignEncode(client.keyPair)
	if err != nil {
		return nil, err
	}

	reply, err := client.veloNodeClient.SubmitVeloTx(ctx, &cenGrpc.VeloTxRequest{
		SignedVeloTxXdr: signedVeloTxB64,
	})
	if err != nil {
		veloNodeErr, ok := status.FromError(err)
		if ok {
			if veloNodeErr.Code() == codes.Unavailable {
				return nil, errors.Wrap(err, "cannot connect to Velo Node")
			}
		}

		return nil, err
	}

	tx, err := txnbuild.TransactionFromXDR(reply.SignedStellarTxXdr)
	if err != nil {
		return nil, err
	}
	tx.Network = client.networkPassphrase

	err = tx.Sign(client.keyPair)
	if err != nil {
		return nil, err
	}

	signedTxB64, err := tx.Base64()
	if err != nil {
		return nil, err
	}

	result, err := client.horizonClient.SubmitTransactionXDR(signedTxB64)
	if err != nil {
		herr, ok := err.(*horizonclient.Error)
		if ok {
			return nil, herr
		}
		return nil, errors.Wrap(err, "cannot connect to horizon")
	}

	return &result, nil
}

func (client *Client) GetExchangeRate(ctx context.Context, request *cenGrpc.GetExchangeRateRequest) (*cenGrpc.GetExchangeRateReply, error) {
	return client.veloNodeClient.GetExchangeRate(ctx, request)
}

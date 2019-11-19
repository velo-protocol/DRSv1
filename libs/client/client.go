package vclient

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	spec "github.com/velo-protocol/DRSv1/grpc"
	"github.com/velo-protocol/DRSv1/libs/convert"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client struct {
	horizonClient     horizonclient.ClientInterface
	networkPassphrase string
	keyPair           *keypair.Full

	veloNodeClient spec.VeloNodeClient
	grpcConnection *grpc.ClientConn
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

	return newClient(horizonclient.DefaultPublicNetClient, network.PublicNetworkPassphrase, spec.NewVeloNodeClient(grpcConn), grpcConn, keyPair), nil
}

func NewTestNetClient(grpcConn *grpc.ClientConn, stellarAccountSecretKey string) (*Client, error) {
	keyPair, err := vconvert.SecretKeyToKeyPair(stellarAccountSecretKey)
	if err != nil {
		return nil, err
	}

	return newClient(horizonclient.DefaultTestNetClient, network.TestNetworkPassphrase, spec.NewVeloNodeClient(grpcConn), grpcConn, keyPair), nil
}

func NewClient(horizonClient horizonclient.ClientInterface, networkPassphrase string, grpcConn *grpc.ClientConn, keyPair *keypair.Full) *Client {
	return newClient(horizonClient, networkPassphrase, spec.NewVeloNodeClient(grpcConn), grpcConn, keyPair)
}

func newClient(horizonClient horizonclient.ClientInterface, networkPassphrase string, veloNodeClient spec.VeloNodeClient, grpcConn *grpc.ClientConn, keyPair *keypair.Full) *Client {
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

func (client *Client) SetKeyPair(keyPair *keypair.Full) {
	client.keyPair = keyPair
}

type WhitelistResult struct {
	HorizonResult  *horizon.TransactionSuccess
	VeloNodeResult *spec.WhitelistOpResponse
}

func (client *Client) Whitelist(ctx context.Context, veloOp vtxnbuild.Whitelist) (WhitelistResult, error) {
	horizonSuccess, veloReply, err := client.executeVeloTx(ctx, &veloOp)
	var veloNodeResult *spec.WhitelistOpResponse
	if veloReply != nil {
		veloNodeResult = veloReply.WhitelistOpResponse
	}

	return WhitelistResult{
		HorizonResult:  horizonSuccess,
		VeloNodeResult: veloNodeResult,
	}, err
}

type SetupCreditResult struct {
	HorizonResult  *horizon.TransactionSuccess
	VeloNodeResult *spec.SetupCreditOpResponse
}

func (client *Client) SetupCredit(ctx context.Context, veloOp vtxnbuild.SetupCredit) (SetupCreditResult, error) {
	horizonSuccess, veloReply, err := client.executeVeloTx(ctx, &veloOp)
	var veloNodeResult *spec.SetupCreditOpResponse
	if veloReply != nil {
		veloNodeResult = veloReply.SetupCreditOpResponse
	}

	return SetupCreditResult{
		HorizonResult:  horizonSuccess,
		VeloNodeResult: veloNodeResult,
	}, err
}

type PriceUpdateResult struct {
	HorizonResult  *horizon.TransactionSuccess
	VeloNodeResult *spec.PriceUpdateOpResponse
}

func (client *Client) PriceUpdate(ctx context.Context, veloOp vtxnbuild.PriceUpdate) (PriceUpdateResult, error) {
	horizonSuccess, veloReply, err := client.executeVeloTx(ctx, &veloOp)
	var veloNodeResult *spec.PriceUpdateOpResponse
	if veloReply != nil {
		veloNodeResult = veloReply.PriceUpdateOpResponse
	}

	return PriceUpdateResult{
		HorizonResult:  horizonSuccess,
		VeloNodeResult: veloNodeResult,
	}, err
}

type MintCreditResult struct {
	HorizonResult  *horizon.TransactionSuccess
	VeloNodeResult *spec.MintCreditOpResponse
}

func (client *Client) MintCredit(ctx context.Context, veloOp vtxnbuild.MintCredit) (MintCreditResult, error) {
	horizonSuccess, veloReply, err := client.executeVeloTx(ctx, &veloOp)
	var veloNodeResult *spec.MintCreditOpResponse
	if veloReply != nil {
		veloNodeResult = veloReply.MintCreditOpResponse
	}

	return MintCreditResult{
		HorizonResult:  horizonSuccess,
		VeloNodeResult: veloNodeResult,
	}, err
}

type RedeemCreditResult struct {
	HorizonResult  *horizon.TransactionSuccess
	VeloNodeResult *spec.RedeemCreditOpResponse
}

func (client *Client) RedeemCredit(ctx context.Context, veloOp vtxnbuild.RedeemCredit) (RedeemCreditResult, error) {
	horizonSuccess, veloReply, err := client.executeVeloTx(ctx, &veloOp)
	var veloNodeResult *spec.RedeemCreditOpResponse
	if veloReply != nil {
		veloNodeResult = veloReply.RedeemCreditOpResponse
	}

	return RedeemCreditResult{
		HorizonResult:  horizonSuccess,
		VeloNodeResult: veloNodeResult,
	}, err
}

type RebalanceReserveResult struct {
	HorizonResult  *horizon.TransactionSuccess
	VeloNodeResult *spec.RebalanceReserveOpResponse
}

func (client *Client) RebalanceReserve(ctx context.Context, veloOp vtxnbuild.RebalanceReserve) (RebalanceReserveResult, error) {
	horizonSuccess, veloReply, err := client.executeVeloTx(ctx, &veloOp)
	var veloNodeResult *spec.RebalanceReserveOpResponse
	if veloReply != nil {
		veloNodeResult = veloReply.RebalanceReserveOpResponse
	}

	return RebalanceReserveResult{
		HorizonResult:  horizonSuccess,
		VeloNodeResult: veloNodeResult,
	}, err
}

func (client *Client) executeVeloTx(ctx context.Context, veloOp vtxnbuild.VeloOp) (*horizon.TransactionSuccess, *spec.VeloTxReply, error) {
	veloTx := vtxnbuild.VeloTx{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: client.keyPair.Address(),
		},
		VeloOp: veloOp,
	}

	signedVeloTxB64, err := veloTx.BuildSignEncode(client.keyPair)
	if err != nil {
		return nil, nil, err
	}

	reply, err := client.veloNodeClient.SubmitVeloTx(ctx, &spec.VeloTxRequest{
		SignedVeloTxXdr: signedVeloTxB64,
	})
	if err != nil {
		veloNodeErr, ok := status.FromError(err)
		if ok {
			if veloNodeErr.Code() == codes.Unavailable {
				return nil, nil, errors.Wrap(err, "cannot connect to Velo Node")
			}
		}

		return nil, nil, err
	}

	tx, err := txnbuild.TransactionFromXDR(reply.SignedStellarTxXdr)
	if err != nil {
		return nil, reply, err
	}
	tx.Network = client.networkPassphrase

	err = tx.Sign(client.keyPair)
	if err != nil {
		return nil, reply, err
	}

	signedTxB64, err := tx.Base64()
	if err != nil {
		return nil, reply, err
	}

	result, err := client.horizonClient.SubmitTransactionXDR(signedTxB64)
	if err != nil {
		herr, ok := err.(*horizonclient.Error)
		if ok {
			return nil, reply, herr
		}
		return nil, reply, errors.Wrap(err, "cannot connect to horizon")
	}

	return &result, reply, nil
}

func (client *Client) GetExchangeRate(ctx context.Context, request *spec.GetExchangeRateRequest) (*spec.GetExchangeRateReply, error) {
	return client.veloNodeClient.GetExchangeRate(ctx, request)
}

func (client *Client) GetCollateralHealthCheck(ctx context.Context, request *spec.GetCollateralHealthCheckRequest) (*spec.GetCollateralHealthCheckReply, error) {
	return client.veloNodeClient.GetCollateralHealthCheck(ctx, request)
}

package vclient

import (
	"context"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	cenGrpc "gitlab.com/velo-labs/cen/grpc"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
)

type ClientInterface interface {
	Close() error
	SetKeyPair(keyPair *keypair.Full)

	Whitelist(ctx context.Context, veloOp vtxnbuild.Whitelist) (*horizon.TransactionSuccess, error)
	SetupCredit(ctx context.Context, veloOp vtxnbuild.SetupCredit) (*horizon.TransactionSuccess, error)
	PriceUpdate(ctx context.Context, veloOp vtxnbuild.PriceUpdate) (*horizon.TransactionSuccess, error)
	MintCredit(ctx context.Context, veloOp vtxnbuild.MintCredit) (*horizon.TransactionSuccess, error)
	RedeemCredit(ctx context.Context, veloOp vtxnbuild.RedeemCredit) (*horizon.TransactionSuccess, error)
	RebalanceReserve(ctx context.Context, veloOp vtxnbuild.RebalanceReserve) (*horizon.TransactionSuccess, error)

	GetExchangeRate(ctx context.Context, request *cenGrpc.GetExchangeRateRequest) (*cenGrpc.GetExchangeRateReply, error)
	GetCollateralHealthCheck(ctx context.Context, request *cenGrpc.GetCollateralHealthCheckRequest) (*cenGrpc.GetCollateralHealthCheckReply, error)
}

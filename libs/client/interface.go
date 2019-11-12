package vclient

import (
	"context"
	"github.com/stellar/go/keypair"
	cenGrpc "gitlab.com/velo-labs/cen/grpc"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
)

type ClientInterface interface {
	Close() error
	SetKeyPair(keyPair *keypair.Full)

	Whitelist(ctx context.Context, veloOp vtxnbuild.Whitelist) (WhitelistResult, error)
	SetupCredit(ctx context.Context, veloOp vtxnbuild.SetupCredit) (SetupCreditResult, error)
	PriceUpdate(ctx context.Context, veloOp vtxnbuild.PriceUpdate) (PriceUpdateResult, error)
	MintCredit(ctx context.Context, veloOp vtxnbuild.MintCredit) (MintCreditResult, error)
	RedeemCredit(ctx context.Context, veloOp vtxnbuild.RedeemCredit) (RedeemCreditResult, error)
	RebalanceReserve(ctx context.Context, veloOp vtxnbuild.RebalanceReserve) (RebalanceReserveResult, error)

	GetExchangeRate(ctx context.Context, request *cenGrpc.GetExchangeRateRequest) (*cenGrpc.GetExchangeRateReply, error)
	GetCollateralHealthCheck(ctx context.Context, request *cenGrpc.GetCollateralHealthCheckRequest) (*cenGrpc.GetCollateralHealthCheckReply, error)
}

package vclient

import (
	"context"
	"github.com/stellar/go/keypair"
	spec "github.com/velo-protocol/DRSv1/grpc"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
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

	GetExchangeRate(ctx context.Context, request *spec.GetExchangeRateRequest) (*spec.GetExchangeRateReply, error)
	GetCollateralHealthCheck(ctx context.Context, request *spec.GetCollateralHealthCheckRequest) (*spec.GetCollateralHealthCheckReply, error)
}

package usecases

import (
	"context"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/node/app/entities"
	"gitlab.com/velo-labs/cen/node/app/errors"
)

type UseCase interface {
	SetupCredit(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*string, nerrors.NodeError)
	CreateWhitelist(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*string, nerrors.NodeError)
	UpdatePrice(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*string, nerrors.NodeError)
	MintCredit(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*entities.MintCreditOutput, nerrors.NodeError)
	RedeemCredit(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*entities.RedeemCreditOutput, nerrors.NodeError)
	GetExchangeRate(ctx context.Context, input *entities.GetExchangeRateInput) (*entities.GetExchangeRateOutPut, nerrors.NodeError)
	GetCollateralHealthCheck(ctx context.Context) (*entities.GetCollateralHealthCheckOutput, nerrors.NodeError)
	RebalanceReserve(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*entities.RebalanceOutput, nerrors.NodeError)
}

package usecases

import (
	"context"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"github.com/velo-protocol/DRSv1/node/app/entities"
	"github.com/velo-protocol/DRSv1/node/app/errors"
)

type UseCase interface {
	SetupCredit(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*entities.SetupCreditOutput, nerrors.NodeError)
	Whitelist(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*entities.WhitelistOutput, nerrors.NodeError)
	UpdatePrice(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*entities.UpdatePriceOutput, nerrors.NodeError)
	MintCredit(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*entities.MintCreditOutput, nerrors.NodeError)
	RedeemCredit(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*entities.RedeemCreditOutput, nerrors.NodeError)
	GetExchangeRate(ctx context.Context, input *entities.GetExchangeRateInput) (*entities.GetExchangeRateOutPut, nerrors.NodeError)
	GetCollateralHealthCheck(ctx context.Context) (*entities.GetCollateralHealthCheckOutput, nerrors.NodeError)
	RebalanceReserve(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*entities.RebalanceOutput, nerrors.NodeError)
}

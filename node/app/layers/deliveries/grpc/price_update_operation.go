package grpc

import (
	"context"
	"fmt"
	spec "github.com/velo-protocol/DRSv1/grpc"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"github.com/velo-protocol/DRSv1/node/app/constants"
	"github.com/velo-protocol/DRSv1/node/app/errors"
)

func (handler *handler) handlePriceUpdateOperation(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*spec.VeloTxReply, error) {
	op := veloTx.TxEnvelope().VeloTx.VeloOp.Body.PriceUpdateOp
	if op == nil {
		return nil, nerrors.ErrInvalidArgument{
			Message: fmt.Sprintf(constants.ErrFormatMissingOperation, constants.VeloOpPriceUpdate),
		}.GRPCError()
	}

	priceUpdateOutput, err := handler.UseCase.UpdatePrice(ctx, veloTx)
	if err != nil {
		return nil, err.GRPCError()
	}

	return &spec.VeloTxReply{
		SignedStellarTxXdr: priceUpdateOutput.SignedStellarTxXdr,
		Message:            constants.ReplyPriceUpdateSuccess,
		PriceUpdateOpResponse: &spec.PriceUpdateOpResponse{
			Asset:                       priceUpdateOutput.Asset,
			Currency:                    priceUpdateOutput.Currency,
			PriceInCurrencyPerAssetUnit: priceUpdateOutput.PriceInCurrencyPerAssetUnit,
		},
	}, nil
}

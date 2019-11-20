package grpc

import (
	"context"
	spec "github.com/velo-protocol/DRSv1/grpc"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"github.com/velo-protocol/DRSv1/libs/xdr"
	"github.com/velo-protocol/DRSv1/node/app/constants"
	"github.com/velo-protocol/DRSv1/node/app/errors"
)

func (handler *handler) SubmitVeloTx(ctx context.Context, req *spec.VeloTxRequest) (*spec.VeloTxReply, error) {
	veloTx, err := vtxnbuild.TransactionFromXDR(req.GetSignedVeloTxXdr())
	if err != nil {
		return nil, nerrors.ErrInvalidArgument{
			Message: err.Error(),
		}.GRPCError()
	}

	switch veloTx.TxEnvelope().VeloTx.VeloOp.Body.Type {
	case vxdr.OperationTypeWhitelist:
		return handler.handleWhitelistOperation(ctx, &veloTx)
	case vxdr.OperationTypeSetupCredit:
		return handler.handleSetupCreditOperation(ctx, &veloTx)
	case vxdr.OperationTypePriceUpdate:
		return handler.handlePriceUpdateOperation(ctx, &veloTx)
	case vxdr.OperationTypeMintCredit:
		return handler.handleMintCreditOperation(ctx, &veloTx)
	case vxdr.OperationTypeRedeemCredit:
		return handler.handleRedeemCreditOperation(ctx, &veloTx)
	case vxdr.OperationTypeRebalanceReserve:
		return handler.handleRebalanceReserve(ctx, &veloTx)
	default: // this case should never occur, if the cen/libs and cen/node is aligned
		return nil, nerrors.ErrInvalidArgument{
			Message: constants.ErrUnknownVeloOperationType,
		}.GRPCError()
	}

}

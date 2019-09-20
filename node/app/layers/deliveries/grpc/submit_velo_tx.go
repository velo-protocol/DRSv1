package grpc

import (
	"context"
	"fmt"
	spec "gitlab.com/velo-labs/cen/grpc"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/errors"
)

func (handler *handler) SubmitVeloTx(ctx context.Context, req *spec.VeloTxRequest) (*spec.VeloTxReply, error) {
	veloTx, err := vtxnbuild.TransactionFromXDR(req.GetSignedVeloTxXdr())
	if err != nil {
		return nil, err
	}

	veloTxEnvelope := veloTx.TxEnvelope()
	switch veloTxEnvelope.VeloTx.VeloOp.Body.Type {
	case vxdr.OperationTypeWhiteList:
		return handler.handleWhiteListOperation(ctx, veloTxEnvelope)
	default: // this case should never occur, if the cen/libs and cen/node is aligned
		return nil, nerrors.ErrInvalidArgument{
			Message: constants.ErrUnknownVeloOperationType,
		}.GRPCError()
	}

}

func (handler *handler) handleWhiteListOperation(ctx context.Context, veloTxEnvelope *vxdr.VeloTxEnvelope) (*spec.VeloTxReply, error) {
	if veloTxEnvelope.VeloTx.VeloOp.Body.WhiteListOp == nil {
		return nil, nerrors.ErrInvalidArgument{
			Message: fmt.Sprintf(constants.ErrFormatMissingOperation, "whitelist"),
		}.GRPCError()
	}

	err := handler.UseCase.CreateWhiteList(ctx, veloTxEnvelope)
	if err != nil {
		return nil, err.GRPCError()
	}

	return &spec.VeloTxReply{SignedStellarTxXdr: ""}, nil
}

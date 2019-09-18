package grpc

import (
	"context"
	spec "gitlab.com/velo-labs/cen/grpc"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/libs/xdr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, status.Errorf(codes.InvalidArgument, "unknown velo operation type")
	}

}

func (handler *handler) handleWhiteListOperation(ctx context.Context, veloTxEnvelope *vxdr.VeloTxEnvelope) (*spec.VeloTxReply, error) {
	if veloTxEnvelope.VeloTx.VeloOp.Body.WhiteListOp == nil {
		return nil, status.Errorf(codes.InvalidArgument, "operation type %s is missing", "whiteList")
	}

	err := handler.UseCase.CreateWhiteList(ctx, veloTxEnvelope)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &spec.VeloTxReply{SignedStellarTxXdr: ""}, nil
}

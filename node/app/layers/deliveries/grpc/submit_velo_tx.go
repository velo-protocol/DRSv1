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
		return nil, nerrors.ErrInvalidArgument{
			Message: err.Error(),
		}.GRPCError()
	}

	switch veloTx.TxEnvelope().VeloTx.VeloOp.Body.Type {
	case vxdr.OperationTypeWhiteList:
		return handler.handleWhiteListOperation(ctx, &veloTx)
	case vxdr.OperationTypeSetupCredit:
		return handler.handleSetupCreditOperation(ctx, &veloTx)
	default: // this case should never occur, if the cen/libs and cen/node is aligned
		return nil, nerrors.ErrInvalidArgument{
			Message: constants.ErrUnknownVeloOperationType,
		}.GRPCError()
	}

}

func (handler *handler) handleWhiteListOperation(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*spec.VeloTxReply, error) {
	op := veloTx.TxEnvelope().VeloTx.VeloOp.Body.WhiteListOp
	if op == nil {
		return nil, nerrors.ErrInvalidArgument{
			Message: fmt.Sprintf(constants.ErrFormatMissingOperation, constants.VeloOpWhiteList),
		}.GRPCError()
	}

	err := handler.UseCase.CreateWhiteList(ctx, veloTx)
	if err != nil {
		return nil, err.GRPCError()
	}

	return &spec.VeloTxReply{
		Message: fmt.Sprintf(constants.ReplyWhiteListSuccess, op.Address.Address(), op.Role),
	}, nil
}

func (handler *handler) handleSetupCreditOperation(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*spec.VeloTxReply, error) {
	op := veloTx.TxEnvelope().VeloTx.VeloOp.Body.SetupCreditOp
	if op == nil {
		return nil, nerrors.ErrInvalidArgument{
			Message: fmt.Sprintf(constants.ErrFormatMissingOperation, constants.VeloOpSetupCredit),
		}.GRPCError()
	}

	signedStellarTxXdr, err := handler.UseCase.SetupCredit(ctx, veloTx)
	if err != nil {
		return nil, err.GRPCError()
	}

	return &spec.VeloTxReply{
		SignedStellarTxXdr: *signedStellarTxXdr,
		Message:            constants.ReplySetupCreditSuccess,
	}, nil
}

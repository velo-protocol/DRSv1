package grpc

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	spec "gitlab.com/velo-labs/cen/grpc"
	"gitlab.com/velo-labs/cen/libs/txnbuild"
	"gitlab.com/velo-labs/cen/node/app/constants"
	"gitlab.com/velo-labs/cen/node/app/errors"
)

func (handler handler) RebalanceReserve(ctx context.Context, req *spec.RebalanceReserveRequest) (*spec.RebalanceReserveReply, error) {
	veloTx, vtxErr := vtxnbuild.TransactionFromXDR(req.GetSignedVeloTxXdr())
	if vtxErr != nil {
		return nil, nerrors.ErrInvalidArgument{
			Message: vtxErr.Error(),
		}.GRPCError()
	}

	op := veloTx.TxEnvelope().VeloTx.VeloOp.Body.RebalanceReserveOp
	if op == nil {
		return nil, nerrors.ErrInvalidArgument{
			Message: fmt.Sprintf(constants.ErrFormatMissingOperation, constants.VeloOpRebalanceReserve),
		}.GRPCError()
	}

	rebalanceReserve, err := handler.UseCase.RebalanceReserve(ctx, &veloTx)
	if err != nil {
		return nil, err
	}

	_ = copier.Copy(&spec.RebalanceReserveReply{}, rebalanceReserve)
	return &spec.RebalanceReserveReply{}, nil
}

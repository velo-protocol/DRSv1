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

func (handler *handler) handleWhitelistOperation(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*spec.VeloTxReply, error) {
	op := veloTx.TxEnvelope().VeloTx.VeloOp.Body.WhitelistOp
	if op == nil {
		return nil, nerrors.ErrInvalidArgument{
			Message: fmt.Sprintf(constants.ErrFormatMissingOperation, constants.VeloOpWhitelist),
		}.GRPCError()
	}

	signedStellarTxXdr, err := handler.UseCase.CreateWhitelist(ctx, veloTx)
	if err != nil {
		return nil, err.GRPCError()
	}

	return &spec.VeloTxReply{
		SignedStellarTxXdr: *signedStellarTxXdr,
		Message:            fmt.Sprintf(constants.ReplyWhitelistSuccess, op.Address.Address(), vxdr.RoleMap[op.Role]),
	}, nil
}

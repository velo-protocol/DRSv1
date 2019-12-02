package grpc

import (
	"context"
	"fmt"
	spec "github.com/velo-protocol/DRSv1/grpc"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"github.com/velo-protocol/DRSv1/libs/xdr"
	"github.com/velo-protocol/DRSv1/node/app/constants"
	"github.com/velo-protocol/DRSv1/node/app/errors"
)

func (handler *handler) handleWhitelistOperation(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*spec.VeloTxReply, error) {
	op := veloTx.TxEnvelope().VeloTx.VeloOp.Body.WhitelistOp
	if op == nil {
		return nil, nerrors.ErrInvalidArgument{
			Message: fmt.Sprintf(constants.ErrFormatMissingOperation, constants.VeloOpWhitelist),
		}.GRPCError()
	}

	output, err := handler.UseCase.Whitelist(ctx, veloTx)
	if err != nil {
		return nil, err.GRPCError()
	}

	return &spec.VeloTxReply{
		SignedStellarTxXdr: output.SignedStellarTxXdr,
		Message:            fmt.Sprintf(constants.ReplyWhitelistSuccess, output.Address, vxdr.RoleMap[vxdr.Role(output.Role)]),
		WhitelistOpResponse: &spec.WhitelistOpResponse{
			Address:                   output.Address,
			Role:                      output.Role,
			Currency:                  output.Currency,
			TrustedPartnerMetaAddress: output.TrustedPartnerMetaAddress,
		},
	}, nil
}

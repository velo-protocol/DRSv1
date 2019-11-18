package grpc

import (
	"context"
	"fmt"
	spec "github.com/velo-protocol/DRSv1/grpc"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"github.com/velo-protocol/DRSv1/node/app/constants"
	"github.com/velo-protocol/DRSv1/node/app/errors"
)

func (handler *handler) handleMintCreditOperation(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*spec.VeloTxReply, error) {
	op := veloTx.TxEnvelope().VeloTx.VeloOp.Body.MintCreditOp
	if op == nil {
		return nil, nerrors.ErrInvalidArgument{
			Message: fmt.Sprintf(constants.ErrFormatMissingOperation, constants.VeloOpMintCredit),
		}.GRPCError()
	}

	mintCreditOutput, err := handler.UseCase.MintCredit(ctx, veloTx)
	if err != nil {
		return nil, err.GRPCError()
	}

	return &spec.VeloTxReply{
		SignedStellarTxXdr: mintCreditOutput.SignedStellarTxXdr,
		Message: fmt.Sprintf(
			constants.ReplyMintCreditSuccess,
			mintCreditOutput.MintAmount.Truncate(7).StringFixed(7),
			mintCreditOutput.MintCurrency,
			mintCreditOutput.CollateralAmount.Truncate(7).StringFixed(7),
			mintCreditOutput.CollateralAsset,
		),
	}, nil
}

package grpc

import (
	"context"
	"fmt"
	spec "github.com/velo-protocol/DRSv1/grpc"
	"github.com/velo-protocol/DRSv1/libs/txnbuild"
	"github.com/velo-protocol/DRSv1/node/app/constants"
	"github.com/velo-protocol/DRSv1/node/app/errors"
)

func (handler *handler) handleSetupCreditOperation(ctx context.Context, veloTx *vtxnbuild.VeloTx) (*spec.VeloTxReply, error) {
	op := veloTx.TxEnvelope().VeloTx.VeloOp.Body.SetupCreditOp
	if op == nil {
		return nil, nerrors.ErrInvalidArgument{
			Message: fmt.Sprintf(constants.ErrFormatMissingOperation, constants.VeloOpSetupCredit),
		}.GRPCError()
	}

	setupCreditOutput, err := handler.UseCase.SetupCredit(ctx, veloTx)
	if err != nil {
		return nil, err.GRPCError()
	}

	return &spec.VeloTxReply{
		SignedStellarTxXdr: setupCreditOutput.SignedStellarTxXdr,
		Message:            constants.ReplySetupCreditSuccess,
		SetupCreditOpResponse: &spec.SetupCreditOpResponse{
			AssetIssuer:      setupCreditOutput.AssetIssuer,
			AssetDistributor: setupCreditOutput.AssetDistributor,
			AssetCode:        setupCreditOutput.AssetCode,
			PeggedValue:      setupCreditOutput.PeggedValue,
			PeggedCurrency:   setupCreditOutput.PeggedCurrency,
		},
	}, nil
}

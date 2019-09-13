package grpc

import (
	"context"
	spec "gitlab.com/velo-labs/cen/grpc"
)

func (handler *handler) SubmitVeloTransaction(context.Context, *spec.VeloTxRequest) (*spec.VeloTxReply, error) {
	return &spec.VeloTxReply{
		ResultVeloXdr: "AAA",
	}, nil
}

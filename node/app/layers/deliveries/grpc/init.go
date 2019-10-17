package grpc

import (
	spec "gitlab.com/velo-labs/cen/grpc"
	"gitlab.com/velo-labs/cen/node/app/layers/usecases"
	"google.golang.org/grpc"
)

type handler struct {
	UseCase usecases.UseCase
}

func Init(grpcServer *grpc.Server, useCase usecases.UseCase) {
	h := InitHandler(useCase)

	// Register Velo Node
	spec.RegisterVeloNodeServer(grpcServer, h)
}

func InitHandler(useCase usecases.UseCase) *handler {
	return &handler{
		UseCase: useCase,
	}
}

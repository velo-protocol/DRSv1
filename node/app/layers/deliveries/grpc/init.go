package grpc

import (
	spec "github.com/velo-protocol/DRSv1/grpc"
	"github.com/velo-protocol/DRSv1/node/app/layers/usecases"
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

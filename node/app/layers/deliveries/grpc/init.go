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
	handler := handler{
		UseCase: useCase,
	}

	// Register Velo Node
	spec.RegisterVeloNodeServer(grpcServer, &handler)
}

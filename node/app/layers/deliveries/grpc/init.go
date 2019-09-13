package grpc

import (
	spec "gitlab.com/velo-labs/cen/grpc"
	"gitlab.com/velo-labs/cen/node/app/layers/usecases"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type handler struct {
	UseCase usecases.UseCase
}

func Init(grpcServer *grpc.Server, useCase usecases.UseCase) {
	handler := handler{
		UseCase: useCase,
	}

	// Register Health check
	grpc_health_v1.RegisterHealthServer(grpcServer, &handler)

	// Register Velo Node
	spec.RegisterVeloNodeServer(grpcServer, &handler)
}

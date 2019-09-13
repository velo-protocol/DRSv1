package main

import (
	"fmt"
	"gitlab.com/velo-labs/cen/node/app/environments"
	"gitlab.com/velo-labs/cen/node/app/extensions"
	grpcDelivery "gitlab.com/velo-labs/cen/node/app/layers/deliveries/grpc"
	_stellarRepo "gitlab.com/velo-labs/cen/node/app/layers/repositories/stellar"
	"gitlab.com/velo-labs/cen/node/app/layers/usecases"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	grpcServer := grpc.NewServer(
		// TODO: Add auth, log, correlation, etc. middleware?
		//grpc.UnaryInterceptor(
		//	grpc_middleware.ChainUnaryServer(
		//		// TODO: Each middleware goes here
		//	),
		//),
	)

	// Extensions
	horizonClient := extensions.GetHorizonClient()
	dbConn := extensions.ConnectDB()
	defer dbConn.Close()

	// Repo
	stellarRepo := _stellarRepo.Init(horizonClient)

	// Use Cases
	useCase := usecases.Init(stellarRepo)

	// Deliveries
	grpcDelivery.Init(grpcServer, useCase)

	initServer(grpcServer)
}

func initServer(grpcServer *grpc.Server) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", env.Port))
	if err != nil {
		panic(err)
	}

	log.Printf("Server is starting at port %s", env.Port)
	err = grpcServer.Serve(listener)
	if err != nil {
		panic(err)
	}
}

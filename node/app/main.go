package main

import (
	"context"
	"fmt"
	"github.com/velo-protocol/DRSv1/node/app/environments"
	"github.com/velo-protocol/DRSv1/node/app/extensions"
	grpcDelivery "github.com/velo-protocol/DRSv1/node/app/layers/deliveries/grpc"
	_stellarRepo "github.com/velo-protocol/DRSv1/node/app/layers/repositories/stellar"
	"github.com/velo-protocol/DRSv1/node/app/layers/subusecases"
	"github.com/velo-protocol/DRSv1/node/app/layers/usecases"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	grpcHealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	grpcHealthServer *grpc.Server
)

func main() {
	env.Init()

	ctx := context.Background()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	group, ctx := errgroup.WithContext(ctx)

	grpcServer := grpc.NewServer(
	// TODO: Add auth, log, correlation, etc. middleware?
	//grpc.UnaryInterceptor(
	//	grpc_middleware.ChainUnaryServer(
	//		// TODO: Each middleware goes here
	//	),
	//),
	)

	healthServer := health.NewServer()
	initHealthServer(group, healthServer)

	// Extensions
	horizonClient := extensions.GetHorizonClient()

	// Repo
	stellarRepo := _stellarRepo.Init(horizonClient)

	// Sub Use Cases
	subUseCase := subusecases.Init(stellarRepo)

	// Use Cases
	useCase := usecases.Init(stellarRepo, subUseCase)

	// Deliveries
	grpcDelivery.Init(grpcServer, useCase)

	initServer(group, grpcServer, healthServer)

	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}

	log.Println("received shutdown signal")

	healthServer.SetServingStatus("", grpcHealth.HealthCheckResponse_NOT_SERVING)

	if grpcServer != nil {
		grpcServer.GracefulStop()
	}
	if grpcHealthServer != nil {
		grpcHealthServer.GracefulStop()
	}
	err := group.Wait()
	if err != nil {
		log.Printf("server returning an error: %s", err)
		os.Exit(2)
	}
}

func initServer(group *errgroup.Group, grpcServer *grpc.Server, healthServer *health.Server) {
	if env.EnableReflectionAPI {
		reflection.Register(grpcServer)
	}

	group.Go(func() error {
		addr := fmt.Sprintf(":%d", env.Port)
		listen, err := net.Listen("tcp", addr)
		if err != nil {
			log.Printf("gRPC server: failed to listen: %s", err)
			os.Exit(2)
		}
		log.Printf("gRPC server serving at %s", addr)
		healthServer.SetServingStatus("", grpcHealth.HealthCheckResponse_SERVING)
		return grpcServer.Serve(listen)
	})
}

func initHealthServer(group *errgroup.Group, healthServer *health.Server) {
	group.Go(func() error {
		grpcHealthServer = grpc.NewServer()

		grpcHealth.RegisterHealthServer(grpcHealthServer, healthServer)

		healthAddr := fmt.Sprintf(":%d", env.HealthPort)
		healthListen, err := net.Listen("tcp", healthAddr)
		if err != nil {
			log.Printf("gRPC Health server: failed to listen: %s", err)
			os.Exit(2)
		}
		log.Printf("gRPC health server serving at %s", healthAddr)
		return grpcHealthServer.Serve(healthListen)
	})
}

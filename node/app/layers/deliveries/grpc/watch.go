package grpc

import (
	"google.golang.org/grpc/health/grpc_health_v1"
)

func (handler *handler) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	return nil
}

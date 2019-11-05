package grpc

import (
	"context"
	spec "gitlab.com/velo-labs/cen/grpc"
	"log"
)

func (handler handler) GetCollateralHealthCheck(ctx context.Context, empty *spec.GetCollateralHealthCheckRequest) (*spec.GetCollateralHealthCheckReply, error) {
	log.Println("Welcome get collateral health check")

	getCollateralHealthCheck, err := handler.UseCase.GetCollateralHealthCheck(ctx)
	if err != nil {
		log.Println(err)
		return nil, err.GRPCError()
	}

	log.Println("result:", getCollateralHealthCheck)

	return &spec.GetCollateralHealthCheckReply{
		AssetCode:      getCollateralHealthCheck.AssetCode,
		AssetIssuer:    getCollateralHealthCheck.AssetIssuer,
		RequiredAmount: getCollateralHealthCheck.RequiredAmount.StringFixed(7),
		PoolAmount:     getCollateralHealthCheck.PoolAmount.StringFixed(7),
	}, nil
}

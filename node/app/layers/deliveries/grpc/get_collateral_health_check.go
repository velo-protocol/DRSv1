package grpc

import (
	"context"
	spec "gitlab.com/velo-labs/cen/grpc"
)

func (handler handler) GetCollateralHealthCheck(ctx context.Context, empty *spec.GetCollateralHealthCheckRequest) (*spec.GetCollateralHealthCheckReply, error) {

	getCollateralHealthCheck, err := handler.UseCase.GetCollateralHealthCheck(ctx)
	if err != nil {
		return nil, err.GRPCError()
	}

	return &spec.GetCollateralHealthCheckReply{
		AssetCode:      getCollateralHealthCheck.AssetCode,
		AssetIssuer:    getCollateralHealthCheck.AssetIssuer,
		RequiredAmount: getCollateralHealthCheck.RequiredAmount.StringFixed(7),
		PoolAmount:     getCollateralHealthCheck.PoolAmount.StringFixed(7),
	}, nil
}

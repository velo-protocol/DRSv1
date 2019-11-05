package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/velo-labs/cen/cmd/gvel/entity"
	spec "gitlab.com/velo-labs/cen/grpc"
)

func (lo *logic) GetCollateralHealthCheck() (*entity.GetCollateralHealthCheckOutput, error) {
	defaultAccount := lo.AppConfig.GetDefaultAccount()
	accountBytes, err := lo.DB.Get([]byte(defaultAccount))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account from db")
	}

	var account entity.StellarAccount
	err = json.Unmarshal(accountBytes, &account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal account")
	}

	result, err := lo.Velo.Client(nil).GetCollateralHealthCheck(context.Background(), &spec.GetCollateralHealthCheckRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get collateral health check")
	}

	asset := fmt.Sprintf("%s (%s...)", result.AssetCode, result.AssetIssuer[0:4])

	return &entity.GetCollateralHealthCheckOutput{
		Asset:          asset,
		RequiredAmount: result.RequiredAmount,
		PoolAmount:     result.PoolAmount,
	}, nil

}

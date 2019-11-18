package logic

import (
	"github.com/pkg/errors"
	"github.com/velo-protocol/DRSv1/cmd/gvel/entity"
)

func (lo *logic) SetDefaultAccount(input *entity.SetDefaultAccountInput) (*entity.SetDefaultAccountOutput, error) {
	_, err := lo.DB.Get([]byte(input.Account))
	if err != nil {
		return nil, errors.Wrapf(err, "address %s is not found in gvel", input.Account)
	}

	err = lo.AppConfig.SetDefaultAccount(input.Account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to write config file")
	}

	return &entity.SetDefaultAccountOutput{
		Account: input.Account,
	}, nil
}

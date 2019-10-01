package models

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"gitlab.com/velo-labs/cen/node/app/entities"
)

type DrsAccountData map[string]string

func (model *DrsAccountData) Entity() (*entities.DrsAccountData, error) {
	drsAccountDataEntity := new(entities.DrsAccountData)
	err := mapstructure.Decode(model, drsAccountDataEntity)
	if err != nil {
		return nil, errors.Wrap(err, "fail to map drs account data to entity")
	}

	err = drsAccountDataEntity.DecodeBase64()
	if err != nil {
		return nil, err
	}
	return drsAccountDataEntity, err
}

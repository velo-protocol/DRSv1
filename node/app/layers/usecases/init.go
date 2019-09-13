package usecases

import (
	"gitlab.com/velo-labs/cen/node/app/layers/repositories/stellar"
)

type useCase struct {
	StellarRepo stellar.Repo
}

func Init(stellarRepo stellar.Repo) UseCase {
	return &useCase{
		StellarRepo: stellarRepo,
	}
}

package usecases

import (
	"gitlab.com/velo-labs/cen/node/app/layers/repositories/stellar"
	"gitlab.com/velo-labs/cen/node/app/layers/subusecases"
)

type useCase struct {
	StellarRepo stellar.Repo
	SubUseCase  subusecases.SubUseCase
}

func Init(stellarRepo stellar.Repo, subUseCase subusecases.SubUseCase) UseCase {
	return &useCase{
		StellarRepo: stellarRepo,
		SubUseCase:  subUseCase,
	}
}

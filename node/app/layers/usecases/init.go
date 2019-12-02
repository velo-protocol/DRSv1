package usecases

import (
	"github.com/velo-protocol/DRSv1/node/app/layers/repositories/stellar"
	"github.com/velo-protocol/DRSv1/node/app/layers/subusecases"
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

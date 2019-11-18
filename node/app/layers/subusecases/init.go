package subusecases

import (
	"github.com/velo-protocol/DRSv1/node/app/layers/repositories/stellar"
)

type subUseCase struct {
	StellarRepo stellar.Repo
}

func Init(stellarRepo stellar.Repo) SubUseCase {
	return &subUseCase{
		StellarRepo: stellarRepo,
	}
}

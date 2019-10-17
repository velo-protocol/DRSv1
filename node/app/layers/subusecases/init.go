package subusecases

import (
	"gitlab.com/velo-labs/cen/node/app/layers/repositories/stellar"
)

type subUseCase struct {
	StellarRepo stellar.Repo
}

func Init(stellarRepo stellar.Repo) SubUseCase {
	return &subUseCase{
		StellarRepo: stellarRepo,
	}
}

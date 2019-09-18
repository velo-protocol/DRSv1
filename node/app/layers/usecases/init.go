package usecases

import (
	"gitlab.com/velo-labs/cen/node/app/layers/repositories/stellar"
	"gitlab.com/velo-labs/cen/node/app/layers/repositories/whitelist"
)

type useCase struct {
	StellarRepo   stellar.Repo
	WhitelistRepo whitelist.Repo
}

func Init(stellarRepo stellar.Repo, whitelistRepo whitelist.Repo) UseCase {
	return &useCase{
		StellarRepo:   stellarRepo,
		WhitelistRepo: whitelistRepo,
	}
}

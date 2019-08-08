package usecases

import (
	"gitlab.com/velo-labs/cen/app/modules/client"
	"gitlab.com/velo-labs/cen/app/modules/stellar"
	"gitlab.com/velo-labs/cen/app/services"
)

type usecase struct {
	Drsops            services.Operation
	StellarRepository stellar.Repository
}

func NewUseCase(drsops services.Operation, stellarRepository stellar.Repository) client.UseCase {
	return &usecase{
		Drsops:            drsops,
		StellarRepository: stellarRepository,
	}
}

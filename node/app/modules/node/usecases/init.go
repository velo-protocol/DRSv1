package usecases

import (
	"gitlab.com/velo-labs/cen/node/app/modules/node"
	"gitlab.com/velo-labs/cen/node/app/modules/stellar"
	"gitlab.com/velo-labs/cen/node/app/services/operation"
)

type usecase struct {
	Drsops            operation.Interface
	NodeRepository    node.Repository
	StellarRepository stellar.Repository
}

func NewNodeUseCase(
	drsops operation.Interface,
	nodeRepository node.Repository,
	stellarRepository stellar.Repository,
) node.UseCase {
	return &usecase{
		Drsops:            drsops,
		NodeRepository:    nodeRepository,
		StellarRepository: stellarRepository,
	}
}

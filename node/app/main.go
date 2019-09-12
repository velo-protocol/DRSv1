package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/velo-labs/cen/node/app/extensions"
	_nodeHttps "gitlab.com/velo-labs/cen/node/app/modules/node/deliveries/https"
	_nodeRepository "gitlab.com/velo-labs/cen/node/app/modules/node/repositories"
	_nodeUsecase "gitlab.com/velo-labs/cen/node/app/modules/node/usecases"
	_stellarRepository "gitlab.com/velo-labs/cen/node/app/modules/stellar/repository"
	_stellarDrsops "gitlab.com/velo-labs/cen/node/app/services/operation/stellar-drs-operations"
)

func main() {
	ginEngine := gin.New()

	ginEngine.Use(gin.Recovery())
	ginEngine.Use(gin.Logger())

	levelConn := extensions.ConnLevelDB()
	defer levelConn.Close()

	horizonclient := extensions.ConnectHorizon()

	nodeRepository := _nodeRepository.NewNodeRepository(levelConn)
	stellarRepository := _stellarRepository.NewHorizonStellarRepository(horizonclient)

	drsops := _stellarDrsops.NewDrsOps(stellarRepository)

	nodeUsecase := _nodeUsecase.NewNodeUseCase(drsops, nodeRepository, stellarRepository)

	_nodeHttps.NewEndpointHttpHandler(ginEngine, nodeUsecase)
}

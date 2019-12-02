package main

import (
	"github.com/stellar/go/clients/horizonclient"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/commands"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/logic"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/repositories/database"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/repositories/stellar"
	"github.com/velo-protocol/DRSv1/cmd/gvel/layers/repositories/velo"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/config"
	"github.com/velo-protocol/DRSv1/cmd/gvel/utils/console"
)

func main() {
	appConfig := config.NewConfiguration()
	appConfig.LoadDefault()

	console.InitLogger()

	var logicInstance logic.Logic
	{
		if appConfig.Exists() {
			// db
			accountDbRepository, err := database.NewLevelDb(appConfig.GetAccountDbPath())
			if err != nil {
				console.ExitWithError(console.ExitError, err)
			}

			// stellar
			var horizonClient *horizonclient.Client
			if appConfig.GetIsTestNet() {
				horizonClient = horizonclient.DefaultTestNetClient
			} else {
				horizonClient = horizonclient.DefaultPublicNetClient
			}
			horizonClient.HorizonURL = appConfig.GetHorizonUrl()
			stellarRepository := stellar.NewStellar(horizonClient)

			// velo
			veloRepository := velo.NewVelo(appConfig.GetVeloNodeUrl(), appConfig.GetHorizonUrl(), appConfig.GetNetworkPassphrase())

			// logic
			logicInstance = logic.NewLogic(accountDbRepository, stellarRepository, veloRepository, appConfig)

		} else {
			logicInstance = logic.NewLogic(&database.LevelDbDatabase{}, nil, nil, appConfig)
		}
	}

	commandHandler := commands.NewGvelHandler(logicInstance, appConfig)
	commandHandler.Init()

	err := commandHandler.RootCommand.Execute()
	if err != nil {
		console.ExitWithError(console.ExitError, err)
	}
}

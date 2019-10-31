package main

import (
	"github.com/stellar/go/clients/horizonclient"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/commands"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/logic"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/repositories/database"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/repositories/stellar"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/repositories/velo"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/config"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
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

package main

import (
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/commands"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/logic"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/repositories/database"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/repositories/friendbot"
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

			// friend bot
			friendBotRepository := friendbot.NewFriendBot(appConfig.GetFriendBotUrl())

			// velo
			veloRepository := velo.NewVelo(appConfig.GetVeloNodeUrl(), appConfig.GetHorizonUrl(), appConfig.GetNetworkPassphrase())

			// logic
			logicInstance = logic.NewLogic(accountDbRepository, friendBotRepository, veloRepository, appConfig)

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

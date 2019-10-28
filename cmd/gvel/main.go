package main

import (
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/commands"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/logic"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/repositories/database"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/repositories/friendbot"
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
			accountDbRepository, err := database.NewLevelDb(appConfig.GetAccountDbPath())
			if err != nil {
				console.ExitWithError(console.ExitError, err)
			}
			friendBotRepository := friendbot.NewFriendBot(appConfig.GetFriendBotUrl())

			logicInstance = logic.NewLogic(accountDbRepository, friendBotRepository, appConfig)
		} else {
			logicInstance = logic.NewLogic(&database.LevelDbDatabase{}, nil, appConfig)
		}
	}

	commandHandler := commands.NewGvelHandler(logicInstance, appConfig)
	commandHandler.Init()

	err := commandHandler.RootCommand.Execute()
	if err != nil {
		console.ExitWithError(console.ExitError, err)
	}
}

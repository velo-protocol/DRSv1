package main

import (
	"fmt"
	"github.com/spf13/viper"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/commands"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/logic"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/repositories/database"
	"gitlab.com/velo-labs/cen/cmd/gvel/layers/repositories/friendbot"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/config"
	"gitlab.com/velo-labs/cen/cmd/gvel/utils/console"
	"os"
)

func main() {
	config.Load()
	console.InitLogger()

	var logicInstance logic.Logic
	{
		if config.Exists() {
			accountDbRepository, err := database.NewLevelDbDatabase(viper.GetString("accountDbPath"))
			if err != nil {
				panic(err)
			}
			friendBotRepository := friendbot.NewFriendBot(viper.GetString("friendBotUrl"))

			logicInstance = logic.NewLogic(accountDbRepository, friendBotRepository)
		} else {
			logicInstance = logic.NewLogic(nil, nil)
		}
	}

	commandHandler := commands.NewGvelHandler(logicInstance)
	commandHandler.Init()

	err := commandHandler.RootCommand.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

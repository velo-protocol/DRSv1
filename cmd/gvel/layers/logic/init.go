package logic

func (lo *logic) Init(configFilePath string) error {
	err := lo.AppConfig.InitConfigFile(configFilePath)
	if err != nil {
		return err
	}

	accountDbPath := lo.AppConfig.GetAccountDbPath()

	err = lo.DB.Init(accountDbPath)
	if err != nil {
		return err
	}

	return nil
}

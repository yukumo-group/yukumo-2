package api

func (a *App) Config() (Config, error) {
	cfg, err := a.config.Get()
	if err != nil {
		return Config{}, err
	}
	return toConfig(cfg), nil
}

func (a *App) SaveConfig(c Config) error {
	return a.config.Save(fromConfig(c))
}

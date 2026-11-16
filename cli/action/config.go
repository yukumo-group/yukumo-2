package action

import (
	"context"

	"yukumo-script/cli/schema"
	"yukumo-script/core/api"
)

func ListConfig(_ context.Context, app *api.App) (api.Config, error) {
	return app.Config()
}

func GetConfigField(_ context.Context, app *api.App, key string) (string, error) {
	cfg, err := app.Config()
	if err != nil {
		return "", err
	}
	entry, err := schema.Lookup(schema.ConfigFields(), key)
	if err != nil {
		return "", err
	}
	return entry.Get(cfg), nil
}

func ApplyConfigEdit(_ context.Context, app *api.App, key, raw string, normalize bool) (EditResult, error) {
	cfg, err := app.Config()
	if err != nil {
		return EditResult{}, err
	}
	entry, err := schema.Lookup(schema.ConfigFields(), key)
	if err != nil {
		return EditResult{}, err
	}
	value := raw
	var warning string
	if normalize {
		value, warning, err = schema.Normalize(entry.Field, raw)
		if err != nil {
			return EditResult{}, err
		}
	}
	if err := entry.Set(&cfg, value); err != nil {
		return EditResult{}, err
	}
	if err := app.SaveConfig(cfg); err != nil {
		return EditResult{}, err
	}
	return EditResult{Value: entry.Get(cfg), Warning: warning}, nil
}

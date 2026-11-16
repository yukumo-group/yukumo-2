package action

import (
	"context"

	"yukumo-script/cli/schema"
	"yukumo-script/core/api"
)

func ListCharacters(_ context.Context, app *api.App) ([]api.Character, error) {
	return app.ListCharacters()
}

func GetCharacter(_ context.Context, app *api.App, name string) (api.Character, error) {
	return app.GetCharacter(name)
}

func AddCharacter(_ context.Context, app *api.App, name, description string) (api.Character, error) {
	c := api.Character{Name: name, Description: description}
	if err := app.AddCharacter(c); err != nil {
		return api.Character{}, err
	}
	return c, nil
}

func ApplyCharacterEdit(_ context.Context, app *api.App, oldName, field, raw string, normalize bool) (api.Character, EditResult, error) {
	c, err := app.GetCharacter(oldName)
	if err != nil {
		return api.Character{}, EditResult{}, err
	}
	entry, err := schema.Lookup(schema.CharacterFields(), field)
	if err != nil {
		return api.Character{}, EditResult{}, err
	}
	value := raw
	var warning string
	if normalize {
		value, warning, err = schema.Normalize(entry.Field, raw)
		if err != nil {
			return api.Character{}, EditResult{}, err
		}
	}
	if err := entry.Set(&c, value); err != nil {
		return api.Character{}, EditResult{}, err
	}
	if err := app.UpdateCharacter(oldName, c); err != nil {
		return api.Character{}, EditResult{}, err
	}
	return c, EditResult{Value: entry.Get(c), Warning: warning}, nil
}

func DeleteCharacter(_ context.Context, app *api.App, name string) error {
	return app.DeleteCharacter(name)
}

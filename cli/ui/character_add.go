package ui

import (
	"errors"

	"yukumo-script/cli/action"
	"yukumo-script/cli/prompt"
	"yukumo-script/cli/render"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"
)

// AddCharacter prompts for missing fields and creates a character.
func AddCharacter(s streams.Streams, app *api.App, name, description string, descriptionSet bool) error {
	c := api.Character{Name: name, Description: description}
	if err := prompt.EnsureString("Name", &c.Name); err != nil {
		if errors.Is(err, prompt.ErrAborted) {
			return nil
		}
		return err
	}
	if !descriptionSet {
		if err := prompt.String("Description (optional)", &c.Description); err != nil {
			if errors.Is(err, prompt.ErrAborted) {
				return nil
			}
			return err
		}
	}
	created, err := action.AddCharacter(s.Ctx, app, c.Name, c.Description)
	if err != nil {
		return err
	}
	return render.Line(s.Out, "added %q\n", created.Name)
}

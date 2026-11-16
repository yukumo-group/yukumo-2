package ui

import (
	"errors"
	"fmt"

	"yukumo-script/cli/action"
	"yukumo-script/cli/prompt"
	"yukumo-script/cli/render"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"
)

// DeleteCharacter picks a character when needed, confirms, and deletes it.
func DeleteCharacter(s streams.Streams, app *api.App, name string) error {
	name, err := PickCharacter(s, app, name)
	if err != nil {
		return err
	}
	ok := false
	if err := prompt.Confirm(fmt.Sprintf("Delete character %q?", name), &ok); err != nil {
		if errors.Is(err, prompt.ErrAborted) {
			return nil
		}
		return err
	}
	if !ok {
		return nil
	}
	if err := action.DeleteCharacter(s.Ctx, app, name); err != nil {
		return err
	}
	return render.Line(s.Out, "deleted %q\n", name)
}

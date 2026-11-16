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

// DeleteVoice picks a voice when needed, confirms, and deletes it.
func DeleteVoice(s streams.Streams, app *api.App, charName, name string) error {
	name, err := PickVoice(s, app, charName, name)
	if err != nil {
		return err
	}
	ok := false
	if err := prompt.Confirm(fmt.Sprintf("Delete voice %q?", name), &ok); err != nil {
		if errors.Is(err, prompt.ErrAborted) {
			return nil
		}
		return err
	}
	if !ok {
		return nil
	}
	if err := action.DeleteVoice(s.Ctx, app, charName, name); err != nil {
		return err
	}
	return render.Line(s.Out, "deleted voice %q\n", name)
}

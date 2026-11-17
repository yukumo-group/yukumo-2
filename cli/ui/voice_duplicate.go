package ui

import (
	"errors"
	"fmt"
	"strings"

	"yukumo-script/cli/action"
	"yukumo-script/cli/prompt"
	"yukumo-script/cli/render"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"
)

// DuplicateVoice picks a source voice, prompts for a new name, and copies it.
func DuplicateVoice(s streams.Streams, app *api.App, charName, src, newName string) error {
	src, err := PickVoice(s, app, charName, src)
	if err != nil {
		return err
	}
	if newName == "" {
		def := defaultCopyName(src)
		newName = def
		if err := prompt.String(fmt.Sprintf("New name (default %q)", def), &newName); err != nil {
			if errors.Is(err, prompt.ErrAborted) {
				return nil
			}
			return err
		}
		if newName == "" {
			newName = def
		}
	}
	final, err := action.DuplicateVoice(s.Ctx, app, charName, src, newName)
	if err != nil {
		return err
	}
	return render.Line(s.Out, "duplicated %q → %q\n", src, final)
}

func defaultCopyName(name string) string {
	return strings.TrimSpace(name) + " - copy"
}

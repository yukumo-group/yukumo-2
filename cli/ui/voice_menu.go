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

// VoiceMenu is the voice submenu under character edit.
func VoiceMenu(s streams.Streams, app *api.App, charName string) error {
	return prompt.With("voices", func() error {
		for {
			labels := []string{
				prompt.LabelBack,
				"list  —  Print voices",
				"add  —  Create a voice",
				"edit  —  Update a voice",
				"duplicate  —  Copy a voice",
				"delete  —  Remove a voice",
			}
			idx, err := prompt.Select("", "↑↓ move · enter", labels, 0)
			if err != nil {
				if errors.Is(err, prompt.ErrAborted) {
					return nil
				}
				return err
			}
			var runErr error
			switch idx {
			case 0:
				return nil
			case 1:
				runErr = prompt.With("list", func() error {
					return showVoices(s, app, charName)
				})
			case 2:
				runErr = prompt.With("add", func() error {
					return AddVoice(s, app, charName, "", 0)
				})
			case 3:
				runErr = prompt.With("edit", func() error {
					return EditVoice(s, app, charName, "")
				})
			case 4:
				runErr = prompt.With("duplicate", func() error {
					return DuplicateVoice(s, app, charName, "", "")
				})
			case 5:
				runErr = prompt.With("delete", func() error {
					return DeleteVoice(s, app, charName, "")
				})
			}
			if runErr != nil {
				_, _ = fmt.Fprintln(s.Err, runErr)
			}
		}
	})
}

func showVoices(s streams.Streams, app *api.App, charName string) error {
	voices, err := action.ListVoices(s.Ctx, app, charName)
	if err != nil {
		return err
	}
	return render.YAML(s.Out, voices)
}

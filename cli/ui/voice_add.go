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

// AddVoice prompts for a voice name, version, and params, then creates it.
func AddVoice(s streams.Streams, app *api.App, charName, name string, version int) error {
	if err := prompt.EnsureString("Voice name", &name); err != nil {
		if errors.Is(err, prompt.ErrAborted) {
			return nil
		}
		return err
	}
	if name == "" {
		return fmt.Errorf("voice name is required")
	}

	if version == 0 {
		if err := SelectVersion(&version); err != nil {
			if errors.Is(err, prompt.ErrAborted) {
				return nil
			}
			return err
		}
	}
	if version != api.Version1 && version != api.Version2 && version != api.Version10 {
		return fmt.Errorf("%w: version must be 1, 2, or 10", api.ErrInvalid)
	}

	p := api.DefaultSynthParam(version)
	err := prompt.With(fmt.Sprintf("param v%d", version), func() error {
		return EditVoiceParam(s, app, version, &p)
	})
	if err != nil {
		if errors.Is(err, prompt.ErrAborted) {
			return nil
		}
		return err
	}

	v := api.Voice{Name: name, Params: []api.SynthParam{p}}
	created, err := action.AddVoice(s.Ctx, app, charName, v)
	if err != nil {
		return err
	}
	return render.Line(s.Out, "added voice %q\n", created.Name)
}

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

// EditVoice opens the per-voice edit menu (name and version params).
func EditVoice(s streams.Streams, app *api.App, charName, voiceName string) error {
	name, err := PickVoice(s, app, charName, voiceName)
	if err != nil {
		return err
	}
	return prompt.With(name, func() error {
		current := name
		for {
			v, err := app.GetVoice(charName, current)
			if err != nil {
				return err
			}
			prompt.Replace(current)
			labels := []string{prompt.LabelBack, "name\n  " + v.Name}
			for _, ver := range []int{api.Version1, api.Version2, api.Version10} {
				_, idx := paramByVersion(&v, ver)
				status := "(absent)"
				if idx >= 0 {
					status = fmt.Sprintf("%s · speed %d",
						presetDisplay(s.Ctx, app, v.Params[idx].Version, v.Params[idx].Preset),
						v.Params[idx].Speed)
				}
				labels = append(labels, fmt.Sprintf("param v%d\n  %s", ver, status))
			}
			idx, err := prompt.Select("", "↑↓ move · enter", labels, 0)
			if err != nil {
				return err
			}
			if idx == 0 {
				return nil
			}
			if idx == 1 {
				err := prompt.With("name", func() error {
					newName := v.Name
					if err := prompt.String("", &newName); err != nil {
						return err
					}
					old := current
					v.Name = newName
					updated, err := action.SaveVoice(s.Ctx, app, charName, old, v)
					if err != nil {
						return err
					}
					current = updated.Name
					return render.Line(s.Out, "name: %s\n", updated.Name)
				})
				if err != nil {
					if errors.Is(err, prompt.ErrAborted) {
						continue
					}
					return err
				}
				continue
			}
			version := []int{api.Version1, api.Version2, api.Version10}[idx-2]
			p, pidx := paramByVersion(&v, version)
			var draft api.SynthParam
			if pidx < 0 {
				draft = api.DefaultSynthParam(version)
				p = &draft
			}
			err = prompt.With(fmt.Sprintf("param v%d", version), func() error {
				return EditVoiceParam(s, app, version, p)
			})
			if err != nil {
				if errors.Is(err, prompt.ErrAborted) {
					continue
				}
				return err
			}
			if pidx < 0 {
				v.Params = append(v.Params, *p)
			} else {
				v.Params[pidx] = *p
			}
			if _, err := action.SaveVoice(s.Ctx, app, charName, current, v); err != nil {
				return err
			}
		}
	})
}

func paramByVersion(v *api.Voice, version int) (*api.SynthParam, int) {
	for i := range v.Params {
		if v.Params[i].Version == version {
			return &v.Params[i], i
		}
	}
	return nil, -1
}

package voice

import (
	"errors"
	"fmt"
	"strconv"

	"yukumo-script/cli/action"
	"yukumo-script/cli/prompt"
	"yukumo-script/cli/render"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"
)

func runEdit(s streams.Streams, app *api.App, opts Options, charName string, args []string) error {
	name := ""
	if len(args) >= 1 {
		name = args[0]
	}
	if opts.Interactive || name == "" {
		if err := ensureVoiceName(app, charName, &name); err != nil {
			return err
		}
	}
	if name == "" {
		return fmt.Errorf("voice name is required")
	}

	if !opts.Interactive && len(args) >= 4 {
		return runOneShotVoiceEdit(s, app, charName, name, args[1:])
	}
	if opts.Interactive {
		return runInteractiveVoiceEdit(s, app, charName, name)
	}
	if len(args) >= 3 {
		return runOneShotVoiceEdit(s, app, charName, name, args[1:])
	}
	return fmt.Errorf("usage: character <name> voice edit <voice> <field> <value>")
}

func runOneShotVoiceEdit(s streams.Streams, app *api.App, charName, voiceName string, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("field and value required")
	}
	field, raw := args[0], args[1]
	switch field {
	case "name":
		_, result, err := action.ApplyVoiceName(s.Ctx, app, charName, voiceName, raw)
		if err != nil {
			return err
		}
		return render.Line(s.Out, "name: %s\n", result.Value)
	case "param", "version":
		return fmt.Errorf("use: edit <voice> <version> <field> <value>")
	default:
		if len(args) >= 3 && (field == "1" || field == "2" || field == "10") {
			ver, _ := strconv.Atoi(field)
			if _, result, err := action.ApplyVoiceParam(s.Ctx, app, charName, voiceName, ver, args[1], args[2], true); err != nil {
				return err
			} else if result.Warning != "" {
				_ = render.Warning(s.Err, result.Warning)
			}
			return render.Line(s.Out, "v%d.%s updated\n", ver, args[1])
		}
		return fmt.Errorf("unknown field %q (want name, or <version> <field> <value>)", field)
	}
}

func runInteractiveVoiceEdit(s streams.Streams, app *api.App, charName, voiceName string) error {
	return prompt.With(voiceName, func() error {
		current := voiceName
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
					name := v.Name
					if err := prompt.String("", &name); err != nil {
						return err
					}
					v.Name = name
					updated, err := action.SaveVoice(s.Ctx, app, charName, current, v)
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
				return editParamInteractive(s, app, p)
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

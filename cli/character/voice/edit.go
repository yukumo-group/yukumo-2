package voice

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"yukumo-script/cli/prompt"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func runEdit(cmd *cobra.Command, app *api.App, opts Options, charName string, args []string) error {
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
		return runOneShotVoiceEdit(cmd, app, charName, name, args[1:])
	}
	if opts.Interactive {
		return runInteractiveVoiceEdit(cmd, app, charName, name)
	}
	if len(args) >= 3 {
		return runOneShotVoiceEdit(cmd, app, charName, name, args[1:])
	}
	return fmt.Errorf("usage: character <name> voice edit <voice> <field> <value>")
}

func runOneShotVoiceEdit(cmd *cobra.Command, app *api.App, charName, voiceName string, args []string) error {
	v, err := app.GetVoice(charName, voiceName)
	if err != nil {
		return err
	}
	if len(args) < 2 {
		return fmt.Errorf("field and value required")
	}
	field, raw := args[0], args[1]
	switch field {
	case "name":
		v.Name = raw
	case "param", "version":
		return fmt.Errorf("use: edit <voice> <version> <field> <value>")
	default:
		if len(args) >= 3 && (field == "1" || field == "2" || field == "10") {
			ver, _ := strconv.Atoi(field)
			return editParamOneShot(cmd, app, charName, &v, ver, args[1], args[2])
		}
		return fmt.Errorf("unknown field %q (want name, or <version> <field> <value>)", field)
	}
	if err := app.UpdateVoice(context.Background(), charName, voiceName, v); err != nil {
		return err
	}
	_, err = fmt.Fprintf(cmd.OutOrStdout(), "name: %s\n", v.Name)
	return err
}

func editParamOneShot(
	cmd *cobra.Command,
	app *api.App,
	charName string,
	v *api.Voice,
	version int,
	field, raw string,
) error {
	p, idx := paramByVersion(v, version)
	if idx < 0 {
		np := api.DefaultSynthParam(version)
		v.Params = append(v.Params, np)
		idx = len(v.Params) - 1
		p = &v.Params[idx]
	}
	if err := applyParamField(p, field, raw, cmd.ErrOrStderr()); err != nil {
		return err
	}
	if err := app.UpdateVoice(context.Background(), charName, v.Name, *v); err != nil {
		return err
	}
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "v%d.%s updated\n", version, field)
	return err
}

func runInteractiveVoiceEdit(cmd *cobra.Command, app *api.App, charName, voiceName string) error {
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
						presetDisplay(app, v.Params[idx].Version, v.Params[idx].Preset),
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
					old := current
					v.Name = name
					if err := app.UpdateVoice(context.Background(), charName, old, v); err != nil {
						return err
					}
					current = v.Name
					_, _ = fmt.Fprintf(cmd.OutOrStdout(), "name: %s\n", v.Name)
					return nil
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
				return editParamInteractive(cmd, app, p)
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
			if err := app.UpdateVoice(context.Background(), charName, current, v); err != nil {
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

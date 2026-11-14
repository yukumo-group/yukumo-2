package voice

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"yukumo-script/internal/cli/prompt"
	"yukumo-script/internal/domain"
	svccharacter "yukumo-script/internal/service/character"

	"github.com/spf13/cobra"
)

func runEdit(cmd *cobra.Command, svc *svccharacter.Service, opts Options, charName string, args []string) error {
	name := ""
	if len(args) >= 1 {
		name = args[0]
	}
	if opts.Interactive || name == "" {
		if err := ensureVoiceName(svc, charName, &name); err != nil {
			return err
		}
	}
	if name == "" {
		return fmt.Errorf("voice name is required")
	}

	if !opts.Interactive && len(args) >= 4 {
		// voice edit <voice> <version|name|paramField> ...
		return runOneShotVoiceEdit(cmd, svc, charName, name, args[1:])
	}
	if opts.Interactive {
		return runInteractiveVoiceEdit(cmd, svc, charName, name)
	}
	if len(args) >= 3 {
		return runOneShotVoiceEdit(cmd, svc, charName, name, args[1:])
	}
	return fmt.Errorf("usage: character <name> voice edit <voice> <field> <value>")
}

func runOneShotVoiceEdit(cmd *cobra.Command, svc *svccharacter.Service, charName, voiceName string, args []string) error {
	v, err := svc.GetVoice(charName, voiceName)
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
		return fmt.Errorf("use: edit <voice> param <version> <field> <value>")
	default:
		if len(args) >= 3 && (field == "1" || field == "2" || field == "10") {
			ver, _ := strconv.Atoi(field)
			return editParamOneShot(cmd, svc, charName, &v, ver, args[1], args[2])
		}
		return fmt.Errorf("unknown field %q (want name, or <version> <field> <value>)", field)
	}
	if err := svc.UpdateVoice(context.Background(), charName, voiceName, v); err != nil {
		return err
	}
	_, err = fmt.Fprintf(cmd.OutOrStdout(), "name: %s\n", v.Name)
	return err
}

func editParamOneShot(
	cmd *cobra.Command,
	svc *svccharacter.Service,
	charName string,
	v *domain.Voice,
	version int,
	field, raw string,
) error {
	p, idx := paramByVersion(v, version)
	if idx < 0 {
		np := domain.DefaultSynthParam(version)
		v.Params = append(v.Params, np)
		idx = len(v.Params) - 1
		p = &v.Params[idx]
	}
	if err := applyParamField(svc, p, field, raw, cmd.ErrOrStderr()); err != nil {
		return err
	}
	if err := svc.UpdateVoice(context.Background(), charName, v.Name, *v); err != nil {
		return err
	}
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "v%d.%s updated\n", version, field)
	return err
}

func runInteractiveVoiceEdit(cmd *cobra.Command, svc *svccharacter.Service, charName, voiceName string) error {
	current := voiceName
	for {
		v, err := svc.GetVoice(charName, current)
		if err != nil {
			return err
		}
		labels := []string{prompt.LabelBack, "name\n  " + v.Name}
		for _, ver := range []int{domain.Version1, domain.Version2, domain.Version10} {
			_, idx := paramByVersion(&v, ver)
			status := "(absent)"
			if idx >= 0 {
				status = fmt.Sprintf("preset=%d speed=%d", v.Params[idx].Preset, v.Params[idx].Speed)
			}
			labels = append(labels, fmt.Sprintf("param v%d\n  %s", ver, status))
		}
		idx, err := prompt.Select("Edit voice", current, labels, 0)
		if err != nil {
			return err
		}
		if idx == 0 {
			return nil
		}
		if idx == 1 {
			name := v.Name
			if err := prompt.String("Voice name", &name); err != nil {
				if errors.Is(err, prompt.ErrAborted) {
					continue
				}
				return err
			}
			old := current
			v.Name = name
			if err := svc.UpdateVoice(context.Background(), charName, old, v); err != nil {
				return err
			}
			current = v.Name
			fmt.Fprintf(cmd.OutOrStdout(), "name: %s\n", v.Name)
			continue
		}
		version := []int{domain.Version1, domain.Version2, domain.Version10}[idx-2]
		p, pidx := paramByVersion(&v, version)
		if pidx < 0 {
			create := true
			if err := prompt.Confirm(fmt.Sprintf("Add version %d param?", version), &create); err != nil {
				return err
			}
			if !create {
				continue
			}
			v.Params = append(v.Params, domain.DefaultSynthParam(version))
			if err := svc.UpdateVoice(context.Background(), charName, current, v); err != nil {
				return err
			}
			v, _ = svc.GetVoice(charName, current)
			p, pidx = paramByVersion(&v, version)
		}
		if err := editParamInteractive(cmd, svc, p); err != nil {
			if errors.Is(err, prompt.ErrAborted) {
				continue
			}
			return err
		}
		v.Params[pidx] = *p
		if err := svc.UpdateVoice(context.Background(), charName, current, v); err != nil {
			return err
		}
	}
}

func paramByVersion(v *domain.Voice, version int) (*domain.SynthParam, int) {
	for i := range v.Params {
		if v.Params[i].Version == version {
			return &v.Params[i], i
		}
	}
	return nil, -1
}

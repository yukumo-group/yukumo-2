package voice

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"yukumo-script/cli/prompt"
	"yukumo-script/cli/schema"
	"yukumo-script/cli/setting"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"
)

func editParamInteractive(s streams.Streams, app *api.App, p *api.SynthParam) error {
	focus := 0
	for {
		entries := schema.VoiceParamFields(p.Version)
		metas := paramMetas(s.Ctx, app, entries, *p)
		labels := []string{prompt.LabelCancel, prompt.LabelSave}
		keys := schema.Keys(entries)
		for _, key := range keys {
			meta := metas[key]
			value := meta.Value
			if value == "" {
				value = "(empty)"
			}
			line := value
			if meta.Description != "" {
				line += "  —  " + meta.Description
			}
			labels = append(labels, key+"\n  "+line)
		}

		idx, err := prompt.Select("", "edit fields, then Save", labels, focus)
		if err != nil {
			return err
		}
		focus = idx
		switch idx {
		case 0:
			return prompt.ErrAborted
		case 1:
			p.Normalize()
			return nil
		}
		key := keys[idx-2]
		err = prompt.With(key, func() error {
			return editParamField(s, app, p, key)
		})
		if err != nil {
			if errors.Is(err, prompt.ErrAborted) {
				continue
			}
			return err
		}
	}
}

func editParamField(s streams.Streams, app *api.App, p *api.SynthParam, key string) error {
	entries := schema.VoiceParamFields(p.Version)
	entry, err := schema.Lookup(entries, key)
	if err != nil {
		return err
	}
	switch key {
	case "preset":
		return pickPreset(s, app, p)
	case "base":
		return pickBase(s, p)
	case "phont_data":
		return pickPhont(s, p)
	}
	value, err := setting.Ask(entry.Field, entry.Get(*p), s.Err)
	if err != nil {
		return err
	}
	if err := entry.Set(p, value); err != nil {
		return err
	}
	p.Normalize()
	_, _ = fmt.Fprintf(s.Out, "%s: %s\n", key, entry.Get(*p))
	return nil
}

func pickPreset(s streams.Streams, app *api.App, p *api.SynthParam) error {
	if app == nil {
		return fmt.Errorf("preset catalog unavailable")
	}
	names, err := app.PresetNames(s.Ctx, p.Version)
	if err != nil {
		return err
	}
	labels := make([]string, 0, len(names))
	indexes := make([]int, 0, len(names))
	focus := 0
	for i, name := range names {
		ok, err := app.HasPreset(s.Ctx, p.Version, i)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}
		if i == p.Preset {
			focus = len(labels)
		}
		labels = append(labels, formatPreset(p.Version, name))
		indexes = append(indexes, i)
	}
	if len(labels) == 0 {
		return fmt.Errorf("no valid presets for version %d", p.Version)
	}
	idx, err := prompt.Select("", "aq_preset_names", labels, focus)
	if err != nil {
		return err
	}
	p.Preset = indexes[idx]
	_, err = fmt.Fprintf(s.Out, "preset: %s\n", formatPreset(p.Version, names[p.Preset]))
	return err
}

func pickBase(s streams.Streams, p *api.SynthParam) error {
	focus := p.Base
	names := api.BaseNames()
	if focus < 0 || focus >= len(names) {
		focus = 0
	}
	idx, err := prompt.Select("", "AQ10 基本素片", names, focus)
	if err != nil {
		return err
	}
	p.Base = idx
	_, err = fmt.Fprintf(s.Out, "base: %s\n", api.BaseName(p.Base))
	return err
}

func pickPhont(s streams.Streams, p *api.SynthParam) error {
	for {
		status := "(none)"
		if len(p.PhontData) > 0 {
			status = fmt.Sprintf("(%d bytes)", len(p.PhontData))
		}
		labels := []string{
			prompt.LabelBack,
			"Load from .phont file",
			"Clear embedded data",
		}
		idx, err := prompt.Select("", status+"  ·  Back returns to param edit", labels, 0)
		if err != nil {
			return err
		}
		switch idx {
		case 0:
			return nil
		case 1:
			path := ""
			if err := prompt.File("", &path); err != nil {
				if errors.Is(err, prompt.ErrAborted) {
					continue
				}
				return err
			}
			if err := api.SetPhontPath(p, path); err != nil {
				_, _ = fmt.Fprintln(s.Err, err)
				continue
			}
			_, _ = fmt.Fprintf(s.Out, "phont_data: %d bytes\n", len(p.PhontData))
		case 2:
			p.PhontData = nil
			_, _ = fmt.Fprintln(s.Out, "phont_data: cleared")
		}
	}
}

func formatPreset(version int, name string) string {
	if name == "" {
		return fmt.Sprintf("v%d", version)
	}
	return fmt.Sprintf("v%d · %s", version, name)
}

func presetDisplay(ctx context.Context, app *api.App, version, preset int) string {
	if app == nil {
		return formatPreset(version, strconv.Itoa(preset))
	}
	names, err := app.PresetNames(ctx, version)
	if err != nil || preset < 0 || preset >= len(names) {
		return formatPreset(version, strconv.Itoa(preset))
	}
	return formatPreset(version, names[preset])
}

func paramMetas(ctx context.Context, app *api.App, entries []schema.Entry[api.SynthParam], p api.SynthParam) map[string]prompt.LeafMeta {
	out := setting.Metas(entries, p)
	if meta, ok := out["preset"]; ok {
		meta.Value = presetDisplay(ctx, app, p.Version, p.Preset)
		out["preset"] = meta
	}
	if meta, ok := out["base"]; ok {
		if name := api.BaseName(p.Base); name != "" {
			meta.Value = name
		}
		out["base"] = meta
	}
	return out
}

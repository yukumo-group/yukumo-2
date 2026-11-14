package voice

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"yukumo-script/internal/cli/prompt"
	"yukumo-script/internal/cli/setting"
	"yukumo-script/internal/domain"
	svcvoice "yukumo-script/internal/service/character/voice"

	"github.com/spf13/cobra"
)

func editParamInteractive(cmd *cobra.Command, svc *svcvoice.Service, p *domain.SynthParam) error {
	focus := 0
	for {
		entries := paramEntries(p.Version)
		metas := paramMetas(svc, entries, *p)
		labels := []string{prompt.LabelCancel, prompt.LabelSave}
		keys := setting.Keys(entries)
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

		idx, err := prompt.Select("Param v"+strconv.Itoa(p.Version), "edit fields, then Save", labels, focus)
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
		if err := editParamField(cmd, svc, p, key); err != nil {
			if errors.Is(err, prompt.ErrAborted) {
				continue
			}
			return err
		}
	}
}

func editParamField(cmd *cobra.Command, svc *svcvoice.Service, p *domain.SynthParam, key string) error {
	entries := paramEntries(p.Version)
	entry, err := setting.Lookup(entries, key)
	if err != nil {
		return err
	}
	switch key {
	case "preset":
		return pickPreset(cmd, svc, p)
	case "base":
		return pickBase(cmd, p)
	case "phont_data":
		return pickPhont(cmd, p)
	}
	value, err := setting.Ask(entry.Field, entry.Get(*p), cmd.ErrOrStderr())
	if err != nil {
		return err
	}
	if err := entry.Set(p, value); err != nil {
		return err
	}
	p.Normalize()
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", key, entry.Get(*p))
	return nil
}

func pickPreset(cmd *cobra.Command, svc *svcvoice.Service, p *domain.SynthParam) error {
	catalog := svc.Presets()
	if catalog == nil {
		return fmt.Errorf("preset catalog unavailable")
	}
	names, err := catalog.PresetNames(context.Background(), int32(p.Version))
	if err != nil {
		return err
	}
	labels := make([]string, 0, len(names))
	indexes := make([]int, 0, len(names))
	focus := 0
	for i, name := range names {
		ok, err := catalog.HasPreset(context.Background(), int32(p.Version), int32(i))
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
	idx, err := prompt.Select("Preset", "aq_preset_names", labels, focus)
	if err != nil {
		return err
	}
	p.Preset = indexes[idx]
	_, err = fmt.Fprintf(cmd.OutOrStdout(), "preset: %s\n", formatPreset(p.Version, names[p.Preset]))
	return err
}

func pickBase(cmd *cobra.Command, p *domain.SynthParam) error {
	focus := p.Base
	if focus < 0 || focus >= len(domain.BaseNames) {
		focus = 0
	}
	idx, err := prompt.Select("Base", "AQ10 基本素片", domain.BaseNames, focus)
	if err != nil {
		return err
	}
	p.Base = idx
	_, err = fmt.Fprintf(cmd.OutOrStdout(), "base: %s\n", domain.BaseName(p.Base))
	return err
}

func pickPhont(cmd *cobra.Command, p *domain.SynthParam) error {
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
		idx, err := prompt.Select("Phont data", status+"  ·  Back returns to param edit", labels, 0)
		if err != nil {
			return err
		}
		switch idx {
		case 0:
			return nil
		case 1:
			path := ""
			if err := prompt.File("Phont file (.phont)", &path); err != nil {
				if errors.Is(err, prompt.ErrAborted) {
					continue
				}
				return err
			}
			if err := svcvoice.SetPhontPath(p, path); err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), err)
				continue
			}
			fmt.Fprintf(cmd.OutOrStdout(), "phont_data: %d bytes\n", len(p.PhontData))
		case 2:
			p.PhontData = nil
			fmt.Fprintln(cmd.OutOrStdout(), "phont_data: cleared")
		}
	}
}

func formatPreset(version int, name string) string {
	if name == "" {
		return fmt.Sprintf("v%d", version)
	}
	return fmt.Sprintf("v%d · %s", version, name)
}

func presetDisplay(svc *svcvoice.Service, version, preset int) string {
	if svc == nil || svc.Presets() == nil {
		return formatPreset(version, strconv.Itoa(preset))
	}
	names, err := svc.Presets().PresetNames(context.Background(), int32(version))
	if err != nil || preset < 0 || preset >= len(names) {
		return formatPreset(version, strconv.Itoa(preset))
	}
	return formatPreset(version, names[preset])
}

func paramMetas(svc *svcvoice.Service, entries []setting.Entry[domain.SynthParam], p domain.SynthParam) map[string]prompt.LeafMeta {
	out := setting.Metas(entries, p)
	if meta, ok := out["preset"]; ok {
		meta.Value = presetDisplay(svc, p.Version, p.Preset)
		out["preset"] = meta
	}
	if meta, ok := out["base"]; ok {
		if name := domain.BaseName(p.Base); name != "" {
			meta.Value = name
		}
		out["base"] = meta
	}
	return out
}

func paramEntries(version int) []setting.Entry[domain.SynthParam] {
	entries := []setting.Entry[domain.SynthParam]{
		{
			Field: setting.Field{Key: "preset", Title: "Preset", Description: "Engine voice preset", Kind: setting.KindString},
			Get:   func(p domain.SynthParam) string { return strconv.Itoa(p.Preset) },
			Set: func(p *domain.SynthParam, v string) error {
				n, err := strconv.Atoi(v)
				if err != nil {
					return fmt.Errorf("invalid preset %q", v)
				}
				p.Preset = n
				return nil
			},
		},
		{
			Field: setting.Field{Key: "speed", Title: "Speed", Description: "Speech rate %", Kind: setting.KindInt, Min: 50, Max: 300},
			Get:   func(p domain.SynthParam) string { return strconv.Itoa(p.Speed) },
			Set: func(p *domain.SynthParam, v string) error {
				n, clamped, err := prompt.ParseRangedInt(v, 50, 300)
				if err != nil {
					return err
				}
				_ = clamped
				p.Speed = n
				return nil
			},
		},
	}
	if version == domain.Version2 {
		entries = append(entries, setting.Entry[domain.SynthParam]{
			Field: setting.Field{Key: "phont_data", Title: "Phont data", Description: "Embedded .phont bytes", Kind: setting.KindFile},
			Get: func(p domain.SynthParam) string {
				if len(p.PhontData) == 0 {
					return "(none)"
				}
				return fmt.Sprintf("(%d bytes)", len(p.PhontData))
			},
			Set: func(p *domain.SynthParam, v string) error {
				return svcvoice.SetPhontPath(p, v)
			},
		})
	}
	if version == domain.Version10 {
		entries = append(entries,
			setting.Entry[domain.SynthParam]{
				Field: setting.Field{Key: "base", Title: "Base", Description: "F1E/F2E/M1E", Kind: setting.KindString},
				Get: func(p domain.SynthParam) string {
					if name := domain.BaseName(p.Base); name != "" {
						return name
					}
					return strconv.Itoa(p.Base)
				},
				Set: func(p *domain.SynthParam, v string) error {
					for i, name := range domain.BaseNames {
						if strings.EqualFold(name, v) || strconv.Itoa(i) == v {
							p.Base = i
							return nil
						}
					}
					n, err := strconv.Atoi(v)
					if err != nil {
						return fmt.Errorf("invalid base %q (want F1E/F2E/M1E or 0-2)", v)
					}
					if n < 0 || n > 2 {
						return fmt.Errorf("invalid base %d", n)
					}
					p.Base = n
					return nil
				},
			},
			intEntry("volume", "Volume", "0-300", 0, 300, func(p domain.SynthParam) int { return p.Volume }, func(p *domain.SynthParam, n int) { p.Volume = n }),
			intEntry("pitch", "Pitch", "20-200", 20, 200, func(p domain.SynthParam) int { return p.Pitch }, func(p *domain.SynthParam, n int) { p.Pitch = n }),
			intEntry("accent", "Accent", "0-200", 0, 200, func(p domain.SynthParam) int { return p.Accent }, func(p *domain.SynthParam, n int) { p.Accent = n }),
			intEntry("lmd", "LMD", "Pitch1 0-200", 0, 200, func(p domain.SynthParam) int { return p.Lmd }, func(p *domain.SynthParam, n int) { p.Lmd = n }),
			intEntry("fsc", "FSC", "Sample-rate scale 50-200", 50, 200, func(p domain.SynthParam) int { return p.Fsc }, func(p *domain.SynthParam, n int) { p.Fsc = n }),
		)
	}
	return entries
}

func intEntry(
	key, title, desc string,
	min, max int,
	get func(domain.SynthParam) int,
	set func(*domain.SynthParam, int),
) setting.Entry[domain.SynthParam] {
	return setting.Entry[domain.SynthParam]{
		Field: setting.Field{Key: key, Title: title, Description: desc, Kind: setting.KindInt, Min: min, Max: max},
		Get:   func(p domain.SynthParam) string { return strconv.Itoa(get(p)) },
		Set: func(p *domain.SynthParam, v string) error {
			n, _, err := prompt.ParseRangedInt(v, min, max)
			if err != nil {
				return err
			}
			set(p, n)
			return nil
		},
	}
}

func applyParamField(p *domain.SynthParam, key, raw string, notify io.Writer) error {
	entries := paramEntries(p.Version)
	entry, err := setting.Lookup(entries, key)
	if err != nil {
		return err
	}
	value, err := setting.Normalize(entry.Field, raw, notify)
	if err != nil {
		return err
	}
	if key == "phont_data" {
		return svcvoice.SetPhontPath(p, value)
	}
	if err := entry.Set(p, value); err != nil {
		return err
	}
	p.Normalize()
	return nil
}

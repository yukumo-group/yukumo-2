package voice

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"yukumo-script/internal/cli/prompt"
	"yukumo-script/internal/cli/setting"
	"yukumo-script/internal/domain"
	svccharacter "yukumo-script/internal/service/character"

	"github.com/spf13/cobra"
)

func editParamInteractive(cmd *cobra.Command, svc *svccharacter.Service, p *domain.SynthParam) error {
	entries := paramEntries(svc, p.Version)
	return setting.BrowseKeys(setting.Keys(entries), func() map[string]prompt.LeafMeta {
		return setting.Metas(entries, *p)
	}, func(key string) error {
		entry, err := setting.Lookup(entries, key)
		if err != nil {
			return err
		}
		if key == "preset" {
			return pickPreset(cmd, svc, p)
		}
		if key == "phont_data" {
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
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", key, entry.Get(*p))
		return err
	}, "", "")
}

func pickPreset(cmd *cobra.Command, svc *svccharacter.Service, p *domain.SynthParam) error {
	catalog := svc.Presets()
	if catalog == nil {
		return fmt.Errorf("preset catalog unavailable")
	}
	names, err := catalog.PresetNames(context.Background(), int32(p.Version))
	if err != nil {
		return err
	}
	if len(names) == 0 {
		return fmt.Errorf("no presets for version %d", p.Version)
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
		labels = append(labels, fmt.Sprintf("%d  %s", i, name))
		indexes = append(indexes, i)
	}
	if len(labels) == 0 {
		return fmt.Errorf("no valid presets for version %d", p.Version)
	}
	idx, err := prompt.Select("Preset", "from aq_preset_names", labels, focus)
	if err != nil {
		return err
	}
	p.Preset = indexes[idx]
	_, err = fmt.Fprintf(cmd.OutOrStdout(), "preset: %d (%s)\n", p.Preset, names[p.Preset])
	return err
}

func pickPhont(cmd *cobra.Command, p *domain.SynthParam) error {
	path := ""
	if err := prompt.File("Phont file (.phont)", &path); err != nil {
		return err
	}
	if err := svccharacter.SetPhontPath(p, path); err != nil {
		return err
	}
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "phont_data: %d bytes\n", len(p.PhontData))
	return err
}

func paramEntries(svc *svccharacter.Service, version int) []setting.Entry[domain.SynthParam] {
	_ = svc
	entries := []setting.Entry[domain.SynthParam]{
		{
			Field: setting.Field{Key: "preset", Title: "Preset", Description: "Engine preset index", Kind: setting.KindString},
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
				return svccharacter.SetPhontPath(p, v)
			},
		})
	}
	if version == domain.Version10 {
		entries = append(entries,
			intEntry("base", "Base", "F1E/F2E/M1E", 0, 2, func(p domain.SynthParam) int { return p.Base }, func(p *domain.SynthParam, n int) { p.Base = n }),
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

func applyParamField(svc *svccharacter.Service, p *domain.SynthParam, key, raw string, notify io.Writer) error {
	entries := paramEntries(svc, p.Version)
	entry, err := setting.Lookup(entries, key)
	if err != nil {
		return err
	}
	value, err := setting.Normalize(entry.Field, raw, notify)
	if err != nil {
		return err
	}
	if key == "phont_data" {
		return svccharacter.SetPhontPath(p, value)
	}
	if err := entry.Set(p, value); err != nil {
		return err
	}
	p.Normalize()
	return nil
}

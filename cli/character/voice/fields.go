package voice

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"yukumo-script/cli/prompt"
	"yukumo-script/cli/setting"
	"yukumo-script/core/api"
)

func paramEntries(version int) []setting.Entry[api.SynthParam] {
	entries := []setting.Entry[api.SynthParam]{
		{
			Field: setting.Field{Key: "preset", Title: "Preset", Description: "Engine voice preset", Kind: setting.KindString},
			Get:   func(p api.SynthParam) string { return strconv.Itoa(p.Preset) },
			Set: func(p *api.SynthParam, v string) error {
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
			Get:   func(p api.SynthParam) string { return strconv.Itoa(p.Speed) },
			Set: func(p *api.SynthParam, v string) error {
				n, _, err := prompt.ParseRangedInt(v, 50, 300)
				if err != nil {
					return err
				}
				p.Speed = n
				return nil
			},
		},
	}
	if version == api.Version2 {
		entries = append(entries, setting.Entry[api.SynthParam]{
			Field: setting.Field{Key: "phont_data", Title: "Phont data", Description: "Embedded .phont bytes", Kind: setting.KindFile},
			Get: func(p api.SynthParam) string {
				if len(p.PhontData) == 0 {
					return "(none)"
				}
				return fmt.Sprintf("(%d bytes)", len(p.PhontData))
			},
			Set: func(p *api.SynthParam, v string) error {
				return api.SetPhontPath(p, v)
			},
		})
	}
	if version == api.Version10 {
		entries = append(entries,
			setting.Entry[api.SynthParam]{
				Field: setting.Field{Key: "base", Title: "Base", Description: "F1E/F2E/M1E", Kind: setting.KindString},
				Get: func(p api.SynthParam) string {
					if name := api.BaseName(p.Base); name != "" {
						return name
					}
					return strconv.Itoa(p.Base)
				},
				Set: func(p *api.SynthParam, v string) error {
					for i, name := range api.BaseNames() {
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
			intEntry("volume", "Volume", "0-300", 0, 300, func(p api.SynthParam) int { return p.Volume }, func(p *api.SynthParam, n int) { p.Volume = n }),
			intEntry("pitch", "Pitch", "20-200", 20, 200, func(p api.SynthParam) int { return p.Pitch }, func(p *api.SynthParam, n int) { p.Pitch = n }),
			intEntry("accent", "Accent", "0-200", 0, 200, func(p api.SynthParam) int { return p.Accent }, func(p *api.SynthParam, n int) { p.Accent = n }),
			intEntry("lmd", "LMD", "Pitch1 0-200", 0, 200, func(p api.SynthParam) int { return p.Lmd }, func(p *api.SynthParam, n int) { p.Lmd = n }),
			intEntry("fsc", "FSC", "Sample-rate scale 50-200", 50, 200, func(p api.SynthParam) int { return p.Fsc }, func(p *api.SynthParam, n int) { p.Fsc = n }),
		)
	}
	return entries
}

func intEntry(
	key, title, desc string,
	min, max int,
	get func(api.SynthParam) int,
	set func(*api.SynthParam, int),
) setting.Entry[api.SynthParam] {
	return setting.Entry[api.SynthParam]{
		Field: setting.Field{Key: key, Title: title, Description: desc, Kind: setting.KindInt, Min: min, Max: max},
		Get:   func(p api.SynthParam) string { return strconv.Itoa(get(p)) },
		Set: func(p *api.SynthParam, v string) error {
			n, _, err := prompt.ParseRangedInt(v, min, max)
			if err != nil {
				return err
			}
			set(p, n)
			return nil
		},
	}
}

func applyParamField(p *api.SynthParam, key, raw string, notify io.Writer) error {
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
		return api.SetPhontPath(p, value)
	}
	if err := entry.Set(p, value); err != nil {
		return err
	}
	p.Normalize()
	return nil
}

package schema

import (
	"fmt"
	"strconv"
	"strings"

	"yukumo-script/cli/parse"
	"yukumo-script/core/api"
)

// VoiceParamFields returns editable synth params for a voice version.
func VoiceParamFields(version int) []Entry[api.SynthParam] {
	entries := []Entry[api.SynthParam]{
		{
			Field: Field{Key: "preset", Title: "Preset", Description: "Engine voice preset", Kind: KindString},
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
	}
	if version == api.Version2 {
		entries = append(entries, Entry[api.SynthParam]{
			Field: Field{Key: "phont_data", Title: "Phont data", Description: "Embedded .phont bytes", Kind: KindFile},
			Get: func(p api.SynthParam) string {
				if len(p.PhontData) == 0 {
					return "(none)"
				}
				return fmt.Sprintf("(%d bytes)", len(p.PhontData))
			},
			Set: func(p *api.SynthParam, v string) error {
				if strings.TrimSpace(v) == "" {
					return api.SetPhontData(p, nil)
				}
				data, err := parse.ReadPhont(v)
				if err != nil {
					return err
				}
				return api.SetPhontData(p, data)
			},
		})
	}
	if version == api.Version10 {
		entries = append(entries,
			Entry[api.SynthParam]{
				Field: Field{Key: "base", Title: "Base", Description: "F1E/F2E/M1E", Kind: KindString},
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
			intEntry("pitch", "Pitch", "20-200", 20, 200, func(p api.SynthParam) int { return p.Pitch }, func(p *api.SynthParam, n int) { p.Pitch = n }),
			intEntry("accent", "Accent", "0-200", 0, 200, func(p api.SynthParam) int { return p.Accent }, func(p *api.SynthParam, n int) { p.Accent = n }),
			intEntry("lmd", "LMD", "Pitch1 0-200", 0, 200, func(p api.SynthParam) int { return p.Lmd }, func(p *api.SynthParam, n int) { p.Lmd = n }),
		)
	}
	return entries
}

func intEntry(
	key, title, desc string,
	min, max int,
	get func(api.SynthParam) int,
	set func(*api.SynthParam, int),
) Entry[api.SynthParam] {
	return Entry[api.SynthParam]{
		Field: Field{Key: key, Title: title, Description: desc, Kind: KindInt, Min: min, Max: max},
		Get:   func(p api.SynthParam) string { return strconv.Itoa(get(p)) },
		Set: func(p *api.SynthParam, v string) error {
			n, _, err := parse.ParseRangedInt(v, min, max)
			if err != nil {
				return err
			}
			set(p, n)
			return nil
		},
	}
}

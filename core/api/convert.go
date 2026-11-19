package api

import "yukumo-script/core/internal/domain"

func toCharacter(c domain.Character) Character {
	out := Character{Name: c.Name, Description: c.Description}
	if c.Voices != nil {
		out.Voices = make([]Voice, len(c.Voices))
		for i, v := range c.Voices {
			out.Voices[i] = toVoice(v)
		}
	}
	return out
}

func fromCharacter(c Character) domain.Character {
	out := domain.Character{Name: c.Name, Description: c.Description}
	if c.Voices != nil {
		out.Voices = make([]domain.Voice, len(c.Voices))
		for i, v := range c.Voices {
			out.Voices[i] = fromVoice(v)
		}
	}
	return out
}

func toVoice(v domain.Voice) Voice {
	out := Voice{Name: v.Name}
	if v.Params != nil {
		out.Params = make([]SynthParam, len(v.Params))
		for i, p := range v.Params {
			out.Params[i] = toParam(p)
		}
	}
	return out
}

func fromVoice(v Voice) domain.Voice {
	out := domain.Voice{Name: v.Name}
	if v.Params != nil {
		out.Params = make([]domain.SynthParam, len(v.Params))
		for i, p := range v.Params {
			out.Params[i] = fromParam(p)
		}
	}
	return out
}

func toParam(p domain.SynthParam) SynthParam {
	return SynthParam{
		Version:   p.Version,
		Preset:    p.Preset,
		Speed:     p.Speed,
		PhontData: copyBytes(p.PhontData),
		Base:      p.Base,
		Volume:    p.Volume,
		Pitch:     p.Pitch,
		Accent:    p.Accent,
		Lmd:       p.Lmd,
		Fsc:       p.Fsc,
	}
}

func fromParam(p SynthParam) domain.SynthParam {
	return domain.SynthParam{
		Version:   p.Version,
		Preset:    p.Preset,
		Speed:     p.Speed,
		PhontData: copyBytes(p.PhontData),
		Base:      p.Base,
		Volume:    p.Volume,
		Pitch:     p.Pitch,
		Accent:    p.Accent,
		Lmd:       p.Lmd,
		Fsc:       p.Fsc,
	}
}

func toConfig(c domain.Config) Config {
	return Config{
		Characters: CharacterConfig{InitDefaultVoice: c.Characters.InitDefaultVoice},
		Paths:      PathConfig{CharacterPath: c.Paths.CharacterPath, OutputPath: c.Paths.OutputPath},
		Synth: SynthConfig{
			OutputFormat: c.Synth.OutputFormat,
			AlwaysWrite:  c.Synth.WriteAlways(),
			ReuseCache:   c.Synth.CacheReuse(),
		},
	}
}

func fromConfig(c Config) domain.Config {
	always := c.Synth.AlwaysWrite
	reuse := c.Synth.ReuseCache
	return domain.Config{
		Characters: domain.CharacterConfig{InitDefaultVoice: c.Characters.InitDefaultVoice},
		Paths:      domain.PathConfig{CharacterPath: c.Paths.CharacterPath, OutputPath: c.Paths.OutputPath},
		Synth: domain.SynthConfig{
			OutputFormat: c.Synth.OutputFormat,
			AlwaysWrite:  &always,
			ReuseCache:   &reuse,
		},
	}
}

func copyBytes(b []byte) []byte {
	if b == nil {
		return nil
	}
	out := make([]byte, len(b))
	copy(out, b)
	return out
}

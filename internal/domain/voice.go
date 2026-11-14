package domain

const (
	Version1  = 1
	Version2  = 2
	Version10 = 10
)

// Voice is a named synthesis profile. Params holds at most one SynthParam per
// engine version (1, 2, 10).
type Voice struct {
	Name   string       `yaml:"name"`
	Params []SynthParam `yaml:"params,omitempty"`
}

// SynthParam mirrors aq_wasm.h AqSynthParam for persistence.
// PhontData is V2-only embedded .phont bytes (set from a path at edit time).
type SynthParam struct {
	Version   int    `yaml:"version"`
	Preset    int    `yaml:"preset"`
	Speed     int    `yaml:"speed"`
	PhontData []byte `yaml:"phont_data,omitempty"`
	Base      int    `yaml:"base,omitempty"`
	Volume    int    `yaml:"volume,omitempty"`
	Pitch     int    `yaml:"pitch,omitempty"`
	Accent    int    `yaml:"accent,omitempty"`
	Lmd       int    `yaml:"lmd,omitempty"`
	Fsc       int    `yaml:"fsc,omitempty"`
}

func DefaultSynthParam(version int) SynthParam {
	p := SynthParam{Version: version, Speed: 100}
	switch version {
	case Version10:
		p.Volume = 100
		p.Pitch = 100
		p.Accent = 100
		p.Lmd = 100
		p.Fsc = 100
	}
	return p
}

func (p *SynthParam) Normalize() {
	switch p.Version {
	case Version1:
		p.PhontData = nil
		p.clearV10()
	case Version2:
		p.clearV10()
	case Version10:
		p.PhontData = nil
	}
	p.Speed = clamp(p.Speed, 50, 300)
	if p.Version == Version10 {
		p.Base = clamp(p.Base, 0, 2)
		p.Volume = clamp(p.Volume, 0, 300)
		p.Pitch = clamp(p.Pitch, 20, 200)
		p.Accent = clamp(p.Accent, 0, 200)
		p.Lmd = clamp(p.Lmd, 0, 200)
		p.Fsc = clamp(p.Fsc, 50, 200)
	}
}

func (p *SynthParam) clearV10() {
	p.Base = 0
	p.Volume = 0
	p.Pitch = 0
	p.Accent = 0
	p.Lmd = 0
	p.Fsc = 0
}

func (v *Voice) Normalize() {
	for i := range v.Params {
		v.Params[i].Normalize()
	}
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

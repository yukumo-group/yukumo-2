package domain

const (
	Version1  = 1
	Version2  = 2
	Version10 = 10
)

// Voice is a named synthesis profile. Params holds at most one SynthParam per
// engine version (1, 2, 10). Speed is shared across all params on the voice.
type Voice struct {
	Name   string       `yaml:"name"`
	Speed  int          `yaml:"speed"`
	Params []SynthParam `yaml:"params,omitempty"`
}

// SynthParam mirrors aq_wasm.h AqSynthParam for persistence.
// PhontData is V2-only embedded .phont bytes (set from a path at edit time).
// Volume and Fsc are not stored; synthesis always uses 100.
type SynthParam struct {
	Version   int    `yaml:"version"`
	Preset    int    `yaml:"preset"`
	PhontData []byte `yaml:"phont_data,omitempty"`
	Base      int    `yaml:"base,omitempty"`
	Pitch     int    `yaml:"pitch,omitempty"`
	Accent    int    `yaml:"accent,omitempty"`
	Lmd       int    `yaml:"lmd,omitempty"`
}

// AQ10 base (基本素片) labels — index matches AquesTalk10 F1E/F2E/M1E.
var BaseNames = []string{"F1E", "F2E", "M1E"}

func BaseName(base int) string {
	if base >= 0 && base < len(BaseNames) {
		return BaseNames[base]
	}
	return ""
}

func DefaultSynthParam(version int) SynthParam {
	p := SynthParam{Version: version}
	switch version {
	case Version10:
		p.Pitch = 100
		p.Accent = 100
		p.Lmd = 100
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
	if p.Version == Version10 {
		p.Base = clamp(p.Base, 0, 2)
		p.Pitch = clamp(p.Pitch, 20, 200)
		p.Accent = clamp(p.Accent, 0, 200)
		p.Lmd = clamp(p.Lmd, 0, 200)
	}
}

func (p *SynthParam) clearV10() {
	p.Base = 0
	p.Pitch = 0
	p.Accent = 0
	p.Lmd = 0
}

func (v *Voice) Normalize() {
	if v.Speed == 0 {
		v.Speed = 100
	}
	v.Speed = clamp(v.Speed, 50, 300)
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

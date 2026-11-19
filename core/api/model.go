package api

import "yukumo-script/core/internal/domain"

// Character is the UI-facing character model.
type Character struct {
	Name        string  `json:"name" yaml:"name"`
	Description string  `json:"description" yaml:"description,omitempty"`
	Voices      []Voice `json:"voices" yaml:"voices,omitempty"`
}

// Voice is the UI-facing voice model.
type Voice struct {
	Name   string       `json:"name" yaml:"name"`
	Speed  int          `json:"speed" yaml:"speed"`
	Params []SynthParam `json:"params" yaml:"params,omitempty"`
}

// SynthParam is the UI-facing synthesis parameter model.
// Volume and Fsc are not stored; synthesis always uses 100.
type SynthParam struct {
	Version   int    `json:"version" yaml:"version"`
	Preset    int    `json:"preset" yaml:"preset"`
	PhontData []byte `json:"phontData,omitempty" yaml:"phont_data,omitempty"`
	Base      int    `json:"base" yaml:"base,omitempty"`
	Pitch     int    `json:"pitch" yaml:"pitch,omitempty"`
	Accent    int    `json:"accent" yaml:"accent,omitempty"`
	Lmd       int    `json:"lmd" yaml:"lmd,omitempty"`
}

// Config is the UI-facing configuration model. Booleans are effective values;
// SaveConfig writes them back as explicit values.
type Config struct {
	Characters CharacterConfig `json:"characters" yaml:"characters"`
	Paths      PathConfig      `json:"paths" yaml:"paths"`
	Synth      SynthConfig     `json:"synth" yaml:"synth"`
}

type CharacterConfig struct {
	InitDefaultVoice bool `json:"initDefaultVoice" yaml:"init_default_voice"`
}

type PathConfig struct {
	CharacterPath string `json:"characterPath" yaml:"character_path"`
	OutputPath    string `json:"outputPath" yaml:"output_path"`
}

type SynthConfig struct {
	OutputFormat string `json:"outputFormat" yaml:"output_format"`
	AlwaysWrite  bool   `json:"alwaysWrite" yaml:"always_write"`
	ReuseCache   bool   `json:"reuseCache" yaml:"reuse_cache"`
}

const (
	Version1  = domain.Version1
	Version2  = domain.Version2
	Version10 = domain.Version10
)

func DefaultSynthParam(version int) SynthParam {
	return toParam(domain.DefaultSynthParam(version))
}

func (p *SynthParam) Normalize() {
	if p == nil {
		return
	}
	d := fromParam(*p)
	d.Normalize()
	*p = toParam(d)
}

func BaseName(base int) string {
	return domain.BaseName(base)
}

func BaseNames() []string {
	return append([]string(nil), domain.BaseNames...)
}

// DefaultConfig returns non-OS defaults. Hosts supply Paths.
func DefaultConfig() Config {
	return toConfig(domain.DefaultConfig(domain.PathConfig{}))
}

// SetPhontData copies embedded phont bytes. An empty slice clears the field.
func SetPhontData(p *SynthParam, data []byte) error {
	if p == nil {
		return ErrInvalid
	}
	if len(data) == 0 {
		p.PhontData = nil
		return nil
	}
	p.PhontData = copyBytes(data)
	return nil
}

package api

import "yukumo-script/core/internal/domain"

// Character is the UI-facing character model.
type Character struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description,omitempty"`
	Voices      []Voice `yaml:"voices,omitempty"`
}

// Voice is the UI-facing voice model.
type Voice struct {
	Name   string       `yaml:"name"`
	Params []SynthParam `yaml:"params,omitempty"`
}

// SynthParam is the UI-facing synthesis parameter model.
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

// Config is the UI-facing configuration model. Booleans are effective values;
// SaveConfig writes them back as explicit values.
type Config struct {
	Characters CharacterConfig `yaml:"characters"`
	Paths      PathConfig      `yaml:"paths"`
	CLI        CLIConfig       `yaml:"cli"`
	Synth      SynthConfig     `yaml:"synth"`
}

type CharacterConfig struct {
	InitDefaultVoice bool `yaml:"init_default_voice"`
}

type PathConfig struct {
	CharacterPath string `yaml:"character_path"`
	OutputPath    string `yaml:"output_path"`
}

type CLIConfig struct {
	DefaultInteractive bool `yaml:"default_interactive"`
	WithAccessibility  bool `yaml:"with_accessibility"`
}

type SynthConfig struct {
	OutputFormat string `yaml:"output_format"`
	AlwaysWrite  bool   `yaml:"always_write"`
	ReuseCache   bool   `yaml:"reuse_cache"`
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

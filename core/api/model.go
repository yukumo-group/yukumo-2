package api

import (
	"yukumo-script/core/internal/domain"
	svcvoice "yukumo-script/core/internal/service/character/voice"
)

// Character is the UI-facing character model.
type Character struct {
	Name        string
	Description string
	Voices      []Voice
}

// Voice is the UI-facing voice model.
type Voice struct {
	Name   string
	Params []SynthParam
}

// SynthParam is the UI-facing synthesis parameter model.
type SynthParam struct {
	Version   int
	Preset    int
	Speed     int
	PhontData []byte
	Base      int
	Volume    int
	Pitch     int
	Accent    int
	Lmd       int
	Fsc       int
}

// Config is the UI-facing configuration model. Booleans are effective values;
// SaveConfig writes them back as explicit values.
type Config struct {
	Characters CharacterConfig
	Paths      PathConfig
	CLI        CLIConfig
	Synth      SynthConfig
}

type CharacterConfig struct {
	InitDefaultVoice bool
}

type PathConfig struct {
	CharacterPath string
	OutputPath    string
}

type CLIConfig struct {
	DefaultInteractive bool
	WithAccessibility  bool
}

type SynthConfig struct {
	OutputFormat string
	AlwaysWrite  bool
	ReuseCache   bool
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

func SetPhontPath(p *SynthParam, path string) error {
	if p == nil {
		return ErrInvalid
	}
	d := fromParam(*p)
	if err := svcvoice.SetPhontPath(&d, path); err != nil {
		return err
	}
	*p = toParam(d)
	return nil
}

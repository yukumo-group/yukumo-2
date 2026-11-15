package domain

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Characters CharacterConfig `yaml:"characters"`
	Paths      PathConfig      `yaml:"paths"`
	CLI        CLIConfig       `yaml:"cli"`
	Synth      SynthConfig     `yaml:"synth"`
}

const DefaultOutputFormat = "c-v-ymd"

// SynthConfig holds synthesis output preferences.
type SynthConfig struct {
	OutputFormat string `yaml:"output_format"`
	// Pointer distinguishes unset from explicit false so missing keys keep defaults.
	AlwaysWrite *bool `yaml:"always_write"`
}

func (c SynthConfig) WriteAlways() bool {
	if c.AlwaysWrite == nil {
		return true
	}
	return *c.AlwaysWrite
}

func (c *SynthConfig) SetAlwaysWrite(v bool) {
	c.AlwaysWrite = &v
}

// Pointers distinguish "unset" from explicit false so missing keys keep defaults.
type CLIConfig struct {
	DefaultInteractive *bool `yaml:"default_interactive"`
	WithAccessibility  *bool `yaml:"with_accessibility"`
}

type PathConfig struct {
	CharacterPath string `yaml:"character_path"`
	OutputPath    string `yaml:"output_path"`
}

type CharacterConfig struct {
	InitDefaultVoice bool `yaml:"init_default_voice"`
}

func DefaultCharacterPath() string {
	base, err := os.UserConfigDir()
	if err != nil {
		panic(fmt.Sprintf("resolve user config dir: %v", err))
	}
	return filepath.Join(base, "yukumo-script", "characters")
}

func DefaultOutputPath() string {
	return filepath.Join(os.TempDir(), "yukumo-script")
}

func DefaultConfig() Config {
	on := true
	off := false
	return Config{
		Paths: PathConfig{
			CharacterPath: DefaultCharacterPath(),
			OutputPath:    DefaultOutputPath(),
		},
		CLI: CLIConfig{
			DefaultInteractive: &on,
			WithAccessibility:  &off,
		},
		Synth: SynthConfig{
			OutputFormat: DefaultOutputFormat,
			AlwaysWrite:  &on,
		},
	}
}

func (c CLIConfig) Interactive() bool {
	if c.DefaultInteractive == nil {
		return true
	}
	return *c.DefaultInteractive
}

func (c CLIConfig) Accessibility() bool {
	if c.WithAccessibility == nil {
		return false
	}
	return *c.WithAccessibility
}

func (c *CLIConfig) SetInteractive(v bool) {
	c.DefaultInteractive = &v
}

func (c *CLIConfig) SetAccessibility(v bool) {
	c.WithAccessibility = &v
}

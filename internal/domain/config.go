package domain

type Config struct {
	Characters       CharacterConfig `yaml:"characters"`
	Paths            PathConfig `yaml:"paths"`
	CLI              CLIConfig  `yaml:"cli"`
}

// Pointers distinguish "unset" from explicit false so missing keys keep defaults.
type CLIConfig struct {
	DefaultInteractive *bool `yaml:"default_interactive"`
	WithAccessibility  *bool `yaml:"with_accessibility"`
}

type PathConfig struct {
	CharacterPath string `yaml:"character_path"`
	DataPath      string `yaml:"data_path"`
}

type CharacterConfig struct {
	InitDefaultVoice bool `yaml:"init_default_voice"`
}

func DefaultConfig() Config {
	on := true
	off := false
	return Config{
		CLI: CLIConfig{
			DefaultInteractive: &on,
			WithAccessibility:  &off,
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

package domain

type Config struct {
	Characters CharacterConfig `yaml:"characters"`
	Paths      PathConfig      `yaml:"paths"`
	CLI        CLIConfig       `yaml:"cli"`
	Synth      SynthConfig     `yaml:"synth"`
}

// DefaultOutputFormat uses braced tokens only.
// Month is {MM}, minute is {mm}; h/m/s/y/d/M run length = zero-padded width.
const DefaultOutputFormat = "{character}:{voice}_{yyyy}-{MM}-{dd}:{hh}-{mm}-{ss}_{ms}"

// SynthConfig holds synthesis output preferences.
type SynthConfig struct {
	OutputFormat string `yaml:"output_format"`
	// Pointers distinguish unset from explicit false so missing keys keep defaults.
	AlwaysWrite *bool `yaml:"always_write"`
	ReuseCache  *bool `yaml:"reuse_cache"`
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

func (c SynthConfig) CacheReuse() bool {
	if c.ReuseCache == nil {
		return true
	}
	return *c.ReuseCache
}

func (c *SynthConfig) SetReuseCache(v bool) {
	c.ReuseCache = &v
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

// DefaultConfig returns non-OS defaults. The platform composition root supplies
// PathConfig (filesystem paths for CLI, logical namespaces for web).
func DefaultConfig(paths PathConfig) Config {
	on := true
	off := false
	return Config{
		Paths: paths,
		CLI: CLIConfig{
			DefaultInteractive: &on,
			WithAccessibility:  &off,
		},
		Synth: SynthConfig{
			OutputFormat: DefaultOutputFormat,
			AlwaysWrite:  &on,
			ReuseCache:   &on,
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

---
description: Public DTO types for characters, voices, params, and config
icon: shapes
---

# Models

UI-facing types in `core/api`. Conversion to/from `core/internal/domain` stays inside the package.

## Character & Voice

```go
type Character struct {
    Name        string  `yaml:"name"`
    Description string  `yaml:"description,omitempty"`
    Voices      []Voice `yaml:"voices,omitempty"`
}

type Voice struct {
    Name   string       `yaml:"name"`
    Params []SynthParam `yaml:"params,omitempty"`
}
```

## SynthParam

```go
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

const (
    Version1  = /* domain.Version1 */
    Version2  = /* domain.Version2 */
    Version10 = /* domain.Version10 */
)

func DefaultSynthParam(version int) SynthParam
func (p *SynthParam) Normalize()
func SetPhontData(p *SynthParam, data []byte) error
func BaseName(base int) string
func BaseNames() []string
```

{% columns %}
{% column %}
### Shared

`preset`, `speed` (50–300)
{% endcolumn %}

{% column %}
### Version-specific

* **v2:** `PhontData`
* **v10:** `Base`, `Volume`, `Pitch`, `Accent`, `Lmd`, `Fsc`
{% endcolumn %}
{% endcolumns %}

## Config

```go
type Config struct {
    Characters CharacterConfig `yaml:"characters"`
    Paths      PathConfig      `yaml:"paths"`
    CLI        CLIConfig       `yaml:"cli"`
    Synth      SynthConfig     `yaml:"synth"`
}

func DefaultConfig() Config // non-OS defaults; hosts supply Paths
```

<details>
<summary>Nested config structs</summary>

```go
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
```
</details>

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

<table><thead><tr><th width="140">Field</th><th width="120">Go / YAML</th><th width="100">Versions</th><th width="140">Range / values</th><th>Description</th></tr></thead><tbody><tr><td>Version</td><td><code>Version</code> / <code>version</code></td><td>1, 2, 10</td><td><code>1</code> / <code>2</code> / <code>10</code></td><td>Engine version for this param slot.</td></tr><tr><td>Preset</td><td><code>Preset</code> / <code>preset</code></td><td>1, 2, 10</td><td>engine catalog</td><td>Voice preset id.</td></tr><tr><td>Speed</td><td><code>Speed</code> / <code>speed</code></td><td>1, 2, 10</td><td>50–300</td><td>Preset urgency; multiplied by <code>SpeedFactor</code> at synth time.</td></tr><tr><td>PhontData</td><td><code>PhontData</code> / <code>phont_data</code></td><td>2</td><td>embedded bytes</td><td>Custom phont payload. Use <code>SetPhontData</code> to set or clear.</td></tr><tr><td>Base</td><td><code>Base</code> / <code>base</code></td><td>10</td><td>0–2 → F1E / F2E / M1E</td><td>AquesTalk10 base voice. Labels via <code>BaseName</code> / <code>BaseNames</code>.</td></tr><tr><td>Volume</td><td><code>Volume</code> / <code>volume</code></td><td>10</td><td>0–300</td><td>Volume scale.</td></tr><tr><td>Pitch</td><td><code>Pitch</code> / <code>pitch</code></td><td>10</td><td>20–200</td><td>Pitch scale.</td></tr><tr><td>Accent</td><td><code>Accent</code> / <code>accent</code></td><td>10</td><td>0–200</td><td>Accent strength.</td></tr><tr><td>Lmd</td><td><code>Lmd</code> / <code>lmd</code></td><td>10</td><td>0–200</td><td>Length / duration modifier.</td></tr><tr><td>Fsc</td><td><code>Fsc</code> / <code>fsc</code></td><td>10</td><td>50–200</td><td>Formant shift control.</td></tr></tbody></table>

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

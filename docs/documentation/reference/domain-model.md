---
description: Character, Voice, and SynthParam models used by the CLI and API
icon: diagram-project
---

# Domain model

Shared models between YAML on disk, the CLI, and `core/api`.

{% columns %}
{% column %}
### Character

* `name` (id)
* `description`
* `voices[]`
* File: `{Sanitize(name)}.yaml`
{% endcolumn %}

{% column %}
### Voice

* `name` (id)
* `params[]`
* At most 3 params per voice
* Unique versions ∈ {1, 2, 10}
{% endcolumn %}
{% endcolumns %}

## SynthParam

Mirrors `AqSynthParam`. One param slot per engine version on a voice.

<table><thead><tr><th width="140">Field</th><th width="100">Versions</th><th width="140">Range / values</th><th>Description</th></tr></thead><tbody><tr><td><code>version</code></td><td>1, 2, 10</td><td><code>1</code> / <code>2</code> / <code>10</code></td><td>Engine version for this param slot.</td></tr><tr><td><code>preset</code></td><td>1, 2, 10</td><td>engine catalog</td><td>Voice preset id.</td></tr><tr><td><code>speed</code></td><td>1, 2, 10</td><td>50–300</td><td>Preset urgency; multiplied by CLI/API speed factor at synth time.</td></tr><tr><td><code>phont_data</code></td><td>2</td><td>embedded bytes</td><td>Custom phont payload for v2.</td></tr><tr><td><code>base</code></td><td>10</td><td>0–2 → F1E / F2E / M1E</td><td>AquesTalk10 base voice.</td></tr><tr><td><code>volume</code></td><td>10</td><td>0–300</td><td>Volume scale.</td></tr><tr><td><code>pitch</code></td><td>10</td><td>20–200</td><td>Pitch scale.</td></tr><tr><td><code>accent</code></td><td>10</td><td>0–200</td><td>Accent strength.</td></tr><tr><td><code>lmd</code></td><td>10</td><td>0–200</td><td>Length / duration modifier.</td></tr><tr><td><code>fsc</code></td><td>10</td><td>50–200</td><td>Formant shift control.</td></tr></tbody></table>

{% hint style="info" icon="sliders" %}
`Normalize` clears irrelevant fields and clamps. Effective speed: `clamp(round(voice.speed * factor), 50, 300)`; `factor≤0` treated as `1`.
{% endhint %}

<details>
<summary>Go types in core/api</summary>

See [Models](https://app.gitbook.com/s/XYoAJYQXysfu26sSYZOc/models) in the API Reference for the public DTO structs (`Character`, `Voice`, `SynthParam`, `Config`).
</details>

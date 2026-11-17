---
description: Manage voices and synthesis params under a character
icon: microphone
---

# voice

Voice commands hang under `character <name> voice …`. Each voice has at most 3 params with unique versions ∈ {1, 2, 10}.

| Verb | Behavior |
| --- | --- |
| `list` | List voices |
| `get [voice]` | Print voice |
| `add [voice]` | Create; interactive picks version and edits param |
| `edit [voice] …` | Rename or `edit <voice> <1\|2\|10> <field> <value>`; interactive edits per-version params |
| `delete [voice]` | Delete |
| `duplicate [voice] [newName]` | Deep copy; default name via `nameutil` |
| (verb omitted, interactive) | Voice menu |

## Synth parameters

<table><thead><tr><th width="140">Field</th><th width="100">Versions</th><th width="140">Range / values</th><th>Description</th></tr></thead><tbody><tr><td><code>version</code></td><td>1, 2, 10</td><td><code>1</code> / <code>2</code> / <code>10</code></td><td>Engine version for this param slot. At most one param per version on a voice.</td></tr><tr><td><code>preset</code></td><td>1, 2, 10</td><td>engine catalog</td><td>Voice preset id. Validated through <code>PresetCatalog</code> (WASM).</td></tr><tr><td><code>speed</code></td><td>1, 2, 10</td><td>50–300</td><td>Preset urgency. Effective synth speed is <code>clamp(round(speed × -s), 50, 300)</code>; factor ≤ 0 treated as 1.</td></tr><tr><td><code>phont_data</code></td><td>2</td><td>embedded bytes</td><td>V2-only custom phont. CLI edit picks a file and embeds it.</td></tr><tr><td><code>base</code></td><td>10</td><td>0–2 → F1E / F2E / M1E</td><td>AquesTalk10 base voice (基本素片).</td></tr><tr><td><code>volume</code></td><td>10</td><td>0–300</td><td>Output volume scale.</td></tr><tr><td><code>pitch</code></td><td>10</td><td>20–200</td><td>Pitch scale.</td></tr><tr><td><code>accent</code></td><td>10</td><td>0–200</td><td>Accent strength.</td></tr><tr><td><code>lmd</code></td><td>10</td><td>0–200</td><td>Length / duration modifier.</td></tr><tr><td><code>fsc</code></td><td>10</td><td>50–200</td><td>Formant shift control.</td></tr></tbody></table>

{% hint style="info" icon="shield-check" %}
`Normalize` clears fields that do not apply to the selected version and clamps in-range values. Presets are validated via WASM when an engine is available.
{% endhint %}

## Examples

```bash
just run character remilia voice list
just run character remilia voice add soft 10
just run character remilia voice edit soft 10 speed 120
just run character remilia voice duplicate soft soft-copy
```

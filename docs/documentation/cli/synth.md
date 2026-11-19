---
description: Synthesize text to WAV with preview, cache, and path overrides
icon: waveform-lines
---

# synth

Synthesize from a configured character (`-v`) or an explicit character YAML (`-P`). Shared flags cover version, speed factor, output path, preview, and cache.

| Command / flags | Behavior |
| --- | --- |
| `synth -v character:voice … TEXT` | Load character from `character_path`; `-s` default `1` |
| `synth -P file.yaml:voice … TEXT` | Character YAML override (split on **last** `:`); **mutex** with `-v` |
| Shared: `[-V 1\|2\|10] [-s factor] [-o path] [-p] [--reuse-cache[=false]]` | Version / speed / output / blocking preview playback / cache (flag overrides config) |
| `synth custom` | Stub, not implemented |

{% hint style="warning" icon="code-branch" %}
**Ref:** `character:voice` (first `:`). **`-P`:** `--voice-path` short flag; path/voice split on the **last** `:`. `-v` and `-P` are mutually exclusive.
{% endhint %}

## Version selection

Omit `-V` → first param among **1→2→10**. If none: `params empty` / `definition invalid`.

## Cache (`output_path/index.yaml`)

{% stepper %}
{% step %}
## Task hash

SHA-256 hex over text, effective speed, full `SynthParam` (incl. phont), and entire YAML bytes when using `-P`.
{% endstep %}

{% step %}
## Index map

Flat map `hash: basename.wav`. Startup drops entries whose files are missing.
{% endstep %}

{% step %}
## Hit behavior

On hit with `reuse_cache` / `--reuse-cache`: skip synth; **copy** to a new path (`-o` or a new unique formatted name) — never move the source (multithread-safe). Source kept.

Hit but file gone: remove index entry and re-synth (no error).
{% endstep %}

{% step %}
## Preview-only

No durable write (e.g. preview with `always_write=false`): do not update index; may still load WAV from cache for preview.
{% endstep %}
{% endstepper %}

{% hint style="success" %}
After each synth, duration is printed after the output path (and `(cached)` if any), e.g. `/tmp/a.wav (cached) 12.34ms`.
{% endhint %}

## Examples

```bash
just run synth -v remilia:default "こんにちは"
just run synth -v remilia:default -p "こんにちは"   # blocking playback
just run synth -v remilia:default -o /tmp/out.wav "テスト"
just run synth -P ./remilia.yaml:default "こんにちは"
just run synth -v remilia:default -V 10 -s 1.2 "こんにちは"
```

<details>
<summary>Write policy (always_write / -o / -p)</summary>

| `always_write` | `-o` | `-p` preview | Behavior |
| --- | --- | --- | --- |
| true | no | any | Write under `output_path` via format |
| true/false | yes | any | Write to `-o` |
| false | no | yes | No durable write; temp file for preview (removed after) |
| false | no | no | Still write via format |

See also [Configuration](../reference/configuration.md).
</details>

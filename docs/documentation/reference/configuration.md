---
description: config.yaml location, keys, filename patterns, and write policy
icon: sliders
---

# Configuration

{% tabs %}
{% tab title="English" %}
Config **directory**: `$YUKUMO_CONFIG_DIR` if set, else `{UserConfigDir}/yukumo-script`; file is `config.yaml` therein.

Startup `MkdirAll` on `character_path` / `output_path`; failure (including permission) **panics**. Empty path fields refill defaults. Bool defaults use pointers so missing YAML ≠ explicit `false`.
{% endtab %}

{% tab title="中文" %}
配置**目录**为 `$YUKUMO_CONFIG_DIR`（若设置）或 `{UserConfigDir}/yukumo-script`；文件为该目录下的 `config.yaml`。

启动时对 `character_path` / `output_path` 做 `MkdirAll`，失败（含无权限）**panic**；空路径字段回填默认值。布尔缺省用指针区分「未设置」与显式 `false`。
{% endtab %}
{% endtabs %}

## Keys

| Key | Default | Meaning |
| --- | --- | --- |
| `paths.character_path` | `{UserConfigDir}/yukumo-script/characters` | Character YAML directory |
| `paths.output_path` | `{TempDir}/yukumo-script` | Synth output directory |
| `characters.init_default_voice` | `false` | Auto-add `default` voice (v1) on `character add` |
| `synth.output_format` | `{character}:{voice}_{yyyy}-{MM}-{dd}:{hh}-{mm}-{ss}_{ms}` | Output basename (braced tokens only) |
| `synth.always_write` | `true` | Always write files (see write policy) |
| `synth.reuse_cache` | `true` | Reuse WAVs listed in `index.yaml` by task hash |
| `cli.default_interactive` | `true` | huh when args incomplete; root browse |
| `cli.with_accessibility` | `false` | huh accessibility |

{% hint style="warning" icon="clone" %}
`init_default_voice` **or** manual `voice add` for the same default name — not both.
{% endhint %}

## Filename pattern

Braced `{…}` only; no compact form.

{% columns %}
{% column %}
### Names

* `{character}` / `{c}`
* `{voice}` / `{v}`
{% endcolumn %}

{% column %}
### Time

Repeated letter length = zero-padded width.

* `y` year, `M` month, `d` day
* `h` / `H` hour (24h), `m` minute, `s` second
* `{ms}` milliseconds (3 digits)
{% endcolumn %}
{% endcolumns %}

{% hint style="info" %}
Default → e.g. `remilia:default_2026-11-15:13-04-05_123.wav` (`{MM}`=month, `{mm}`=minute). On collision: `name.wav` → `name-1.wav` …
{% endhint %}

## Write policy

| `always_write` | `-o` | `-p` preview | Behavior |
| --- | --- | --- | --- |
| true | no | any | Write under `output_path` via format |
| true/false | yes | any | Write to `-o` |
| false | no | yes | No durable write; temp file for preview (removed after) |
| false | no | no | Still write via format |

<a href="../cli/config.md" class="button primary">Edit via CLI</a><a href="../cli/synth.md" class="button secondary">Synth flags & cache</a>

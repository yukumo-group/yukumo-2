---
description: "Breaking: interactive CLI mode removed; the CLI is now scripted and cobra-only"
tags:
  - breaking
  - feature
---

# Remove interactive mode

The interactive CLI (huh prompts, root browse menu, `synth voice`) has been **completely removed**. Yukumo Script is now a scripted, cobra-only CLI: every command takes explicit arguments, and missing arguments exit with an error.

## Breaking changes

* Bare `yukumo-script`, `config`, and `character` print help instead of opening a menu.
* Missing arguments no longer open prompts — they error out.
* `synth voice` is removed. Use `synth -v character:voice … TEXT` or `synth -P path.yaml:voice … TEXT`.
* `config edit` requires `<key> <value>`; the interactive tree is gone.
* The `cli.default_interactive` and `cli.with_accessibility` config keys are removed, along with the whole `cli` config section.
* `core/api`: `Config.CLI` and `CLIConfig` are removed.
* `-p` preview is blocking playback.

## Command reference

| Command | Required arguments |
| --- | --- |
| `config edit` | `<key> <value>` |
| `character get` / `delete` | `<name>` |
| `character add` | `<name> [description]` |
| `character edit` | `<name> <field> <value>` |
| `character <name> voice <verb>` | verb is required: `list`, `get`, `add`, `edit`, `delete`, `duplicate` |
| `character <name> voice add` | `<voice> [1\|2\|10]` (version defaults to `1`) |

## Migration

{% stepper %}
{% step %}
## Remove `cli.*` from `config.yaml`

Delete the `cli:` section if present. Remaining keys are `paths.*`, `characters.init_default_voice`, and `synth.output_format` / `always_write` / `reuse_cache`.
{% endstep %}

{% step %}
## Pass arguments explicitly

```bash
just run config edit synth.reuse_cache false
just run character add remilia "example"
just run character remilia voice add default 1
just run synth -v remilia:default "こんにちは"
```
{% endstep %}

{% step %}
## Update Go hosts

Drop any use of `Config.CLI` / `CLIConfig` from code that calls `App.Config` or `App.SaveConfig`.
{% endstep %}
{% endstepper %}

<a href="https://app.gitbook.com/s/A67wXyFHjvGTd6AJa6xT/cli/overview" class="button primary">CLI overview</a><a href="https://app.gitbook.com/s/A67wXyFHjvGTd6AJa6xT/getting-started/quickstart" class="button secondary">Quickstart</a>

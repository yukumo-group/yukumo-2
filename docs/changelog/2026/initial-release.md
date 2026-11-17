---
description: First documented release of the Yukumo Script CLI and core/api facade
tags:
  - feature
  - docs
---

# Initial release

Yukumo Script is a CLI for character/voice management and speech synthesis via AquesTalk WASM (`third-party/aq_wasm`), plus a Go DTO facade in `core/api` for CLI hosts and future GUIs.

## CLI

{% columns %}
{% column %}
### Commands

* `config` — list / get / edit settings
* `character` — CRUD characters
* `voice` — CRUD, duplicate, param edit
* `synth` — synthesize with preview & cache
{% endcolumn %}

{% column %}
### Modes

* Interactive root browse (`just run`)
* Scripted one-shots via `just run …`
* Optional accessibility for huh prompts
{% endcolumn %}
{% endcolumns %}

<a href="https://app.gitbook.com/s/A67wXyFHjvGTd6AJa6xT/getting-started/quickstart" class="button primary">Quickstart</a><a href="https://app.gitbook.com/s/A67wXyFHjvGTd6AJa6xT/cli/overview" class="button secondary">CLI overview</a>

## Architecture

Modules depend downward only: CLI → `core/api` → services / `runtime/aq` → domain & store. Yukw packages multiple AquesTalk product WASMs into one `aq_wasm.wasm`; `aq.Pool` leases exclusive runtimes for parallel synth.

<details>
<summary>Domain & config</summary>

* Character / Voice / SynthParam with versions 1, 2, 10
* Effective speed: `clamp(round(voice.speed * factor), 50, 300)`
* `config.yaml` under `$YUKUMO_CONFIG_DIR` or `{UserConfigDir}/yukumo-script`
* Synth cache keyed by SHA-256 task hash; hits copy rather than move
</details>

## API

`core/api.App` exposes config, character, voice, and synth use cases. `cliopen.Open` wires FS stores and `aq.Pool` for local hosts.

<a href="https://app.gitbook.com/s/XYoAJYQXysfu26sSYZOc/" class="button primary">API Reference</a>

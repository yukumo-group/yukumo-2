---
description: Build requirements, just recipes, and how to run Yukumo Script
icon: download
---

# Install & build

Requires a usable `third-party/aq_wasm.wasm` (or set `YUKUMO_AQ_WASM`). Development commands use `just`. `just lint` / `just build` need golangci-lint v2 (`just install-tools`).

{% hint style="warning" icon="cube" %}
Without the WASM module, synthesis and preset validation will fail. Point `YUKUMO_AQ_WASM` at the module if it is not at the default path.
{% endhint %}

## Recipes

```bash
just build          # depguard + staticcheck, then build
just test           # all modules
just ci             # replace check + lint + build + test + vet
```

## Running the CLI

All later examples use `just run …`, which forwards args to the CLI (`go -C cli run ./cmd/yukumo-script …`). Equivalents:

```bash
just run                          # interactive root menu
just run synth -v remilia:default "こんにちは"
go run ./cli/cmd/yukumo-script …  # same without just
./yukumo-script …                 # after just build
```

{% columns %}
{% column %}
### Interactive

With no subcommand and interactive mode on, the root browse menu opens.

```bash
just run
```
{% endcolumn %}

{% column %}
### Scripted

Disable interactive prompts first when automating.

```bash
just run config edit cli.default_interactive false
```
{% endcolumn %}
{% endcolumns %}

{% hint style="success" icon="circle-check" %}
Next: [Quickstart](quickstart.md) — create a character and synthesize your first line.
{% endhint %}

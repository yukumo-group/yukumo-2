---
description: Build requirements, just recipes, and how to run Yukumo Script
icon: download
---

# Install & build

Requires a usable `third-party/aq_wasm.wasm` (or set `YUKUMO_AQ_WASM`). Development commands use `just`. `just lint` / `just build` need golangci-lint v2 (`just install-tools`).

{% hint style="warning" icon="cube" %}
Without the WASM module, synthesis and preset validation will fail. Point `YUKUMO_AQ_WASM` at the module if it is not at the default path.
{% endhint %}

## Run from Source

```bash
just build          # depguard + staticcheck, then build
```

## Running the CLI

All later examples use `just run …`, which forwards args to the CLI (`go -C cli run ./cmd/yukumo-script …`). Equivalents:

```bash
just run                          # bare invocation prints help
go run ./cli/cmd/yukumo-script …  # same without just
./yukumo-script …                 # after just build
```

The CLI is scripted (cobra-only): there are no interactive prompts. Missing arguments exit with an error, and a bare `just run` (or bare `config` / `character`) prints help.

```bash
just run config list
just run synth -v remilia:default "こんにちは"
```

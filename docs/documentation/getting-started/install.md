---
description: Build requirements, just recipes, and how to run Yukumo Script
icon: download
---

# Install & build

{% tabs %}
{% tab title="English" %}
Requires a usable `third-party/aq_wasm.wasm` (or set `YUKUMO_AQ_WASM`). Development commands use `just`. `just lint` / `just build` need golangci-lint v2 (`just install-tools`).
{% endtab %}

{% tab title="中文" %}
需可用的 `third-party/aq_wasm.wasm`（或设置 `YUKUMO_AQ_WASM`）。开发命令使用 `just`；`just lint` / `just build` 需要 golangci-lint v2（可用 `just install-tools`）。
{% endtab %}
{% endtabs %}

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

{% tabs %}
{% tab title="English" %}
All later examples use `just run …`, which forwards args to the CLI (`go -C cli run ./cmd/yukumo-script …`). Equivalents:

```bash
just run                          # interactive root menu
just run synth -v remilia:default "こんにちは"
go run ./cli/cmd/yukumo-script …  # same without just
./yukumo-script …                 # after just build
```
{% endtab %}

{% tab title="中文" %}
后续示例一律用 `just run …`（把参数转发给 CLI：`go -C cli run ./cmd/yukumo-script …`）。等价写法：

```bash
just run                          # 交互根菜单
just run synth -v remilia:default "こんにちは"
go run ./cli/cmd/yukumo-script …  # 不用 just
./yukumo-script …                 # just build 之后
```
{% endtab %}
{% endtabs %}

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

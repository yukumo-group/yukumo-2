---
description: Module layers, Yukw, multithreading, and decoupling boundaries
icon: layer-group
---

# Architecture

Modules (dependencies point downward only):

```
cli/cmd/yukumo-script (main)
  ↓
cli  ── cobra / huh / lipgloss
  ↓
core/api  ── DTO use-case facade + composition root
  ↓                    ↘
core/internal/service   runtime/aq → aq_wasm.wasm (Yukw) → AquesTalk
  ↓
core/internal/{domain, store}
```

| Module | Responsibility |
| --- | --- |
| `core/internal/domain` | Pure models, defaults, normalization |
| `core/internal/store` | Named byte IO, YAML encode/decode |
| `core/internal/service` | Use-cases; depends on store/runtime/other services via interfaces |
| `runtime/aq` | wazero runtime, `aq_*` exports, multi-goroutine pool |
| `core/api` | DTO facade + use-case API for CLI / future GUIs |
| `cli` | cobra commands, huh TUI, audio preview |

{% hint style="info" icon="gears" %}
**Yukw** is a Rust compatibility layer that links multiple AquesTalk product WASMs into one `aq_wasm.wasm`, exporting a unified `aq_*` synth/query API plus host heap helpers (`malloc`, `free`, …). This repo’s `runtime/aq` loads that module through wazero.
{% endhint %}

## Multithreading

`aq.Pool` allows parallel use across goroutines; preview playback can also run in a background goroutine.

* The WASM module is **not** thread-safe: one `Runtime` may be used by only one goroutine at a time.
* `aq.Pool` creates instances on demand: `Acquire` leases an exclusive `Runtime`, `Release` returns it; `Do` / `Synthesize` / `HasPreset` / `PresetNames` use the same lease path.
* Interactive `-p` / preview: playback runs in a background goroutine and does not block huh; a new synth stops the previous player process first.

## Decoupling

<details>
<summary>Store, service, and CLI boundaries</summary>

* `store.Store` plus `Reader` / `Writer` / `ExistenceChecker` / `Lister` / `Deleter` let CRUD swap backends. Tests use `Memory`; production uses `FS`.
* `core/internal/service/character/voice` depends only on `Characters` and `PresetCatalog` — not disk paths or a global WASM singleton.
* `core/internal/service/synth` depends on `Config`, `Voices`, and `Engine` only. Write policy stays in the service; printing paths and preview stay in the CLI (`runRequest`).
* `cli/setting` describes fields and model get/set; `cli/prompt` asks the user. Config / character / voice params share Ask, Browse, and Normalize.
* `cli/browse` mounts menus on branch commands when `cli.default_interactive=true`; leaf work remains each command’s `RunE`.
* CLI imports only `core/api`; depguard checks the remaining layer directions.
</details>

{% hint style="success" icon="book-open" %}
Public Go surface for hosts and GUIs: [API Reference](https://app.gitbook.com/s/XYoAJYQXysfu26sSYZOc/).
{% endhint %}

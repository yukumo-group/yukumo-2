# core

Domain, store, and service use-cases live under `core/internal/*`. Public Go module for those packages is `yukumo-script/core`; the CLI-facing facade is a separate module at [`api/`](api/AGENTS.md).

## Layout

| Package | Role |
|---|---|
| [`internal/domain`](internal/domain/AGENTS.md) | Pure models, defaults, normalization |
| [`internal/store`](internal/store/AGENTS.md) | Named byte IO, YAML encode/decode |
| [`internal/service`](internal/service/AGENTS.md) | Use-cases; depend on store/runtime/other services via interfaces |
| [`api`](api/AGENTS.md) | DTO facade + composition root (own module) |
| `port` | Shared port interfaces (Store, PresetCatalog, Synthesizer, …) |

## Visibility

- `internal` is unreachable from the `cli` module; external callers go through `core/api`.
- Services depend downward on domain/store/ports—not on CLI.

## Related

- Root map: [../AGENTS.md](../AGENTS.md)

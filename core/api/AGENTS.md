# core/api

DTO use-case facade + composition root for CLI and future GUIs. Module: `yukumo-script/core/api`.

## Responsibility

- Expose typed use-case APIs (`character`, `voice`, `synth`, `config`, …) over DTOs.
- Wire concrete services, stores, and `runtime/aq` at the composition root (`cliopen`, app construction).
- JSON-RPC method registry lives in `rpc/` (`Bind` + `Handle`); schemas/Dart clients from `cmd/gen`.
- Keep callers off `core/internal/*`.

## Decoupling

- CLI (and GUI) import **only** this module for core behavior.
- Implementation details stay under `core/internal` and `runtime`; this package is the stable boundary.

## Related

- [../AGENTS.md](../AGENTS.md)
- [../../cli/AGENTS.md](../../cli/AGENTS.md)
- [../../runtime/aq/AGENTS.md](../../runtime/aq/AGENTS.md)

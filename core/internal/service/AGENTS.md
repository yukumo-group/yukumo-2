# core/internal/service

Use-cases. Depend on store / runtime / other services **via interfaces**, not concrete CLI or global WASM singletons.

## Subpackages

| Package | Notes |
|---|---|
| `character` | Character CRUD; owns on-disk naming helpers used by character store wiring |
| [`character/voice`](character/voice/AGENTS.md) | Voice CRUD behind `Characters` + `PresetCatalog` |
| `config` | Config load/update |
| [`synth`](synth/AGENTS.md) | Synthesis, cache index, write policy |
| `nameutil` | Copy/default name helpers |

## Related

- [../../AGENTS.md](../../AGENTS.md)
- [../domain/AGENTS.md](../domain/AGENTS.md)
- [../store/AGENTS.md](../store/AGENTS.md)

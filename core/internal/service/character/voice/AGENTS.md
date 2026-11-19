# core/internal/service/character/voice

Voice CRUD for a character.

## Decoupling

- Depends only on:
  - **`Characters`**: load/update character
  - **`PresetCatalog`** (`port`): ask whether a preset is valid
- Does **not** depend on character disk paths or a global WASM singleton.
- Same voice logic can wire a real character service + `aq.Pool`, or in-memory fakes in tests.

## Related

- [../../AGENTS.md](../../AGENTS.md)
- [../../../../../runtime/aq/AGENTS.md](../../../../../runtime/aq/AGENTS.md)

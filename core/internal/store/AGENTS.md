# core/internal/store

Named byte IO and YAML encode/decode. Must not import `os` / `path/filepath` or service packages (FS adapter lives in `store/fs`).

## Decoupling

- `store.Store` plus `Reader` / `Writer` / `ExistenceChecker` / `Lister` / `Deleter` (aliases of `core/port`) let CRUD swap backends.
- Tests use `Memory` (no real FS); production uses `FS`.
- Services depend on these interfaces—not on a concrete filesystem path layout beyond what the wired store provides.

## Related

- [../../AGENTS.md](../../AGENTS.md)
- [../service/AGENTS.md](../service/AGENTS.md)

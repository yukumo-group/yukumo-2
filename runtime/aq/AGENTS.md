# runtime/aq

wazero runtime over `aq_wasm.wasm`, `aq_*` exports, and a multi-goroutine instance pool.

## Yukw

Yukw is a Rust compatibility layer that links multiple AquesTalk product WASMs into one `aq_wasm.wasm`, exporting a unified `aq_*` synth/query API plus host heap helpers (`malloc`, `free`, …). This package loads that module through wazero (`third-party/aq_wasm.wasm`, or `YUKUMO_AQ_WASM`).

## Multithreading

`aq.Pool` allows parallel use across goroutines.

- The WASM module is **not** thread-safe: one `Runtime` (one instance) may be used by only one goroutine at a time; never share an instance across goroutines.
- `aq.Pool` creates instances on demand: `Acquire` leases an exclusive `Runtime`, `Release` returns it; `Do` / `Synthesize` / `HasPreset` / `PresetNames` use the same lease path, so concurrent synth/queries each get their own instance.

## Related

- [../AGENTS.md](../AGENTS.md)
- [../../core/api/AGENTS.md](../../core/api/AGENTS.md)

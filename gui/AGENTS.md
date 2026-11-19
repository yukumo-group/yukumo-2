# gui

Flutter app. Talks to core only through JSON-RPC 2 (generated `lib/gen/` + `lib/bridge/`).

## Desktop

- Transport: `StdioRpcTransport` → spawn `yukumo-script serve` (NDJSON).
- Binary: `YUKUMO_RPC_BIN`, else `../yukumo-script` / `yukumo-script` on PATH.
- Smoke: `YUKUMO_RPC_BIN=../yukumo-script dart run tool/smoke_stdio.dart`

## Mobile / Web (deferred)

- `WasmRpcTransport` is a stub. Future: wasip1 `core.wasm` + `host_synth` → `aq_wasm.wasm`.
- See [docs/refactorings/4-flutter-gui.md](../docs/refactorings/4-flutter-gui.md) §9.

## Codegen

- Do not hand-edit `lib/gen/`. Run `just gen-protocol` from repo root.

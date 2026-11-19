# gui

Flutter app. Talks to core only through JSON-RPC 2 (generated `lib/gen/` + `lib/bridge/`).

## Desktop

- Transport: `StdioRpcTransport` → spawn `yukumo-script serve` (NDJSON).
- Binary: `YUKUMO_RPC_BIN`, else `../yukumo-script` / `yukumo-script` on PATH.
- Smoke: `YUKUMO_RPC_BIN=../yukumo-script dart run tool/smoke_stdio.dart`
- UI: `features/shell` — narrow floating bottom dock, wide (≥840) left `NavigationRail` (`SideNavRail`, extended ≥1280); `features/synth` home; `features/characters` list/edit (swipe clone/delete, `M3EMenuAnchor` context menu). Wide master–detail uses reusable `widgets/multi_column_view.dart` (`MultiColumnView` / `multiColumnModeOf`). Prefer `m3e_core` only — it bundles menu/dismissible/dropdown.

## Mobile / Web (deferred)

- `WasmRpcTransport` is a stub. Future: wasip1 `core.wasm` + `host_synth` → `aq_wasm.wasm`.
- See [docs/refactorings/4-flutter-gui.md](../docs/refactorings/4-flutter-gui.md) §9.

## Codegen

- Do not hand-edit `lib/gen/`. Run `just gen-protocol` from repo root.

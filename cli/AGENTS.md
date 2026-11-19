# cli

Interactive layer: cobra commands, huh TUI, audio preview. Entry: `cli/cmd/yukumo-script`.

## Responsibility

- Wire cobra command trees (`config`, `character`, `synth`, …).
- Interactive root / branch menus when `cli.default_interactive=true`.
- Print synth paths and run preview playback; do **not** own write policy (that stays in `core/internal/service/synth`).

## Decoupling

- **Import boundary**: CLI imports only `core/api`. Go’s `internal` rule makes `core/internal/*` unreachable from this module; remaining layer directions are checked by depguard.
- **Browse vs leaf work**: `enableInteractiveRoot` / browse mounts menus on branch commands when interactive is on; leaf work remains each command’s `RunE`. Command tree and interactive shell stay separate (`browse.go` + `cli/ui`, `cli/menu`).
- **setting vs prompt**: field tables and model get/set live in [`setting/`](setting/AGENTS.md); asking the user lives in [`prompt/`](prompt/AGENTS.md). Config / character / voice params share Ask, Browse, and Normalize—change widgets without touching field tables, and vice versa.
- **Preview**: interactive `-p` / preview runs playback in a background goroutine and does not block huh; a new synth stops the previous player process first.

## Related

- Parent map: [../AGENTS.md](../AGENTS.md)
- Facade: [../core/api/AGENTS.md](../core/api/AGENTS.md)

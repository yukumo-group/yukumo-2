# cli

CLI layer: cobra commands and audio preview. Entry: `cli/cmd/yukumo-script`.

## Responsibility

- Wire cobra command trees (`config`, `character`, `synth`, …).
- Print synth paths and run preview playback; do **not** own write policy (that stays in `core/internal/service/synth`).

## Decoupling

- **Import boundary**: CLI imports only `core/api`. Go’s `internal` rule makes `core/internal/*` unreachable from this module; remaining layer directions are checked by depguard.
- **Preview**: `-p` / preview runs blocking playback via the system audio player.

## Related

- Parent map: [../AGENTS.md](../AGENTS.md)
- Facade: [../core/api/AGENTS.md](../core/api/AGENTS.md)

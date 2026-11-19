<!-- gitbook-agent-instructions:start -->

## GitBook Documentation Editing

This repository contains documentation synced with GitBook via Git Sync.

Before editing GitBook-synced Markdown, YAML, or asset files, make sure the GitBook skill is available and up to date in your local agent environment. Prefer installing or updating it with:

```bash
npx skills add gitbookio/gitbook-skills
```

This command may add or update local agent skill files. Use them only as local agent instructions; do not commit those installed skill files or any tool-generated agent configuration unless the user explicitly asks for it.

If `npx` is unavailable, load the skill from:

https://gitbook.com/docs/skill.md

When making changes, preserve GitBook sync metadata such as frontmatter, `SUMMARY.md`, `gitbook-docs.yaml`, `.gitbook/`, and asset links unless the requested edit explicitly requires changing them.

<!-- gitbook-agent-instructions:end -->

## Architecture (agent map)

Modules (dependencies point downward only):

```
cli/cmd/yukumo-script (main)
  ↓
cli  ── cobra
  ↓
core/api  ── DTO use-case facade + composition root
  ↓                    ↘
core/internal/service   runtime/aq → aq_wasm.wasm (Yukw) → AquesTalk
  ↓
core/internal/{domain, store}
```

Read the package `AGENTS.md` that matches the area you are changing. Prefer the deepest matching file over this root map.

## TOC — package AGENTS.md

| Path | Read when… |
|---|---|
| [cli/AGENTS.md](cli/AGENTS.md) | CLI commands, preview, `core/api` boundary |
| [core/AGENTS.md](core/AGENTS.md) | Core module layout and `internal` visibility |
| [core/api/AGENTS.md](core/api/AGENTS.md) | DTO facade, composition root, CLI/GUI entry |
| [core/internal/domain/AGENTS.md](core/internal/domain/AGENTS.md) | Character / Voice / SynthParam model |
| [core/internal/store/AGENTS.md](core/internal/store/AGENTS.md) | Named byte IO, YAML, swappable backends |
| [core/internal/service/AGENTS.md](core/internal/service/AGENTS.md) | Use-case layer overview |
| [core/internal/service/character/voice/AGENTS.md](core/internal/service/character/voice/AGENTS.md) | Voice CRUD ports (`Characters`, `PresetCatalog`) |
| [core/internal/service/synth/AGENTS.md](core/internal/service/synth/AGENTS.md) | Synth ports, write policy vs CLI preview |
| [runtime/AGENTS.md](runtime/AGENTS.md) | Runtime module overview |
| [runtime/aq/AGENTS.md](runtime/aq/AGENTS.md) | Yukw WASM, wazero, `aq.Pool` / threading |

Human-facing build/CLI/config docs stay in [README.md](README.md) / [README.zh.md](README.zh.md).

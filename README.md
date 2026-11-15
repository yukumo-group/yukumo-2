# Yukumo Script

CLI for character/voice management and speech synthesis via AquesTalk WASM (`third-party/aq_wasm`).

[中文](README.zh.md)

## Architecture

Layers (dependencies point downward only):

```
cli → service → domain
              ↘ store
              ↘ runtime(aq) → aq_wasm.wasm (Yukw) → AquesTalk engines
```

| Layer | Responsibility |
|---|---|
| `domain` | Pure models, defaults, normalization |
| `store` | Named byte IO, YAML encode/decode |
| `runtime/aq` | wazero runtime, `aq_*` exports, multi-goroutine pool |
| `service/*` | Use-cases; depends on store/runtime/other services via interfaces |
| `cli/*` | cobra commands, huh TUI, audio preview |

**Yukw**: a Rust compatibility layer that links multiple AquesTalk product WASMs into one `aq_wasm.wasm`, exporting a unified `aq_*` synth/query API plus host heap helpers (`malloc`, `free`, …). This repo’s `runtime/aq` loads that module through wazero.

### Multithreading

Yes — `aq.Pool` allows parallel use across goroutines; preview playback can also run in a background goroutine.

- The WASM module is **not** thread-safe: one `Runtime` (one instance) may be used by only one goroutine at a time; never share an instance across goroutines.
- `aq.Pool` creates instances on demand: `Acquire` leases an exclusive `Runtime`, `Release` returns it; `Do` / `Synthesize` / `HasPreset` / `PresetNames` use the same lease path, so concurrent synth/queries each get their own instance.
- Interactive `-p` / preview: playback runs in a background goroutine and does not block huh; a new synth stops the previous player process first.

### Domain model

- `Character`: `name` (id), `description`, `voices[]`; file `{Sanitize(name)}.yaml`.
- `Voice`: `name` (id) + `params[]`; at most 3 params per voice, unique versions ∈ {1,2,10}.
- `SynthParam`: mirrors `AqSynthParam`. Shared: `preset`, `speed` (50–300). v2: `phont_data`. v10: `base` (F1E/F2E/M1E), `volume`, `pitch`, `accent`, `lmd`, `fsc`. `Normalize` clears irrelevant fields and clamps. `speed` stays on the voice as that preset’s urgency; synthesis multiplies by CLI `-s`.
- Effective speed: `clamp(round(voice.speed * factor), 50, 300)`; `factor≤0` treated as `1`.

### Decoupling

- `store.Store` plus `Reader` / `Writer` / `ExistenceChecker` / `Lister` / `Deleter` let CRUD swap backends. Tests use `Memory` (no real FS); production uses `FS`.
- `service/character/voice` depends only on `Characters` (load/update character) and `PresetCatalog` (ask whether a preset is valid)—not on character disk paths or a global WASM singleton. The same voice logic can wire a real character service + `aq.Pool`, or in-memory fakes in tests.
- `service/synth` depends on `Config`, `Voices`, and `Engine` only—who supplies config, voices, and synthesis is replaceable. Write policy stays in the service; printing paths and preview stay in the CLI (`runRequest`).
- `cli/setting` describes fields and model get/set; `cli/prompt` asks the user. Config / character / voice params share Ask, Browse, and Normalize—change widgets without touching field tables, and vice versa.
- `cli/browse` mounts menus on branch commands when `cli.default_interactive=true`; leaf work remains each command’s `RunE`.

## Build & run

Requires `third-party/aq_wasm.wasm` (or `YUKUMO_AQ_WASM`).

```bash
go build -o yukumo-script .
go test ./test/...

# Interactive root menu (cli.default_interactive defaults to true)
go run .

# Create character + voice (pick one; do not add the same default voice twice)
# A: auto default voice (v1) on character add
go run . config edit characters.init_default_voice true
go run . character add remilia "example"
# B: add voice manually
# go run . character add remilia "example"
# go run . character remilia voice add default 1
# Note: with default interactive, voice add opens the param editor; accept defaults to continue.
# For scripting first: go run . config edit cli.default_interactive false

go run . synth -v remilia:default "こんにちは"
go run . synth -v remilia:default -p "こんにちは"          # preview
go run . synth -v remilia:default -o /tmp/out.wav "テスト" # explicit output
go run . synth -P ./remilia.yaml:default "こんにちは"      # override character YAML path
```

## CLI

With no subcommand and interactive on, the root menu opens. Under interactive config/character flows, `get` is **hidden** (use edit/browse).

### `config`

| Command | Behavior |
|---|---|
| `config list` | Print all keys and values |
| `config get <key>` | Print one key (non-interactive) |
| `config edit <key> <value>` | One-shot set |
| `config edit` / `config edit <key>` | Interactive tree over `paths.*` / `characters.*` / `synth.*` / `cli.*`; bool toggles on select |

### `character`

| Command | Behavior |
|---|---|
| `character list` | JSON list |
| `character get [name]` | Print character; interactive pick if name omitted |
| `character add [name] [description]` | Create; may honor `init_default_voice` |
| `character edit [name] [field] [value]` | Edit `name`/`description`; interactive field menu includes `voices` |
| `character delete [name]` | Delete |
| `character <name> voice <verb> …` | Voice subcommands (below) |

### `character <name> voice …`

| Verb | Behavior |
|---|---|
| `list` | List voices |
| `get [voice]` | Print voice |
| `add [voice]` | Create; interactive picks version and edits param |
| `edit [voice] …` | Rename or `edit <voice> <1\|2\|10> <field> <value>`; interactive edits per-version params |
| `delete [voice]` | Delete |
| `duplicate [voice] [newName]` | Deep copy; default name via `nameutil` |
| (verb omitted, interactive) | Voice menu |

Param fields: `preset`, `speed`; v2 +`phont_data` (file → embed); v10 +`base`,`volume`,`pitch`,`accent`,`lmd`,`fsc`. Presets validated via `PresetCatalog` (WASM).

### `synth`

| Command / flags | Behavior |
|---|---|
| `synth -v character:voice … TEXT` | Load character from configured `character_path`; `-s` default `1` |
| `synth -P file.yaml:voice … TEXT` | Character YAML override (split on **last** `:`); **mutex** with `-v` |
| Shared: `[-V 1\|2\|10] [-s factor] [-o path] [-p] [--reuse-cache[=false]]` | Version / speed factor / output / preview / cache (flag overrides config; omit → `synth.reuse_cache`) |
| `synth voice` | **Interactive only**; voice tree keys are `character:voice`; **async** preview |
| `synth custom` | Stub, not implemented |

**Version**: omit `-V` → first param among **1→2→10**; else `params empty` / `definition invalid`.  
**Ref**: `character:voice` (first `:`). `-P` is the short flag for `--voice-path` (cobra allows one-letter shorts only).

**Cache (`output_path/index.yaml`)**

- Task hash (SHA-256 hex): text, effective speed, full `SynthParam` (incl. phont), and entire YAML bytes when using `-P`.
- Index is a flat map `hash: basename.wav`. Startup drops entries whose files are missing.
- On hit with `reuse_cache` / `--reuse-cache`: skip synth; **copy** to a new path (`-o` or a new unique formatted name)—never move the source (multithread-safe). Source kept.
- Hit but file gone: remove index entry and re-synth (no error).
- No durable write (e.g. preview-only with `always_write=false`): do not update index; may still load WAV from cache for preview.
- After each synth, duration is printed after the output path (and `(cached)` if any), e.g. `/tmp/a.wav (cached) 12.34ms`.

## Config (`config.yaml`)

Config **directory**: `$YUKUMO_CONFIG_DIR` if set, else `{UserConfigDir}/yukumo-script`; file is `config.yaml` therein.  
Startup `MkdirAll` on `character_path` / `output_path`; failure (including permission) **panics**. Empty path fields refill defaults.

| Key | Default | Meaning |
|---|---|---|
| `paths.character_path` | `{UserConfigDir}/yukumo-script/characters` | Character YAML directory |
| `paths.output_path` | `{TempDir}/yukumo-script` | Synth output directory |
| `characters.init_default_voice` | `false` | Auto-add `default` voice (v1) on `character add`; **either** this **or** manual `voice add`, not both for the same name |
| `synth.output_format` | `{character}:{voice}_{yyyy}-{MM}-{dd}:{hh}-{mm}-{ss}_{ms}` | Output basename (braced tokens only) |
| `synth.always_write` | `true` | Always write files (see write policy) |
| `synth.reuse_cache` | `true` | Reuse WAVs listed in `index.yaml` by task hash |
| `cli.default_interactive` | `true` | huh when args incomplete; root browse |
| `cli.with_accessibility` | `false` | huh accessibility |

Bool defaults use pointers so missing YAML ≠ explicit `false`.

**Filename pattern** (braced `{…}` only; no compact form)

- Names: `{character}`/`{c}`, `{voice}`/`{v}`
- Time: repeated letter length = zero-padded width. `y` year, `M` month, `d` day, `h`/`H` hour (24h), `m` minute, `s` second; `{ms}` milliseconds (3 digits)
- Default → e.g. `remilia:default_2026-11-15:13-04-05_123.wav` (`{MM}`=month, `{mm}`=minute)
- On collision: `name.wav` → `name-1.wav` …

**Write policy**

| `always_write` | `-o` | `-p` preview | Behavior |
|---|---|---|---|
| true | no | any | Write under `output_path` via format |
| true/false | yes | any | Write to `-o` |
| false | no | yes | No durable write; temp file for preview (removed after) |
| false | no | no | Still write via format |

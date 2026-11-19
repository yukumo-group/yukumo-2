# Yukumo Script

CLI for character/voice management and speech synthesis via AquesTalk WASM (`third-party/aq_wasm`).

[中文](README.zh.md)

## Build & run

Requires `third-party/aq_wasm.wasm` (or `YUKUMO_AQ_WASM`). Development commands use `just`; `just lint` / `just build` require golangci-lint v2 (`just install-tools`).

```bash
just build          # depguard + staticcheck, then build
just test           # all modules
just ci             # replace check + lint + build + test + vet

# Run the CLI (bare invocation prints help)
just run
# or: go run ./cli/cmd/yukumo-script

# Create character + voice (pick one; do not add the same default voice twice)
# A: auto default voice (v1) on character add
go run ./cli/cmd/yukumo-script config edit characters.init_default_voice true
go run ./cli/cmd/yukumo-script character add remilia "example"
# B: add voice manually
# go run ./cli/cmd/yukumo-script character add remilia "example"
# go run ./cli/cmd/yukumo-script character remilia voice add default 1
# voice add requires a voice name; optional version 1|2|10 (default 1)

go run ./cli/cmd/yukumo-script synth -v remilia:default "こんにちは"
go run ./cli/cmd/yukumo-script synth -v remilia:default -p "こんにちは"          # preview (blocking playback)
go run ./cli/cmd/yukumo-script synth -v remilia:default -o /tmp/out.wav "テスト" # explicit output
go run ./cli/cmd/yukumo-script synth -P ./remilia.yaml:default "こんにちは"      # override character YAML path
```

## CLI

The CLI is scripted (cobra-only): there are no interactive prompts or menus. Missing required arguments exit with an error. Bare `yukumo-script`, `config`, and `character` print help.

### `config`

| Command | Behavior |
|---|---|
| `config list` | Print all keys and values |
| `config get <key>` | Print one key |
| `config edit <key> <value>` | One-shot set (both args required) |

### `character`

| Command | Behavior |
|---|---|
| `character list` | JSON list |
| `character get <name>` | Print character |
| `character add <name> [description]` | Create; may honor `init_default_voice` |
| `character edit <name> <field> <value>` | Edit `name` / `description` |
| `character delete <name>` | Delete |
| `character <name> voice <verb> …` | Voice subcommands (below) |

### `character <name> voice …`

| Verb | Behavior |
|---|---|
| `list` | List voices |
| `get <voice>` | Print voice |
| `add <voice> [1\|2\|10]` | Create with default params for the version (default `1`) |
| `edit <voice> name <newName>` or `edit <voice> <1\|2\|10> <field> <value>` | Rename or set a per-version param field |
| `delete <voice>` | Delete |
| `duplicate <voice> [newName]` | Deep copy; default name via `nameutil` |

The verb is required; `character <name> voice` without one errors.

Voice fields: `speed` (50–300). Param fields: `preset`; v2 +`phont_data` (file → embed); v10 +`base`,`pitch`,`accent`,`lmd` (`volume`/`fsc` fixed at 100). Presets validated via `PresetCatalog` (WASM).

### `synth`

| Command / flags | Behavior |
|---|---|
| `synth -v character:voice … TEXT` | Load character from configured `character_path`; `-s` default `1` |
| `synth -P file.yaml:voice … TEXT` | Character YAML override (split on **last** `:`); **mutex** with `-v` |
| Shared: `[-V 1\|2\|10] [-s factor] [-o path] [-p] [--reuse-cache[=false]]` | Version / speed factor / output / blocking preview playback / cache (flag overrides config; omit → `synth.reuse_cache`) |
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

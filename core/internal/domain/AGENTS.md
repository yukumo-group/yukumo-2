# core/internal/domain

Pure models, defaults, and normalization. No IO, no WASM, no CLI.

## Model

- **Character**: `name` (id), `description`, `voices[]`; file `{Sanitize(name)}.yaml`.
- **Voice**: `name` (id), `speed` (50–300), `params[]`; at most 3 params per voice, unique versions ∈ {1, 2, 10}.
- **SynthParam**: mirrors `AqSynthParam`. Shared: `preset`. v2: `phont_data`. v10: `base` (F1E/F2E/M1E), `pitch`, `accent`, `lmd`. `volume`/`fsc` are not stored; synth always sends 100. `Normalize` clears irrelevant fields and clamps. Voice `speed` is the preset urgency; synthesis multiplies by CLI `-s`.
- Effective speed: `clamp(round(voice.speed * factor), 50, 300)`; `factor≤0` treated as `1`.

## Related

- [../../AGENTS.md](../../AGENTS.md)
- [../service/AGENTS.md](../service/AGENTS.md)

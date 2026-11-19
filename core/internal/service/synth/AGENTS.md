# core/internal/service/synth

Speech synthesis use-case: resolve voice, call engine, apply write/cache policy.

## Decoupling

- Depends on replaceable ports:
  - **`Config`**: who supplies config
  - **`Voices`**: who supplies voices
  - **`port.Synthesizer`**: who synthesizes
  - **`port.ArtifactStore` / `port.OutputWriter`**: cache index and durable output (optional for no-write paths)
- Write policy (whether/where to write) stays in this service.
- Printing paths and preview playback stay in the CLI (`runRequest` / play)—not here.

## Related

- [../AGENTS.md](../AGENTS.md)
- [../../../../cli/AGENTS.md](../../../../cli/AGENTS.md)
- [../../../../runtime/aq/AGENTS.md](../../../../runtime/aq/AGENTS.md)

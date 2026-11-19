# cli/setting

Describes **which fields exist** and how to get/set them on models (config, character, voice params).

## Decoupling

- Owns field metadata and model accessors—not huh widgets.
- Paired with [`../prompt`](../prompt/AGENTS.md): shared Ask / Browse / Normalize flows. Change field tables without rewriting prompts; change prompts without rewriting field tables.

## Related

- [../AGENTS.md](../AGENTS.md)
- [../prompt/AGENTS.md](../prompt/AGENTS.md)

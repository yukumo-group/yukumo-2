# cli/prompt

Owns **how to ask the user** (huh prompts, trees, nav).

## Decoupling

- Does not define business field tables; those live in [`../setting`](../setting/AGENTS.md).
- Config / character / voice param UIs reuse Ask, Browse, and Normalize so interactive controls can change independently of schema tables.

## Related

- [../AGENTS.md](../AGENTS.md)
- [../setting/AGENTS.md](../setting/AGENTS.md)

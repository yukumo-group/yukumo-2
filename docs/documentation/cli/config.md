---
description: List, get, and edit config.yaml from the CLI
icon: gear
---

# config

Manage settings under `paths.*`, `characters.*`, `synth.*`, and `cli.*`.

| Command | Behavior |
| --- | --- |
| `config list` | Print all keys and values |
| `config get <key>` | Print one key (non-interactive) |
| `config edit <key> <value>` | One-shot set |
| `config edit` / `config edit <key>` | Interactive tree; bool toggles on select |

{% stepper %}
{% step %}
## Inspect current values

```bash
just run config list
just run config get paths.character_path
```
{% endstep %}

{% step %}
## One-shot edit

```bash
just run config edit characters.init_default_voice true
just run config edit cli.default_interactive false
```
{% endstep %}

{% step %}
## Interactive tree

```bash
just run config edit
# or start at a prefix:
just run config edit synth
```
{% endstep %}
{% endstepper %}

{% hint style="success" icon="book" %}
Full key table, filename patterns, and write policy: [Configuration](../reference/configuration.md).
{% endhint %}

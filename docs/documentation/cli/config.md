---
description: List, get, and edit config.yaml from the CLI
icon: gear
---

# config

Manage settings under `paths.*`, `characters.*`, and `synth.*`. Every command is one-shot; `edit` requires both key and value.

| Command | Behavior |
| --- | --- |
| `config list` | Print all keys and values |
| `config get <key>` | Print one key |
| `config edit <key> <value>` | One-shot set |

A bare `config` prints help.

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
just run config edit synth.reuse_cache false
```
{% endstep %}
{% endstepper %}

{% hint style="success" icon="book" %}
Full key table, filename patterns, and write policy: [Configuration](../reference/configuration.md).
{% endhint %}

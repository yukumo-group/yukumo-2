---
description: Create, list, edit, and delete characters
icon: user
---

# character

Characters are stored as YAML under `paths.character_path`. Filename is `{Sanitize(name)}.yaml`.

| Command | Behavior |
| --- | --- |
| `character list` | JSON list |
| `character get [name]` | Print character; interactive pick if name omitted |
| `character add [name] [description]` | Create; may honor `init_default_voice` |
| `character edit [name] [field] [value]` | Edit `name`/`description`; interactive field menu includes `voices` |
| `character delete [name]` | Delete |
| `character <name> voice <verb> …` | Voice subcommands — see [voice](voice.md) |

{% hint style="warning" icon="clone" %}
`characters.init_default_voice` auto-adds a `default` voice (v1) on `character add`. Use **either** that **or** a manual `voice add` for the same name — not both.
{% endhint %}

## Examples

{% tabs %}
{% tab title="Add with auto voice" %}
```bash
just run config edit characters.init_default_voice true
just run character add remilia "example"
```
{% endtab %}

{% tab title="Add without voice" %}
```bash
just run character add remilia "example"
just run character remilia voice add default 1
```
{% endtab %}

{% tab title="List / get / edit" %}
```bash
just run character list
just run character get remilia
just run character edit remilia description "updated"
```
{% endtab %}
{% endtabs %}

---
description: How the scripted cobra CLI tree and command groups fit together
icon: terminal
---

# CLI overview

The CLI is scripted and cobra-only: there are no interactive prompts or menus. Commands require explicit arguments and **error out** when any are missing. A bare `yukumo-script`, `config`, or `character` prints help.

## Command groups

<table data-view="cards"><thead><tr><th></th><th></th><th></th><th data-hidden data-card-target data-type="content-ref"></th></tr></thead><tbody><tr><td><h4><i class="fa-gear" style="color:$primary;">:gear:</i></h4></td><td><h4>config</h4></td><td>List, get, and edit <code>config.yaml</code> keys.</td><td><a href="config.md">config.md</a></td></tr><tr><td><h4><i class="fa-user" style="color:$primary;">:user:</i></h4></td><td><h4>character</h4></td><td>CRUD characters and enter the voice subtree.</td><td><a href="character.md">character.md</a></td></tr><tr><td><h4><i class="fa-microphone" style="color:$primary;">:microphone:</i></h4></td><td><h4>voice</h4></td><td>Add, edit, duplicate, and delete voices under a character.</td><td><a href="voice.md">voice.md</a></td></tr><tr><td><h4><i class="fa-waveform-lines" style="color:$primary;">:waveform-lines:</i></h4></td><td><h4>synth</h4></td><td>Synthesize text to WAV with preview and cache.</td><td><a href="synth.md">synth.md</a></td></tr></tbody></table>

## Argument conventions

| Command | Required arguments |
| --- | --- |
| `config edit` | `<key> <value>` |
| `character get` / `delete` | `<name>` |
| `character add` | `<name> [description]` |
| `character edit` | `<name> <field> <value>` |
| `character <name> voice` | verb: `list`, `get`, `add`, `edit`, `delete`, or `duplicate` |
| `synth` | `-v character:voice` or `-P file.yaml:voice`, plus `TEXT` |

{% hint style="info" icon="code" %}
The CLI imports only `core/api`. Go’s `internal` rule keeps `core/internal/*` unreachable from the CLI module.
{% endhint %}

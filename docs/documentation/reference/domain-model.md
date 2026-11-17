---
description: Character, Voice, and SynthParam models used by the CLI and API
icon: diagram-project
---

# Domain model

{% tabs %}
{% tab title="English" %}
Shared models between YAML on disk, the CLI, and `core/api`.
{% endtab %}

{% tab title="中文" %}
磁盘 YAML、CLI 与 `core/api` 共用的数据模型。
{% endtab %}
{% endtabs %}

{% columns %}
{% column %}
### Character

* `name` (id)
* `description`
* `voices[]`
* File: `{Sanitize(name)}.yaml`
{% endcolumn %}

{% column %}
### Voice

* `name` (id)
* `params[]`
* At most 3 params per voice
* Unique versions ∈ {1, 2, 10}
{% endcolumn %}
{% endcolumns %}

## SynthParam

Mirrors `AqSynthParam`.

| Scope | Fields |
| --- | --- |
| Shared | `preset`, `speed` (50–300) |
| v2 | `phont_data` |
| v10 | `base` (F1E/F2E/M1E), `volume`, `pitch`, `accent`, `lmd`, `fsc` |

{% hint style="info" icon="sliders" %}
`Normalize` clears irrelevant fields and clamps. `speed` stays on the voice as that preset’s urgency; synthesis multiplies by CLI `-s`.

Effective speed: `clamp(round(voice.speed * factor), 50, 300)`; `factor≤0` treated as `1`.
{% endhint %}

<details>
<summary>Go types in core/api</summary>

See [Models](https://app.gitbook.com/s/XYoAJYQXysfu26sSYZOc/models) in the API Reference for the public DTO structs (`Character`, `Voice`, `SynthParam`, `Config`).
</details>

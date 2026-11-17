---
description: Manage voices and synthesis params under a character
icon: microphone
---

# voice

{% tabs %}
{% tab title="English" %}
Voice commands hang under `character <name> voice …`. Each voice has at most 3 params with unique versions ∈ {1, 2, 10}.
{% endtab %}

{% tab title="中文" %}
音色子命令位于 `character <name> voice …`。同一 voice 最多 3 条 param，版本唯一 ∈ {1, 2, 10}。
{% endtab %}
{% endtabs %}

| Verb | Behavior |
| --- | --- |
| `list` | List voices |
| `get [voice]` | Print voice |
| `add [voice]` | Create; interactive picks version and edits param |
| `edit [voice] …` | Rename or `edit <voice> <1\|2\|10> <field> <value>`; interactive edits per-version params |
| `delete [voice]` | Delete |
| `duplicate [voice] [newName]` | Deep copy; default name via `nameutil` |
| (verb omitted, interactive) | Voice menu |

## Param fields

{% columns %}
{% column %}
### Shared (all versions)

* `preset`
* `speed` (50–300)
{% endcolumn %}

{% column %}
### Version-specific

* **v2:** `phont_data` (file → embed)
* **v10:** `base` (F1E/F2E/M1E), `volume`, `pitch`, `accent`, `lmd`, `fsc`
{% endcolumn %}
{% endcolumns %}

{% hint style="info" icon="shield-check" %}
Presets are validated via `PresetCatalog` (WASM). `Normalize` clears irrelevant fields and clamps values per version.
{% endhint %}

## Examples

```bash
just run character remilia voice list
just run character remilia voice add soft 10
just run character remilia voice edit soft 10 speed 120
just run character remilia voice duplicate soft soft-copy
```

<details>
<summary>Effective speed at synth time</summary>

`speed` on the voice is that preset’s urgency. Synthesis multiplies by CLI `-s`:

`clamp(round(voice.speed * factor), 50, 300)` — `factor≤0` treated as `1`.
</details>

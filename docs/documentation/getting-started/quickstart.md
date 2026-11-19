---
description: Create a character, add a voice, and synthesize your first WAV
icon: rocket-launch
---

# Quickstart

Create a character with a default voice, then synthesize Japanese text. Pick **one** of the character-creation paths below — do not add the same default voice twice.

{% hint style="info" icon="terminal" %}
Commands below assume [Install & build](install.md) — use `just run …` (or the equivalents noted there).
{% endhint %}

{% stepper %}
{% step %}
## Enable auto default voice (option A)

```bash
just run config edit characters.init_default_voice true
just run character add remilia "example"
```

<details>
<summary>Option B — add the voice manually</summary>

```bash
just run character add remilia "example"
just run character remilia voice add default 1
```

{% hint style="info" %}
`voice add` requires the voice name; the version (`1`, `2`, or `10`) is optional and defaults to `1`. The voice is created with default params for that version.
{% endhint %}
</details>
{% endstep %}

{% step %}
## Synthesize

```bash
just run synth -v remilia:default "こんにちは"
```
{% endstep %}

{% step %}
## Preview or write an explicit path

{% tabs %}
{% tab title="Preview (-p)" %}
```bash
just run synth -v remilia:default -p "こんにちは"
```
{% endtab %}

{% tab title="Explicit output (-o)" %}
```bash
just run synth -v remilia:default -o /tmp/out.wav "テスト"
```
{% endtab %}

{% tab title="YAML override (-P)" %}
```bash
just run synth -P ./remilia.yaml:default "こんにちは"
```
{% endtab %}
{% endtabs %}
{% endstep %}
{% endstepper %}

{% hint style="success" icon="waveform-lines" %}
After each synth, the CLI prints the output path, optional `(cached)`, and duration — e.g. `/tmp/a.wav (cached) 12.34ms`.
{% endhint %}

## Next steps

<table data-view="cards"><thead><tr><th></th><th></th><th data-hidden data-card-target data-type="content-ref"></th></tr></thead><tbody><tr><td><h4><i class="fa-sliders" style="color:$primary;">:sliders:</i></h4></td><td>Tune <code>config.yaml</code> paths, cache, and write policy.</td><td><a href="../reference/configuration.md">configuration.md</a></td></tr><tr><td><h4><i class="fa-microphone" style="color:$primary;">:microphone:</i></h4></td><td>Edit voice params (preset, speed, v2/v10 fields).</td><td><a href="../cli/voice.md">voice.md</a></td></tr><tr><td><h4><i class="fa-wave-square" style="color:$primary;">:wave-square:</i></h4></td><td>Cache, version selection, and write policy for <code>synth</code>.</td><td><a href="../cli/synth.md">synth.md</a></td></tr></tbody></table>

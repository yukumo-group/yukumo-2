---
description: Release notes for Yukumo Script
icon: clock-rotate-left
layout:
  width: wide
  title:
    visible: true
  description:
    visible: true
  tableOfContents:
    visible: false
  outline:
    visible: true
  pagination:
    visible: false
  metadata:
    visible: false
  tags:
    visible: true
  actions:
    visible: true
  anchors:
    visible: true
---

# Changelog

{% updates format="full" %}
{% update date="2026-11-16" tags="feature,docs" %}
## Initial release

Yukumo Script ships as a CLI for character/voice management and AquesTalk WASM speech synthesis, with a Go use-case facade for hosts and future GUIs.

### Highlights

* Interactive and scripted CLI (`config`, `character`, `voice`, `synth`)
* Layered `core/api` + `core/internal/*` + `runtime/aq` architecture
* Synth cache via `output_path/index.yaml` with multithread-safe copy-on-hit
* Docs site: Documentation, API Reference, and Changelog

<a href="2026/initial-release.md" class="button primary">Read full update</a><a href="https://app.gitbook.com/s/A67wXyFHjvGTd6AJa6xT/" class="button secondary">Documentation</a>

<details>
<summary>Included</summary>

* Character / voice CRUD with per-version synth params (1 / 2 / 10)
* Preview playback (`-p`) and explicit output (`-o`)
* `cliopen.Open` composition root for FS + `aq.Pool`
* Config keys for paths, cache, write policy, and interactive defaults
</details>
{% endupdate %}
{% endupdates %}

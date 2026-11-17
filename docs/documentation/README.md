---
description: CLI docs for character and voice management, plus AquesTalk speech synthesis
icon: book-open
---

# Documentation

**Yukumo Script** is a CLI for character/voice management and speech synthesis via AquesTalk WASM (`third-party/aq_wasm`).

Use these docs to install the tool, create characters and voices, synthesize audio, and configure paths and cache behavior.

<button type="button" class="button primary" data-action="ask" data-icon="gitbook-assistant">Ask a question…</button>

<button type="button" class="button secondary" data-action="ask" data-query="How do I synthesize speech with the CLI?" data-icon="waveform-lines">Synthesize speech</button><button type="button" class="button secondary" data-action="ask" data-query="How do I create a character and voice?" data-icon="user-plus">Create a character</button><button type="button" class="button secondary" data-action="ask" data-query="Where is config.yaml and what keys exist?" data-icon="sliders">Configure the CLI</button>

***

{% hint style="info" icon="terminal" %}
Prefer embedding the Go API instead of the CLI? See the [API Reference](https://app.gitbook.com/s/XYoAJYQXysfu26sSYZOc/).
{% endhint %}

## Where to start

<table data-card-size="large" data-view="cards"><thead><tr><th></th><th></th><th></th><th data-hidden data-card-target data-type="content-ref"></th></tr></thead><tbody><tr><td><h4><i class="fa-download" style="color:$primary;">:download:</i></h4></td><td><h4>Install &amp; build</h4></td><td>Requirements, <code>just</code> recipes, and how to run the binary.</td><td><a href="getting-started/install.md">install.md</a></td></tr><tr><td><h4><i class="fa-rocket-launch" style="color:$primary;">:rocket-launch:</i></h4></td><td><h4>Quickstart</h4></td><td>Create a character, add a voice, and synthesize your first WAV.</td><td><a href="getting-started/quickstart.md">quickstart.md</a></td></tr><tr><td><h4><i class="fa-terminal" style="color:$primary;">:terminal:</i></h4></td><td><h4>CLI commands</h4></td><td><code>config</code>, <code>character</code>, <code>voice</code>, and <code>synth</code> reference.</td><td><a href="cli/overview.md">overview.md</a></td></tr><tr><td><h4><i class="fa-gear" style="color:$primary;">:gear:</i></h4></td><td><h4>Configuration</h4></td><td><code>config.yaml</code> keys, filename patterns, and write policy.</td><td><a href="reference/configuration.md">configuration.md</a></td></tr></tbody></table>

## Popular tasks

<table data-view="cards"><thead><tr><th></th><th data-type="content-ref"></th><th data-type="content-ref"></th><th data-type="content-ref"></th></tr></thead><tbody><tr><td><h4>Day-to-day CLI</h4></td><td><a href="cli/synth.md">synth.md</a></td><td><a href="cli/character.md">character.md</a></td><td><a href="cli/voice.md">voice.md</a></td></tr><tr><td><h4>Setup</h4></td><td><a href="cli/config.md">config.md</a></td><td><a href="reference/configuration.md">configuration.md</a></td><td><a href="getting-started/install.md">install.md</a></td></tr><tr><td><h4>Internals</h4></td><td><a href="reference/architecture.md">architecture.md</a></td><td><a href="reference/domain-model.md">domain-model.md</a></td><td><a href="https://app.gitbook.com/s/XYoAJYQXysfu26sSYZOc/">API Reference</a></td></tr></tbody></table>

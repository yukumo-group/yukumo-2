---
description: Go use-case facade in core/api for CLI hosts and future GUIs
icon: code
---

# API Reference

`core/api` is the DTO use-case facade. The CLI (and future GUIs) call this package only — `core/internal/*` stays behind Go’s `internal` boundary.

{% hint style="info" icon="box" %}
Import path: `yukumo-script/core/api`. For a filesystem + AquesTalk composition root, use `yukumo-script/core/api/cliopen`.
{% endhint %}

## Surface map

<table data-card-size="large" data-view="cards"><thead><tr><th></th><th></th><th></th><th data-hidden data-card-target data-type="content-ref"></th></tr></thead><tbody><tr><td><h4><i class="fa-cube" style="color:$primary;">:cube:</i></h4></td><td><h4>App &amp; deps</h4></td><td><code>New</code>, <code>Deps</code>, <code>Close</code>, <code>CleanIndex</code>.</td><td><a href="app.md">app.md</a></td></tr><tr><td><h4><i class="fa-shapes" style="color:$primary;">:shapes:</i></h4></td><td><h4>Models</h4></td><td><code>Character</code>, <code>Voice</code>, <code>SynthParam</code>, <code>Config</code>.</td><td><a href="models.md">models.md</a></td></tr><tr><td><h4><i class="fa-sliders" style="color:$primary;">:sliders:</i></h4></td><td><h4>Config</h4></td><td><code>Config</code> / <code>SaveConfig</code>.</td><td><a href="config.md">config.md</a></td></tr><tr><td><h4><i class="fa-users" style="color:$primary;">:users:</i></h4></td><td><h4>Characters &amp; voices</h4></td><td>CRUD plus presets.</td><td><a href="characters.md">characters.md</a></td></tr><tr><td><h4><i class="fa-waveform-lines" style="color:$primary;">:waveform-lines:</i></h4></td><td><h4>Synth</h4></td><td><code>SynthRequest</code> → <code>SynthResult</code>.</td><td><a href="synth.md">synth.md</a></td></tr><tr><td><h4><i class="fa-folder-tree" style="color:$primary;">:folder-tree:</i></h4></td><td><h4>cliopen</h4></td><td>FS stores, OS paths, <code>aq.Pool</code> wiring.</td><td><a href="cliopen.md">cliopen.md</a></td></tr></tbody></table>

{% columns %}
{% column %}
### Host-supplied ports

`Deps` requires `ConfigStore` and `CharacterStore`. Optional: `Artifacts`, `Outputs`, `Engine`.
{% endcolumn %}

{% column %}
### CLI docs

Command-line usage lives in [Documentation](https://app.gitbook.com/s/A67wXyFHjvGTd6AJa6xT/).
{% endcolumn %}
{% endcolumns %}

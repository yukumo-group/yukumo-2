---
description: Voice CRUD, duplicate, and preset catalog queries
icon: microphone
---

# Voices

```go
func (a *App) ListVoices(charName string) ([]Voice, error)
func (a *App) GetVoice(charName, voiceName string) (Voice, error)
func (a *App) AddVoice(ctx context.Context, charName string, v Voice) error
func (a *App) UpdateVoice(ctx context.Context, charName, oldName string, v Voice) error
func (a *App) DeleteVoice(charName, voiceName string) error
func (a *App) DuplicateVoice(ctx context.Context, charName, srcName, newName string) (string, error)

func (a *App) PresetNames(ctx context.Context, version int) ([]string, error)
func (a *App) HasPreset(ctx context.Context, version, preset int) (bool, error)
```

{% columns %}
{% column %}
### CRUD

Voices are always scoped to a character name. `DuplicateVoice` returns the resolved new name.
{% endcolumn %}

{% column %}
### Presets

`PresetNames` / `HasPreset` require an engine that exposes a preset catalog. Without one → `ErrInvalid` (“preset catalog unavailable”).
{% endcolumn %}
{% endcolumns %}

{% hint style="info" icon="shield-check" %}
Add/update validate presets through the engine when available — same `PresetCatalog` path as the CLI.
{% endhint %}

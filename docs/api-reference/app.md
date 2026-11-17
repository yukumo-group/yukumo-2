---
description: Construct and shut down the use-case facade
icon: cube
---

# App & deps

```go
type Deps struct {
    ConfigStore    port.Store    // required
    CharacterStore port.Store    // required
    Artifacts      port.ArtifactStore
    Outputs        port.OutputWriter
    Engine         port.Engine
    Defaults       Config
}

type App struct { /* unexported services */ }

func New(deps Deps) (*App, error)
func (a *App) CleanIndex() error
func (a *App) Close(ctx context.Context) error
```

{% stepper %}
{% step %}
## Provide stores and defaults

`New` performs **no OS IO**. Hosts supply stores, optional engine/artifacts/outputs, and `Defaults`.
{% endstep %}

{% step %}
## Call use cases

Character, voice, config, and synth methods hang off `*App`.
{% endstep %}

{% step %}
## Close

`Close` closes the engine (if `port.Closer`) and flushes stores that implement `port.Flusher`.
{% endstep %}
{% endstepper %}

{% hint style="warning" %}
Missing `ConfigStore` or `CharacterStore` → `New` returns an error.
{% endhint %}

{% hint style="info" %}
`CleanIndex` drops synth cache entries whose artifact file is missing. `cliopen.Open` calls this after wiring.
{% endhint %}

---
description: Filesystem composition root used by the CLI
icon: folder-tree
---

# cliopen

Package `yukumo-script/core/api/cliopen` opens a fully wired `*api.App` for local use: FS stores, OS paths, and `aq.Pool`.

```go
type Options struct {
    ConfigDir     string
    CharacterPath string // optional override
    OutputPath    string // optional override
}

func Open(opts Options) (*api.App, error)
func EnsurePaths(app *api.App) error
```

{% stepper %}
{% step %}
## Resolve config dir

`Options.ConfigDir`, else `$YUKUMO_CONFIG_DIR`, else `{UserConfigDir}/yukumo-script`.
{% endstep %}

{% step %}
## Wire stores + engine

FS config/character/artifact stores, `fs.Output`, and `aq.NewPool()` as the engine.
{% endstep %}

{% step %}
## Ensure paths & clean index

`MkdirAll` on character/output paths, then `CleanIndex`. Failures close the app/pool and return.
{% endstep %}
{% endstepper %}

{% hint style="success" %}
```go
app, err := cliopen.Open(cliopen.Options{})
if err != nil {
    return err
}
defer app.Close(context.Background())
```
{% endhint %}

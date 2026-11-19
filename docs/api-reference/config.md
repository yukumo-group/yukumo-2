---
description: Read and persist Config through App
icon: sliders
---

# Config

```go
func (a *App) Config() (Config, error)
func (a *App) SaveConfig(c Config) error
```

{% hint style="info" %}
Booleans on `Config` are **effective** values. `SaveConfig` writes them back as explicit values.
{% endhint %}

{% tabs %}
{% tab title="Read" %}
```go
cfg, err := app.Config()
if err != nil {
    return err
}
fmt.Println(cfg.Paths.CharacterPath)
```
{% endtab %}

{% tab title="Save" %}
```go
cfg, err := app.Config()
if err != nil {
    return err
}
cfg.Synth.ReuseCache = false
return app.SaveConfig(cfg)
```
{% endtab %}
{% endtabs %}

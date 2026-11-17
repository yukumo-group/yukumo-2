---
description: Character CRUD on App
icon: user
---

# Characters

```go
func (a *App) ListCharacters() ([]Character, error)
func (a *App) GetCharacter(name string) (Character, error)
func (a *App) AddCharacter(c Character) error
func (a *App) UpdateCharacter(oldName string, c Character) error
func (a *App) DeleteCharacter(name string) error
```

{% stepper %}
{% step %}
## Add

```go
err := app.AddCharacter(api.Character{
    Name:        "remilia",
    Description: "example",
})
```
{% endstep %}

{% step %}
## Get / list

```go
c, err := app.GetCharacter("remilia")
list, err := app.ListCharacters()
```
{% endstep %}

{% step %}
## Update / delete

```go
err := app.UpdateCharacter("remilia", c)
err = app.DeleteCharacter("remilia")
```
{% endstep %}
{% endstepper %}

{% hint style="warning" %}
Whether `AddCharacter` auto-creates a default voice depends on `characters.init_default_voice` in config — same rule as the CLI.
{% endhint %}

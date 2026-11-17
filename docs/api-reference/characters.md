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

### Add

```go
err := app.AddCharacter(api.Character{
    Name:        "remilia",
    Description: "example",
})
```

### Get / list

```go
c, err := app.GetCharacter("remilia")
list, err := app.ListCharacters()
```

### Update / delete

```go
err := app.UpdateCharacter("remilia", c)
err = app.DeleteCharacter("remilia")
```

{% hint style="warning" %}
Whether `AddCharacter` auto-creates a default voice depends on `characters.init_default_voice` in config — same rule as the CLI.
{% endhint %}

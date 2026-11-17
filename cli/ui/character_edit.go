package ui

import (
	"errors"

	"yukumo-script/cli/action"
	"yukumo-script/cli/prompt"
	"yukumo-script/cli/render"
	"yukumo-script/cli/schema"
	"yukumo-script/cli/setting"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"
)

// EditCharacter edits one field or opens the character field menu.
func EditCharacter(s streams.Streams, app *api.App, partialName, partialField, _ string) error {
	name, err := PickCharacter(s, app, partialName)
	if err != nil {
		return err
	}

	if partialField != "" {
		c, err := app.GetCharacter(name)
		if err != nil {
			return err
		}
		entry, err := schema.Lookup(schema.CharacterFields(), partialField)
		if err != nil {
			return err
		}
		var updatedName string
		err = prompt.With(name, func() error {
			return prompt.With(partialField, func() error {
				value, err := setting.Ask(entry.Field, entry.Get(c), s.Err)
				if err != nil {
					return err
				}
				var applyErr error
				updatedName, applyErr = applyCharacterEdit(s, app, name, partialField, value, false)
				return applyErr
			})
		})
		if err != nil {
			return err
		}
		name = updatedName
	}

	return prompt.With(name, func() error {
		currentName := name
		for {
			c, err := app.GetCharacter(currentName)
			if err != nil {
				return err
			}
			prompt.Replace(currentName)
			metas := setting.Metas(schema.CharacterFields(), c)
			labels := []string{prompt.LabelBack}
			keys := schema.Keys(schema.CharacterFields())
			for _, key := range keys {
				labels = append(labels, leafLabel(key, metas[key]))
			}
			labels = append(labels, "voices/\n  Edit voice profiles")

			idx, err := prompt.Select("", "↑↓ move · enter", labels, 0)
			if err != nil {
				if errors.Is(err, prompt.ErrAborted) {
					return nil
				}
				return err
			}
			if idx == 0 {
				return nil
			}
			if idx == len(labels)-1 {
				if err := VoiceMenu(s, app, currentName); err != nil {
					return err
				}
				continue
			}
			key := keys[idx-1]
			entry, err := schema.Lookup(schema.CharacterFields(), key)
			if err != nil {
				return err
			}
			err = prompt.With(key, func() error {
				value, err := setting.Ask(entry.Field, entry.Get(c), s.Err)
				if err != nil {
					return err
				}
				var applyErr error
				currentName, applyErr = applyCharacterEdit(s, app, currentName, key, value, false)
				return applyErr
			})
			if err != nil {
				if errors.Is(err, prompt.ErrAborted) {
					continue
				}
				return err
			}
		}
	})
}

func applyCharacterEdit(s streams.Streams, app *api.App, oldName, key, raw string, normalize bool) (string, error) {
	updated, result, err := action.ApplyCharacterEdit(s.Ctx, app, oldName, key, raw, normalize)
	if err != nil {
		return "", err
	}
	if result.Warning != "" {
		_ = render.Warning(s.Err, result.Warning)
	}
	if err := render.Line(s.Out, "%s: %s\n", key, result.Value); err != nil {
		return "", err
	}
	return updated.Name, nil
}

func leafLabel(name string, meta prompt.LeafMeta) string {
	value := meta.Value
	if value == "" {
		value = "(empty)"
	}
	line := value
	if meta.Description != "" {
		line += "  —  " + meta.Description
	}
	return name + "\n  " + line
}

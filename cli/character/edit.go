package character

import (
	"errors"
	"fmt"

	"yukumo-script/cli/character/voice"
	"yukumo-script/cli/internal/cobraio"
	"yukumo-script/cli/prompt"
	"yukumo-script/cli/schema"
	"yukumo-script/cli/setting"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func newEditCmd(app *api.App, opts Options) *cobra.Command {
	return &cobra.Command{
		Use:   "edit [name] [field] [value]",
		Short: "Update a character field",
		Args: func(cmd *cobra.Command, args []string) error {
			if opts.Interactive {
				return cobra.MaximumNArgs(3)(cmd, args)
			}
			return cobra.ExactArgs(3)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			s := cobraio.FromCmd(cmd)
			if opts.Interactive && len(args) < 3 {
				return runInteractiveEdit(s, app, args)
			}
			return runOneShotEdit(s, app, args)
		},
	}
}

func runOneShotEdit(s streams.Streams, app *api.App, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("name, field, and value are required")
	}
	c, err := app.GetCharacter(args[0])
	if err != nil {
		return err
	}
	return applyEdit(s, app, args[0], &c, args[1], args[2], true)
}

func runInteractiveEdit(s streams.Streams, app *api.App, args []string) error {
	name := ""
	if len(args) >= 1 {
		name = args[0]
	}
	if err := ensureName(app, &name); err != nil {
		return err
	}

	if len(args) >= 2 {
		c, err := app.GetCharacter(name)
		if err != nil {
			return err
		}
		entry, err := schema.Lookup(schema.CharacterFields(), args[1])
		if err != nil {
			return err
		}
		err = prompt.With(name, func() error {
			return prompt.With(args[1], func() error {
				value, err := setting.Ask(entry.Field, entry.Get(c), s.Err)
				if err != nil {
					return err
				}
				return applyEdit(s, app, name, &c, args[1], value, false)
			})
		})
		if err != nil {
			return err
		}
		name = c.Name
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
				if err := voice.RunMenu(s, app, currentName); err != nil {
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
				return applyEdit(s, app, currentName, &c, key, value, false)
			})
			if err != nil {
				if errors.Is(err, prompt.ErrAborted) {
					continue
				}
				return err
			}
			currentName = c.Name
		}
	})
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

func applyEdit(
	s streams.Streams,
	app *api.App,
	oldName string,
	c *api.Character,
	key, raw string,
	normalize bool,
) error {
	entry, err := schema.Lookup(schema.CharacterFields(), key)
	if err != nil {
		return err
	}
	value := raw
	if normalize {
		var warning string
		value, warning, err = schema.Normalize(entry.Field, raw)
		if err != nil {
			return err
		}
		if warning != "" {
			_, _ = fmt.Fprint(s.Err, warning)
		}
	}
	if err := entry.Set(c, value); err != nil {
		return err
	}
	if err := app.UpdateCharacter(oldName, *c); err != nil {
		return err
	}
	_, err = fmt.Fprintf(s.Out, "%s: %s\n", key, entry.Get(*c))
	return err
}

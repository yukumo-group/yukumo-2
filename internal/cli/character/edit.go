package character

import (
	"errors"
	"fmt"

	"yukumo-script/internal/cli/character/voice"
	"yukumo-script/internal/cli/prompt"
	"yukumo-script/internal/cli/setting"
	"yukumo-script/internal/domain"
	svccharacter "yukumo-script/internal/service/character"

	"github.com/spf13/cobra"
)

func newEditCmd(svc *svccharacter.Service, opts Options) *cobra.Command {
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
			if opts.Interactive && len(args) < 3 {
				return runInteractiveEdit(cmd, svc, args)
			}
			return runOneShotEdit(cmd, svc, args)
		},
	}
}

func runOneShotEdit(cmd *cobra.Command, svc *svccharacter.Service, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("name, field, and value are required")
	}
	c, err := svc.Get(args[0])
	if err != nil {
		return err
	}
	return applyEdit(cmd, svc, args[0], &c, args[1], args[2], true)
}

func runInteractiveEdit(cmd *cobra.Command, svc *svccharacter.Service, args []string) error {
	name := ""
	if len(args) >= 1 {
		name = args[0]
	}
	if err := ensureName(svc, &name); err != nil {
		return err
	}

	if len(args) >= 2 {
		c, err := svc.Get(name)
		if err != nil {
			return err
		}
		entry, err := setting.Lookup(fields, args[1])
		if err != nil {
			return err
		}
		value, err := setting.Ask(entry.Field, entry.Get(c), cmd.ErrOrStderr())
		if err != nil {
			return err
		}
		if err := applyEdit(cmd, svc, name, &c, args[1], value, false); err != nil {
			return err
		}
		name = c.Name
	}

	currentName := name
	for {
		c, err := svc.Get(currentName)
		if err != nil {
			return err
		}
		metas := setting.Metas(fields, c)
		labels := []string{prompt.LabelBack}
		keys := setting.Keys(fields)
		for _, key := range keys {
			labels = append(labels, leafLabel(key, metas[key]))
		}
		labels = append(labels, "voices/\n  Edit voice profiles")

		idx, err := prompt.Select("Edit character", currentName, labels, 0)
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
			if err := voice.RunMenu(cmd, svc, currentName); err != nil {
				return err
			}
			continue
		}
		key := keys[idx-1]
		entry, err := setting.Lookup(fields, key)
		if err != nil {
			return err
		}
		value, err := setting.Ask(entry.Field, entry.Get(c), cmd.ErrOrStderr())
		if err != nil {
			if errors.Is(err, prompt.ErrAborted) {
				continue
			}
			return err
		}
		if err := applyEdit(cmd, svc, currentName, &c, key, value, false); err != nil {
			return err
		}
		currentName = c.Name
	}
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
	cmd *cobra.Command,
	svc *svccharacter.Service,
	oldName string,
	c *domain.Character,
	key, raw string,
	normalize bool,
) error {
	entry, err := setting.Lookup(fields, key)
	if err != nil {
		return err
	}
	value := raw
	if normalize {
		value, err = setting.Normalize(entry.Field, raw, cmd.ErrOrStderr())
		if err != nil {
			return err
		}
	}
	if err := entry.Set(c, value); err != nil {
		return err
	}
	if err := svc.Update(oldName, *c); err != nil {
		return err
	}
	_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", key, entry.Get(*c))
	return err
}

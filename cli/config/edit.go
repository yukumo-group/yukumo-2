package config

import (
	"errors"
	"fmt"

	"yukumo-script/cli/action"
	"yukumo-script/cli/internal/cobraio"
	"yukumo-script/cli/prompt"
	"yukumo-script/cli/render"
	"yukumo-script/cli/schema"
	"yukumo-script/cli/setting"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func newEditCmd(app *api.App, opts Options) *cobra.Command {
	return &cobra.Command{
		Use:   "edit [key] [value]",
		Short: "Update a single setting",
		Args: func(cmd *cobra.Command, args []string) error {
			if opts.Interactive {
				return cobra.MaximumNArgs(2)(cmd, args)
			}
			return cobra.ExactArgs(2)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			s := cobraio.FromCmd(cmd)
			if opts.Interactive && len(args) < 2 {
				return runInteractiveEdit(s, app, args)
			}
			return runOneShotEdit(s, app, args)
		},
	}
}

func runOneShotEdit(s streams.Streams, app *api.App, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("key and value are required")
	}
	return applyConfigEdit(s, app, args[0], args[1], true)
}

func runInteractiveEdit(s streams.Streams, app *api.App, args []string) error {
	startPrefix, focusKey := "", ""
	if len(args) >= 1 {
		cfg, err := app.Config()
		if err != nil {
			return err
		}
		entry, err := schema.Lookup(schema.ConfigFields(), args[0])
		if err != nil {
			return err
		}
		value, err := setting.Ask(entry.Field, entry.Get(cfg), s.Err)
		if err != nil {
			return err
		}
		if err := applyConfigEdit(s, app, args[0], value, false); err != nil {
			return err
		}
		startPrefix = prompt.ParentPrefix(args[0])
		focusKey = args[0]
	}

	err := setting.BrowseKeys(schema.Keys(schema.ConfigFields()), func() map[string]prompt.LeafMeta {
		cfg, err := app.Config()
		if err != nil {
			return map[string]prompt.LeafMeta{}
		}
		return setting.Metas(schema.ConfigFields(), cfg)
	}, func(key string) error {
		cfg, err := app.Config()
		if err != nil {
			return err
		}
		entry, err := schema.Lookup(schema.ConfigFields(), key)
		if err != nil {
			return err
		}
		value, err := setting.Ask(entry.Field, entry.Get(cfg), s.Err)
		if err != nil {
			return err
		}
		return applyConfigEdit(s, app, key, value, false)
	}, startPrefix, focusKey)
	if errors.Is(err, prompt.ErrAborted) {
		return nil
	}
	return err
}

func applyConfigEdit(s streams.Streams, app *api.App, key, raw string, normalize bool) error {
	result, err := action.ApplyConfigEdit(s.Ctx, app, key, raw, normalize)
	if err != nil {
		return err
	}
	if result.Warning != "" {
		_ = render.Warning(s.Err, result.Warning)
	}
	return render.Line(s.Out, "%s: %s\n", key, result.Value)
}

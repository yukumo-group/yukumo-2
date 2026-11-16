package config

import (
	"errors"
	"fmt"

	"yukumo-script/cli/prompt"
	"yukumo-script/cli/setting"
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
			if opts.Interactive && len(args) < 2 {
				return runInteractiveEdit(cmd, app, args)
			}
			return runOneShotEdit(cmd, app, args)
		},
	}
}

func runOneShotEdit(cmd *cobra.Command, app *api.App, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("key and value are required")
	}
	cfg, err := app.Config()
	if err != nil {
		return err
	}
	return applyEdit(cmd, app, cfg, args[0], args[1], true)
}

func runInteractiveEdit(cmd *cobra.Command, app *api.App, args []string) error {
	startPrefix, focusKey := "", ""
	if len(args) >= 1 {
		cfg, err := app.Config()
		if err != nil {
			return err
		}
		entry, err := setting.Lookup(fields, args[0])
		if err != nil {
			return err
		}
		value, err := setting.Ask(entry.Field, entry.Get(cfg), cmd.ErrOrStderr())
		if err != nil {
			return err
		}
		if err := applyEdit(cmd, app, cfg, args[0], value, false); err != nil {
			return err
		}
		startPrefix = prompt.ParentPrefix(args[0])
		focusKey = args[0]
	}

	err := setting.BrowseKeys(setting.Keys(fields), func() map[string]prompt.LeafMeta {
		cfg, err := app.Config()
		if err != nil {
			return map[string]prompt.LeafMeta{}
		}
		return setting.Metas(fields, cfg)
	}, func(key string) error {
		cfg, err := app.Config()
		if err != nil {
			return err
		}
		entry, err := setting.Lookup(fields, key)
		if err != nil {
			return err
		}
		value, err := setting.Ask(entry.Field, entry.Get(cfg), cmd.ErrOrStderr())
		if err != nil {
			return err
		}
		return applyEdit(cmd, app, cfg, key, value, false)
	}, startPrefix, focusKey)
	if errors.Is(err, prompt.ErrAborted) {
		return nil
	}
	return err
}

func applyEdit(
	cmd *cobra.Command,
	app *api.App,
	cfg api.Config,
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
	if err := entry.Set(&cfg, value); err != nil {
		return err
	}
	if err := app.SaveConfig(cfg); err != nil {
		return err
	}
	_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", key, entry.Get(cfg))
	return err
}

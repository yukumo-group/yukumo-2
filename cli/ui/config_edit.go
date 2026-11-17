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

// EditConfig edits one setting or opens the settings menu.
func EditConfig(s streams.Streams, app *api.App, partialKey, _ string) error {
	startPrefix, focusKey := "", ""
	if partialKey != "" {
		cfg, err := app.Config()
		if err != nil {
			return err
		}
		entry, err := schema.Lookup(schema.ConfigFields(), partialKey)
		if err != nil {
			return err
		}
		value, err := setting.Ask(entry.Field, entry.Get(cfg), s.Err)
		if err != nil {
			return err
		}
		if err := applyConfigEdit(s, app, partialKey, value, false); err != nil {
			return err
		}
		startPrefix = prompt.ParentPrefix(partialKey)
		focusKey = partialKey
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

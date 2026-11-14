package config

import (
	"errors"
	"fmt"

	"yukumo-script/internal/cli/prompt"
	"yukumo-script/internal/cli/setting"
	"yukumo-script/internal/domain"
	svcconfig "yukumo-script/internal/service/config"

	"github.com/spf13/cobra"
)

func newEditCmd(svc *svcconfig.Service, opts Options) *cobra.Command {
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
				return runInteractiveEdit(cmd, svc, args)
			}
			return runOneShotEdit(cmd, svc, args)
		},
	}
}

func runOneShotEdit(cmd *cobra.Command, svc *svcconfig.Service, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("key and value are required")
	}
	cfg, err := svc.Get()
	if err != nil {
		return err
	}
	return applyEdit(cmd, svc, cfg, args[0], args[1], true)
}

func runInteractiveEdit(cmd *cobra.Command, svc *svcconfig.Service, args []string) error {
	startPrefix, focusKey := "", ""
	if len(args) >= 1 {
		cfg, err := svc.Get()
		if err != nil {
			return err
		}
		entry, err := lookupField(args[0])
		if err != nil {
			return err
		}
		value, err := setting.Ask(entry.Field, entry.Get(cfg), cmd.ErrOrStderr())
		if err != nil {
			return err
		}
		if err := applyEdit(cmd, svc, cfg, args[0], value, false); err != nil {
			return err
		}
		startPrefix = prompt.ParentPrefix(args[0])
		focusKey = args[0]
	}

	err := setting.BrowseKeys(fieldKeys(), func() map[string]string {
		cfg, err := svc.Get()
		if err != nil {
			return map[string]string{}
		}
		return fieldValues(cfg)
	}, func(key string) error {
		cfg, err := svc.Get()
		if err != nil {
			return err
		}
		entry, err := lookupField(key)
		if err != nil {
			return err
		}
		value, err := setting.Ask(entry.Field, entry.Get(cfg), cmd.ErrOrStderr())
		if err != nil {
			return err
		}
		return applyEdit(cmd, svc, cfg, key, value, false)
	}, startPrefix, focusKey)
	if errors.Is(err, prompt.ErrAborted) {
		return nil
	}
	return err
}

func applyEdit(
	cmd *cobra.Command,
	svc *svcconfig.Service,
	cfg domain.Config,
	key, raw string,
	normalize bool,
) error {
	entry, err := lookupField(key)
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
	if err := svc.Save(cfg); err != nil {
		return err
	}
	_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", key, entry.Get(cfg))
	return err
}

func fieldValues(cfg domain.Config) map[string]string {
	out := make(map[string]string, len(fields))
	for _, f := range fields {
		out[f.Field.Key] = f.Get(cfg)
	}
	return out
}

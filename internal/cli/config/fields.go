package config

import (
	"fmt"
	"strconv"

	"yukumo-script/internal/cli/setting"
	"yukumo-script/internal/domain"
)

type fieldEntry struct {
	Field setting.Field
	Get   func(domain.Config) string
	Set   func(*domain.Config, string) error
}

var fields = []fieldEntry{
	{
		Field: setting.Field{
			Key:         "paths.character_path",
			Title:       "Character settings path",
			Description: "Directory for character YAML files",
			Kind:        setting.KindDir,
		},
		Get: func(c domain.Config) string { return c.Paths.CharacterPath },
		Set: func(c *domain.Config, v string) error {
			c.Paths.CharacterPath = v
			return nil
		},
	},
	{
		Field: setting.Field{
			Key:         "paths.data_path",
			Title:       "Data path",
			Description: "Directory for app data and caches",
			Kind:        setting.KindDir,
		},
		Get: func(c domain.Config) string { return c.Paths.DataPath },
		Set: func(c *domain.Config, v string) error {
			c.Paths.DataPath = v
			return nil
		},
	},
	{
		Field: setting.Field{
			Key:         "characters.init_default_voice",
			Title:       "Initialize default voice",
			Description: "Create a default voice when adding a character",
			Kind:        setting.KindBool,
		},
		Get: func(c domain.Config) string {
			return strconv.FormatBool(c.Characters.InitDefaultVoice)
		},
		Set: func(c *domain.Config, v string) error {
			b, err := strconv.ParseBool(v)
			if err != nil {
				return err
			}
			c.Characters.InitDefaultVoice = b
			return nil
		},
	},
	{
		Field: setting.Field{
			Key:         "cli.default_interactive",
			Title:       "Default interactive CLI",
			Description: "Enter huh menus when command args are incomplete",
			Kind:        setting.KindBool,
		},
		Get: func(c domain.Config) string {
			return strconv.FormatBool(c.CLI.Interactive())
		},
		Set: func(c *domain.Config, v string) error {
			b, err := strconv.ParseBool(v)
			if err != nil {
				return err
			}
			c.CLI.SetInteractive(b)
			return nil
		},
	},
	{
		Field: setting.Field{
			Key:         "cli.with_accessibility",
			Title:       "Huh accessibility mode",
			Description: "Use huh accessible prompts for screen readers",
			Kind:        setting.KindBool,
		},
		Get: func(c domain.Config) string {
			return strconv.FormatBool(c.CLI.Accessibility())
		},
		Set: func(c *domain.Config, v string) error {
			b, err := strconv.ParseBool(v)
			if err != nil {
				return err
			}
			c.CLI.SetAccessibility(b)
			return nil
		},
	},
}

func fieldKeys() []string {
	keys := make([]string, len(fields))
	for i, f := range fields {
		keys[i] = f.Field.Key
	}
	return keys
}

func lookupField(key string) (fieldEntry, error) {
	for _, f := range fields {
		if f.Field.Key == key {
			return f, nil
		}
	}
	return fieldEntry{}, fmt.Errorf("unknown key %q (want %v)", key, fieldKeys())
}

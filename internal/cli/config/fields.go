package config

import (
	"strconv"

	"yukumo-script/internal/cli/setting"
	"yukumo-script/internal/domain"
)

var fields = []setting.Entry[domain.Config]{
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
			Key:         "paths.output_path",
			Title:       "Output path",
			Description: "Directory for generated output files",
			Kind:        setting.KindDir,
		},
		Get: func(c domain.Config) string { return c.Paths.OutputPath },
		Set: func(c *domain.Config, v string) error {
			c.Paths.OutputPath = v
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
			Key:         "synth.output_format",
			Title:       "Synth output filename format",
			Description: "Token pattern for output names (c=character, v=voice, ymd=date)",
			Kind:        setting.KindString,
		},
		Get: func(c domain.Config) string { return c.Synth.OutputFormat },
		Set: func(c *domain.Config, v string) error {
			c.Synth.OutputFormat = v
			return nil
		},
	},
	{
		Field: setting.Field{
			Key:         "synth.always_write",
			Title:       "Always write synth output",
			Description: "Write WAV even when only previewing",
			Kind:        setting.KindBool,
		},
		Get: func(c domain.Config) string {
			return strconv.FormatBool(c.Synth.WriteAlways())
		},
		Set: func(c *domain.Config, v string) error {
			b, err := strconv.ParseBool(v)
			if err != nil {
				return err
			}
			c.Synth.SetAlwaysWrite(b)
			return nil
		},
	},
	{
		Field: setting.Field{
			Key:         "synth.reuse_cache",
			Title:       "Reuse synth cache",
			Description: "Skip synthesis when index hash matches an existing WAV",
			Kind:        setting.KindBool,
		},
		Get: func(c domain.Config) string {
			return strconv.FormatBool(c.Synth.CacheReuse())
		},
		Set: func(c *domain.Config, v string) error {
			b, err := strconv.ParseBool(v)
			if err != nil {
				return err
			}
			c.Synth.SetReuseCache(b)
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

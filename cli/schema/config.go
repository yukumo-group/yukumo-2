package schema

import (
	"strconv"

	"yukumo-script/core/api"
)

var configFields = []Entry[api.Config]{
	{
		Field: Field{
			Key:         "paths.character_path",
			Title:       "Character settings path",
			Description: "Directory for character YAML files",
			Kind:        KindDir,
		},
		Get: func(c api.Config) string { return c.Paths.CharacterPath },
		Set: func(c *api.Config, v string) error {
			c.Paths.CharacterPath = v
			return nil
		},
	},
	{
		Field: Field{
			Key:         "paths.output_path",
			Title:       "Output path",
			Description: "Directory for generated output files",
			Kind:        KindDir,
		},
		Get: func(c api.Config) string { return c.Paths.OutputPath },
		Set: func(c *api.Config, v string) error {
			c.Paths.OutputPath = v
			return nil
		},
	},
	{
		Field: Field{
			Key:         "characters.init_default_voice",
			Title:       "Initialize default voice",
			Description: "Create a default voice when adding a character",
			Kind:        KindBool,
		},
		Get: func(c api.Config) string {
			return strconv.FormatBool(c.Characters.InitDefaultVoice)
		},
		Set: func(c *api.Config, v string) error {
			b, err := strconv.ParseBool(v)
			if err != nil {
				return err
			}
			c.Characters.InitDefaultVoice = b
			return nil
		},
	},
	{
		Field: Field{
			Key:         "synth.output_format",
			Title:       "Synth output filename format",
			Description: "Braced tokens e.g. {character}:{voice}_{yyyy}-{MM}-{dd}:{hh}-{mm}-{ss}_{ms}",
			Kind:        KindString,
		},
		Get: func(c api.Config) string { return c.Synth.OutputFormat },
		Set: func(c *api.Config, v string) error {
			c.Synth.OutputFormat = v
			return nil
		},
	},
	{
		Field: Field{
			Key:         "synth.always_write",
			Title:       "Always write synth output",
			Description: "Write WAV even when only previewing",
			Kind:        KindBool,
		},
		Get: func(c api.Config) string {
			return strconv.FormatBool(c.Synth.AlwaysWrite)
		},
		Set: func(c *api.Config, v string) error {
			b, err := strconv.ParseBool(v)
			if err != nil {
				return err
			}
			c.Synth.AlwaysWrite = b
			return nil
		},
	},
	{
		Field: Field{
			Key:         "synth.reuse_cache",
			Title:       "Reuse synth cache",
			Description: "Skip synthesis when index hash matches an existing WAV",
			Kind:        KindBool,
		},
		Get: func(c api.Config) string {
			return strconv.FormatBool(c.Synth.ReuseCache)
		},
		Set: func(c *api.Config, v string) error {
			b, err := strconv.ParseBool(v)
			if err != nil {
				return err
			}
			c.Synth.ReuseCache = b
			return nil
		},
	},
}

// ConfigFields returns the editable config fields.
func ConfigFields() []Entry[api.Config] { return configFields }

package character

import (
	"yukumo-script/cli/setting"
	"yukumo-script/core/api"
)

var fields = []setting.Entry[api.Character]{
	{
		Field: setting.Field{
			Key:         "name",
			Title:       "Name",
			Description: "Unique character name",
			Kind:        setting.KindString,
		},
		Get: func(c api.Character) string { return c.Name },
		Set: func(c *api.Character, v string) error {
			c.Name = v
			return nil
		},
	},
	{
		Field: setting.Field{
			Key:         "description",
			Title:       "Description",
			Description: "Optional notes",
			Kind:        setting.KindString,
		},
		Get: func(c api.Character) string { return c.Description },
		Set: func(c *api.Character, v string) error {
			c.Description = v
			return nil
		},
	},
}

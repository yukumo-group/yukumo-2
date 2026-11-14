package character

import (
	"yukumo-script/internal/cli/setting"
	"yukumo-script/internal/domain"
)

var fields = []setting.Entry[domain.Character]{
	{
		Field: setting.Field{
			Key:         "name",
			Title:       "Name",
			Description: "Unique character name",
			Kind:        setting.KindString,
		},
		Get: func(c domain.Character) string { return c.Name },
		Set: func(c *domain.Character, v string) error {
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
		Get: func(c domain.Character) string { return c.Description },
		Set: func(c *domain.Character, v string) error {
			c.Description = v
			return nil
		},
	},
}

package schema

import (
	"yukumo-script/core/api"
)

var characterFields = []Entry[api.Character]{
	{
		Field: Field{
			Key:         "name",
			Title:       "Name",
			Description: "Unique character name",
			Kind:        KindString,
		},
		Get: func(c api.Character) string { return c.Name },
		Set: func(c *api.Character, v string) error {
			c.Name = v
			return nil
		},
	},
	{
		Field: Field{
			Key:         "description",
			Title:       "Description",
			Description: "Optional notes",
			Kind:        KindString,
		},
		Get: func(c api.Character) string { return c.Description },
		Set: func(c *api.Character, v string) error {
			c.Description = v
			return nil
		},
	},
}

// CharacterFields returns the editable character fields.
func CharacterFields() []Entry[api.Character] { return characterFields }

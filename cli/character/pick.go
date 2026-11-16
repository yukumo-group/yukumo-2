package character

import (
	"fmt"

	"yukumo-script/cli/prompt"
	"yukumo-script/core/api"

	"github.com/charmbracelet/lipgloss"
)

var dimMeta = lipgloss.NewStyle().Faint(true)

func ensureName(app *api.App, name *string) error {
	if name == nil {
		return fmt.Errorf("name is nil")
	}
	if *name != "" {
		return nil
	}
	chars, err := app.ListCharacters()
	if err != nil {
		return err
	}
	if len(chars) == 0 {
		return fmt.Errorf("no characters")
	}
	labels := make([]string, len(chars))
	for i, c := range chars {
		labels[i] = characterLabel(c)
	}
	idx, err := prompt.Select("character", "↑↓ move · enter", labels, 0)
	if err != nil {
		return err
	}
	*name = chars[idx].Name
	return nil
}

func characterLabel(c api.Character) string {
	line := c.Name
	if c.Description != "" {
		line += "\n  " + dimMeta.Render(c.Description)
	}
	return line
}

package ui

import (
	"fmt"

	"yukumo-script/cli/prompt"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"

	"github.com/charmbracelet/lipgloss"
)

var dimCharacter = lipgloss.NewStyle().Faint(true)

// PickCharacter returns name, or prompts for a character when it is empty.
func PickCharacter(_ streams.Streams, app *api.App, name string) (string, error) {
	if name != "" {
		return name, nil
	}
	chars, err := app.ListCharacters()
	if err != nil {
		return "", err
	}
	if len(chars) == 0 {
		return "", fmt.Errorf("no characters")
	}
	labels := make([]string, len(chars))
	for i, c := range chars {
		labels[i] = characterLabel(c)
	}
	idx, err := prompt.Select("character", "↑↓ move · enter", labels, 0)
	if err != nil {
		return "", err
	}
	return chars[idx].Name, nil
}

func characterLabel(c api.Character) string {
	line := c.Name
	if c.Description != "" {
		line += "\n  " + dimCharacter.Render(c.Description)
	}
	return line
}

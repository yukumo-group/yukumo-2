package character

import (
	"fmt"

	"yukumo-script/internal/cli/prompt"
	"yukumo-script/internal/domain"
	svccharacter "yukumo-script/internal/service/character"

	"github.com/charmbracelet/lipgloss"
)

var dimMeta = lipgloss.NewStyle().Faint(true)

func ensureName(svc *svccharacter.Service, name *string) error {
	if name == nil {
		return fmt.Errorf("name is nil")
	}
	if *name != "" {
		return nil
	}
	chars, err := svc.List()
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

func characterLabel(c domain.Character) string {
	line := c.Name
	if c.Description != "" {
		line += "\n  " + dimMeta.Render(c.Description)
	}
	return line
}

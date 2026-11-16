package synth

import (
	"fmt"

	"yukumo-script/cli/prompt"
	"yukumo-script/core/api"

	"github.com/charmbracelet/lipgloss"
)

var dimMeta = lipgloss.NewStyle().Faint(true)

func ensureVoiceRef(app *api.App, interactive bool, ref *string) error {
	if ref == nil {
		return fmt.Errorf("voice ref is nil")
	}
	if *ref != "" {
		return nil
	}
	if !interactive {
		return fmt.Errorf("voice is required (-v character:voice)")
	}
	list, err := app.ListCharacters()
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return fmt.Errorf("no characters")
	}
	labels := make([]string, len(list))
	for i, c := range list {
		labels[i] = characterLabel(c)
	}
	idx, err := prompt.Select("character", "↑↓ move · enter", labels, 0)
	if err != nil {
		return err
	}
	charName := list[idx].Name

	vs, err := app.ListVoices(charName)
	if err != nil {
		return err
	}
	if len(vs) == 0 {
		return fmt.Errorf("no voices for %q", charName)
	}
	vlabels := make([]string, len(vs))
	for i, v := range vs {
		vlabels[i] = voiceLabel(v)
	}
	vidx, err := prompt.Select("voice", "↑↓ move · enter", vlabels, 0)
	if err != nil {
		return err
	}
	*ref = charName + ":" + vs[vidx].Name
	return nil
}

func characterLabel(c api.Character) string {
	line := c.Name
	if c.Description != "" {
		line += "\n  " + dimMeta.Render(c.Description)
	}
	return line
}

func voiceLabel(v api.Voice) string {
	line := v.Name
	if n := len(v.Params); n > 0 {
		line += "\n  " + dimMeta.Render(fmt.Sprintf("%d param(s)", n))
	}
	return line
}

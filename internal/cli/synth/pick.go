package synth

import (
	"fmt"

	"yukumo-script/internal/cli/prompt"
	"yukumo-script/internal/domain"
	svccharacter "yukumo-script/internal/service/character"
	svcvoice "yukumo-script/internal/service/character/voice"

	"github.com/charmbracelet/lipgloss"
)

var dimMeta = lipgloss.NewStyle().Faint(true)

func ensureVoiceRef(chars *svccharacter.Service, voices *svcvoice.Service, interactive bool, ref *string) error {
	if ref == nil {
		return fmt.Errorf("voice ref is nil")
	}
	if *ref != "" {
		return nil
	}
	if !interactive {
		return fmt.Errorf("voice is required (-v character.voice)")
	}
	list, err := chars.List()
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

	vs, err := voices.List(charName)
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
	*ref = charName + "." + vs[vidx].Name
	return nil
}

func characterLabel(c domain.Character) string {
	line := c.Name
	if c.Description != "" {
		line += "\n  " + dimMeta.Render(c.Description)
	}
	return line
}

func voiceLabel(v domain.Voice) string {
	line := v.Name
	if n := len(v.Params); n > 0 {
		line += "\n  " + dimMeta.Render(fmt.Sprintf("%d param(s)", n))
	}
	return line
}

package voice

import (
	"fmt"

	"yukumo-script/internal/cli/prompt"
	"yukumo-script/internal/domain"
	svcvoice "yukumo-script/internal/service/character/voice"

	"github.com/charmbracelet/lipgloss"
)

var dimMeta = lipgloss.NewStyle().Faint(true)

func ensureVoiceName(svc *svcvoice.Service, charName string, name *string) error {
	if name == nil {
		return fmt.Errorf("voice name is nil")
	}
	if *name != "" {
		return nil
	}
	voices, err := svc.List(charName)
	if err != nil {
		return err
	}
	if len(voices) == 0 {
		return fmt.Errorf("no voices")
	}
	labels := make([]string, len(voices))
	for i, v := range voices {
		labels[i] = voiceLabel(v)
	}
	idx, err := prompt.Select("voice", "↑↓ move · enter", labels, 0)
	if err != nil {
		return err
	}
	*name = voices[idx].Name
	return nil
}

func voiceLabel(v domain.Voice) string {
	line := v.Name
	if n := len(v.Params); n > 0 {
		line += "\n  " + dimMeta.Render(fmt.Sprintf("%d param(s)", n))
	}
	return line
}

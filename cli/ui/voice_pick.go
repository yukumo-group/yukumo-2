package ui

import (
	"fmt"

	"yukumo-script/cli/prompt"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"

	"github.com/charmbracelet/lipgloss"
)

var dimVoice = lipgloss.NewStyle().Faint(true)

// PickVoice returns name, or prompts for a voice when it is empty.
func PickVoice(_ streams.Streams, app *api.App, charName, name string) (string, error) {
	if name != "" {
		return name, nil
	}
	voices, err := app.ListVoices(charName)
	if err != nil {
		return "", err
	}
	if len(voices) == 0 {
		return "", fmt.Errorf("no voices")
	}
	labels := make([]string, len(voices))
	for i, v := range voices {
		labels[i] = voiceLabel(v)
	}
	idx, err := prompt.Select("voice", "↑↓ move · enter", labels, 0)
	if err != nil {
		return "", err
	}
	return voices[idx].Name, nil
}

func voiceLabel(v api.Voice) string {
	line := v.Name
	if n := len(v.Params); n > 0 {
		line += "\n  " + dimVoice.Render(fmt.Sprintf("%d param(s)", n))
	}
	return line
}

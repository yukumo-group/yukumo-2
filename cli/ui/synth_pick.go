package ui

import (
	"fmt"

	"yukumo-script/cli/prompt"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"
)

// EnsureSynthText prompts for synth text when it is empty.
func EnsureSynthText(_ streams.Streams, text *string) error {
	return prompt.EnsureString("text", text)
}

// PickVoiceRef prompts for a character and voice when ref is empty.
func PickVoiceRef(_ streams.Streams, app *api.App, ref *string) error {
	if ref == nil {
		return fmt.Errorf("voice ref is nil")
	}
	if *ref != "" {
		return nil
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

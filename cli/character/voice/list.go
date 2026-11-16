package voice

import (
	"fmt"

	"yukumo-script/cli/streams"
	"yukumo-script/core/api"

	"gopkg.in/yaml.v3"
)

func runList(s streams.Streams, app *api.App, charName string, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: character <name> voice list")
	}
	voices, err := app.ListVoices(charName)
	if err != nil {
		return err
	}
	enc := yaml.NewEncoder(s.Out)
	enc.SetIndent(2)
	defer func() { _ = enc.Close() }()
	return enc.Encode(voices)
}

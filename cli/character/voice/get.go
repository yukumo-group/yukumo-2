package voice

import (
	"fmt"

	"yukumo-script/cli/streams"
	"yukumo-script/core/api"

	"gopkg.in/yaml.v3"
)

func runGet(s streams.Streams, app *api.App, opts Options, charName string, args []string) error {
	name := ""
	if len(args) >= 1 {
		name = args[0]
	}
	if opts.Interactive || name == "" {
		if err := ensureVoiceName(app, charName, &name); err != nil {
			return err
		}
	}
	if name == "" {
		return fmt.Errorf("usage: character <name> voice get <voice>")
	}
	v, err := app.GetVoice(charName, name)
	if err != nil {
		return err
	}
	enc := yaml.NewEncoder(s.Out)
	enc.SetIndent(2)
	defer func() { _ = enc.Close() }()
	return enc.Encode(v)
}

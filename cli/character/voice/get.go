package voice

import (
	"fmt"

	"yukumo-script/cli/action"
	"yukumo-script/cli/render"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"
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
	v, err := action.GetVoice(s.Ctx, app, charName, name)
	if err != nil {
		return err
	}
	return render.YAML(s.Out, v)
}

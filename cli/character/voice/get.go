package voice

import (
	"fmt"

	"yukumo-script/cli/action"
	"yukumo-script/cli/render"
	"yukumo-script/cli/streams"
	"yukumo-script/cli/ui"
	"yukumo-script/core/api"
)

func runGet(s streams.Streams, app *api.App, interactive bool, charName string, args []string) error {
	name := ""
	if len(args) >= 1 {
		name = args[0]
	}
	if interactive {
		var err error
		name, err = ui.PickVoice(s, app, charName, name)
		if err != nil {
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

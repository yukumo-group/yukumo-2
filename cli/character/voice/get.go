package voice

import (
	"fmt"

	"yukumo-script/cli/action"
	"yukumo-script/cli/render"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"
)

func runGet(s streams.Streams, app *api.App, charName string, args []string) error {
	if len(args) < 1 || args[0] == "" {
		return fmt.Errorf("usage: character <name> voice get <voice>")
	}
	v, err := action.GetVoice(s.Ctx, app, charName, args[0])
	if err != nil {
		return err
	}
	return render.YAML(s.Out, v)
}

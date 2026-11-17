package voice

import (
	"fmt"

	"yukumo-script/cli/action"
	"yukumo-script/cli/render"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"
)

func runList(s streams.Streams, app *api.App, charName string, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: character <name> voice list")
	}
	voices, err := action.ListVoices(s.Ctx, app, charName)
	if err != nil {
		return err
	}
	return render.YAML(s.Out, voices)
}

package voice

import (
	"fmt"

	"yukumo-script/cli/action"
	"yukumo-script/cli/render"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"
)

func runDelete(s streams.Streams, app *api.App, charName string, args []string) error {
	if len(args) < 1 || args[0] == "" {
		return fmt.Errorf("voice name is required")
	}
	name := args[0]
	if err := action.DeleteVoice(s.Ctx, app, charName, name); err != nil {
		return err
	}
	return render.Line(s.Out, "deleted voice %q\n", name)
}

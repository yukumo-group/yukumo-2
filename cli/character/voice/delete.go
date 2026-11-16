package voice

import (
	"fmt"

	"yukumo-script/cli/action"
	"yukumo-script/cli/render"
	"yukumo-script/cli/streams"
	"yukumo-script/cli/ui"
	"yukumo-script/core/api"
)

func runDelete(s streams.Streams, app *api.App, interactive bool, charName string, args []string) error {
	name := ""
	if len(args) >= 1 {
		name = args[0]
	}
	if interactive {
		return ui.DeleteVoice(s, app, charName, name)
	}
	if name == "" {
		return fmt.Errorf("voice name is required")
	}
	if err := action.DeleteVoice(s.Ctx, app, charName, name); err != nil {
		return err
	}
	return render.Line(s.Out, "deleted voice %q\n", name)
}

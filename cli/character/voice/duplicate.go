package voice

import (
	"fmt"

	"yukumo-script/cli/action"
	"yukumo-script/cli/render"
	"yukumo-script/cli/streams"
	"yukumo-script/cli/ui"
	"yukumo-script/core/api"
)

func runDuplicate(s streams.Streams, app *api.App, interactive bool, charName string, args []string) error {
	src := ""
	if len(args) >= 1 {
		src = args[0]
	}
	newName := ""
	if len(args) >= 2 {
		newName = args[1]
	}
	if interactive {
		return ui.DuplicateVoice(s, app, charName, src, newName)
	}
	if src == "" {
		return fmt.Errorf("voice name is required")
	}
	final, err := action.DuplicateVoice(s.Ctx, app, charName, src, newName)
	if err != nil {
		return err
	}
	return render.Line(s.Out, "duplicated %q → %q\n", src, final)
}

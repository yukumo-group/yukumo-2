package voice

import (
	"fmt"

	"yukumo-script/cli/action"
	"yukumo-script/cli/render"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"
)

func runDuplicate(s streams.Streams, app *api.App, charName string, args []string) error {
	if len(args) < 1 || args[0] == "" {
		return fmt.Errorf("voice name is required")
	}
	src := args[0]
	newName := ""
	if len(args) >= 2 {
		newName = args[1]
	}
	final, err := action.DuplicateVoice(s.Ctx, app, charName, src, newName)
	if err != nil {
		return err
	}
	return render.Line(s.Out, "duplicated %q → %q\n", src, final)
}

package voice

import (
	"fmt"

	"yukumo-script/cli/action"
	"yukumo-script/cli/parse"
	"yukumo-script/cli/render"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"
)

func runAdd(s streams.Streams, app *api.App, charName string, args []string) error {
	if len(args) < 1 || args[0] == "" {
		return fmt.Errorf("voice name is required")
	}
	name := args[0]
	version := api.Version1
	if len(args) >= 2 {
		n, _, err := parse.ParseRangedInt(args[1], 1, 10)
		if err != nil {
			return err
		}
		version = n
	}
	if version != api.Version1 && version != api.Version2 && version != api.Version10 {
		return fmt.Errorf("%w: version must be 1, 2, or 10", api.ErrInvalid)
	}
	p := api.DefaultSynthParam(version)
	created, err := action.AddVoice(s.Ctx, app, charName, api.Voice{
		Name:   name,
		Speed:  100,
		Params: []api.SynthParam{p},
	})
	if err != nil {
		return err
	}
	return render.Line(s.Out, "added voice %q\n", created.Name)
}

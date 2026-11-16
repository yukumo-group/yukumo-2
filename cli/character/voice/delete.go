package voice

import (
	"errors"
	"fmt"

	"yukumo-script/cli/prompt"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"
)

func runDelete(s streams.Streams, app *api.App, opts Options, charName string, args []string) error {
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
		return fmt.Errorf("voice name is required")
	}
	if opts.Interactive {
		ok := false
		if err := prompt.Confirm(fmt.Sprintf("Delete voice %q?", name), &ok); err != nil {
			if errors.Is(err, prompt.ErrAborted) {
				return nil
			}
			return err
		}
		if !ok {
			return nil
		}
	}
	if err := app.DeleteVoice(charName, name); err != nil {
		return err
	}
	_, err := fmt.Fprintf(s.Out, "deleted voice %q\n", name)
	return err
}

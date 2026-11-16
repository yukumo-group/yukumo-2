package voice

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"yukumo-script/cli/prompt"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"
)

func runDuplicate(s streams.Streams, app *api.App, opts Options, charName string, args []string) error {
	src := ""
	newName := ""
	if len(args) >= 1 {
		src = args[0]
	}
	if len(args) >= 2 {
		newName = args[1]
	}
	if opts.Interactive || src == "" {
		if err := ensureVoiceName(app, charName, &src); err != nil {
			return err
		}
	}
	if src == "" {
		return fmt.Errorf("voice name is required")
	}
	if opts.Interactive && len(args) < 2 {
		def := defaultCopyName(src)
		newName = def
		if err := prompt.String(fmt.Sprintf("New name (default %q)", def), &newName); err != nil {
			if errors.Is(err, prompt.ErrAborted) {
				return nil
			}
			return err
		}
		if newName == "" {
			newName = def
		}
	}
	final, err := app.DuplicateVoice(context.Background(), charName, src, newName)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(s.Out, "duplicated %q → %q\n", src, final)
	return err
}

func defaultCopyName(name string) string {
	name = strings.TrimSpace(name)
	return name + " - copy"
}

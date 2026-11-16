package voice

import (
	"context"
	"errors"
	"fmt"

	"yukumo-script/cli/prompt"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func runAdd(cmd *cobra.Command, app *api.App, opts Options, charName string, args []string) error {
	name := ""
	if len(args) >= 1 {
		name = args[0]
	}
	if opts.Interactive {
		if err := prompt.EnsureString("Voice name", &name); err != nil {
			if errors.Is(err, prompt.ErrAborted) {
				return nil
			}
			return err
		}
	}
	if name == "" {
		return fmt.Errorf("voice name is required")
	}

	v := api.Voice{Name: name}
	version := api.Version1
	if len(args) >= 2 {
		n, clamped, err := prompt.ParseRangedInt(args[1], 1, 10)
		if err != nil {
			return err
		}
		_ = clamped
		version = n
	} else if opts.Interactive {
		if err := SelectVersion(&version); err != nil {
			if errors.Is(err, prompt.ErrAborted) {
				return nil
			}
			return err
		}
	}
	if version != api.Version1 && version != api.Version2 && version != api.Version10 {
		return fmt.Errorf("%w: version must be 1, 2, or 10", api.ErrInvalid)
	}

	p := api.DefaultSynthParam(version)
	if opts.Interactive {
		err := prompt.With(fmt.Sprintf("param v%d", version), func() error {
			return editParamInteractive(cmd, app, &p)
		})
		if err != nil {
			if errors.Is(err, prompt.ErrAborted) {
				return nil
			}
			return err
		}
	}
	v.Params = []api.SynthParam{p}

	if err := app.AddVoice(context.Background(), charName, v); err != nil {
		return err
	}
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "added voice %q\n", v.Name)
	return err
}

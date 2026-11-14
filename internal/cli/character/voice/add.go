package voice

import (
	"context"
	"errors"
	"fmt"

	"yukumo-script/internal/cli/prompt"
	"yukumo-script/internal/domain"
	svcvoice "yukumo-script/internal/service/character/voice"

	"github.com/spf13/cobra"
)

func runAdd(cmd *cobra.Command, svc *svcvoice.Service, opts Options, charName string, args []string) error {
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

	v := domain.Voice{Name: name}
	version := domain.Version1
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
	if version != domain.Version1 && version != domain.Version2 && version != domain.Version10 {
		return fmt.Errorf("%w: version must be 1, 2, or 10", domain.ErrInvalid)
	}

	p := domain.DefaultSynthParam(version)
	if opts.Interactive {
		err := prompt.With(fmt.Sprintf("param v%d", version), func() error {
			return editParamInteractive(cmd, svc, &p)
		})
		if err != nil {
			if errors.Is(err, prompt.ErrAborted) {
				return nil
			}
			return err
		}
	}
	v.Params = []domain.SynthParam{p}

	if err := svc.Add(context.Background(), charName, v); err != nil {
		return err
	}
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "added voice %q\n", v.Name)
	return err
}

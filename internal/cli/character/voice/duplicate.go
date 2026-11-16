package voice

import (
	"context"
	"errors"
	"fmt"

	"yukumo-script/internal/cli/prompt"
	svcvoice "yukumo-script/internal/service/character/voice"
	"yukumo-script/internal/service/nameutil"

	"github.com/spf13/cobra"
)

func runDuplicate(cmd *cobra.Command, svc *svcvoice.Service, opts Options, charName string, args []string) error {
	src := ""
	newName := ""
	if len(args) >= 1 {
		src = args[0]
	}
	if len(args) >= 2 {
		newName = args[1]
	}
	if opts.Interactive || src == "" {
		if err := ensureVoiceName(svc, charName, &src); err != nil {
			return err
		}
	}
	if src == "" {
		return fmt.Errorf("voice name is required")
	}
	if opts.Interactive && len(args) < 2 {
		def := nameutil.ResolveCopy(src, "")
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
	if err := svc.Duplicate(context.Background(), charName, src, newName); err != nil {
		return err
	}
	final := nameutil.ResolveCopy(src, newName)
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "duplicated %q → %q\n", src, final)
	return err
}

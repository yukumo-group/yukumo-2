package voice

import (
	"errors"
	"fmt"

	"yukumo-script/internal/cli/prompt"
	svccharacter "yukumo-script/internal/service/character"

	"github.com/spf13/cobra"
)

func runDelete(cmd *cobra.Command, svc *svccharacter.Service, opts Options, charName string, args []string) error {
	name := ""
	if len(args) >= 1 {
		name = args[0]
	}
	if opts.Interactive || name == "" {
		if err := ensureVoiceName(svc, charName, &name); err != nil {
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
	if err := svc.DeleteVoice(charName, name); err != nil {
		return err
	}
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "deleted voice %q\n", name)
	return err
}

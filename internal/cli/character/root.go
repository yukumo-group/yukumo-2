package character

import (
	"fmt"

	"yukumo-script/internal/cli/character/voice"
	svccharacter "yukumo-script/internal/service/character"

	"github.com/spf13/cobra"
)

type Options struct {
	Interactive bool
}

func NewCmd(svc *svccharacter.Service, opts Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "character",
		Short: "Manage characters",
		Long: `Manage characters and their voices.

Voice commands (non-interactive):
  character <name> voice list
  character <name> voice get <voice>
  character <name> voice add <voice> [version]
  character <name> voice edit <voice> ...
  character <name> voice delete <voice>`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) >= 2 && args[1] == "voice" {
				return voice.Dispatch(cmd, svc, voice.Options{Interactive: opts.Interactive}, args[0], args[2:])
			}
			if len(args) == 0 {
				return cmd.Help()
			}
			return fmt.Errorf("unknown character args %v (try list|add|edit|delete or <name> voice ...)", args)
		},
	}
	get := newGetCmd(svc, opts)
	if opts.Interactive {
		get.Hidden = true
	}
	cmd.AddCommand(
		newListCmd(svc),
		get,
		newAddCmd(svc, opts),
		newEditCmd(svc, opts),
		newDeleteCmd(svc, opts),
	)
	return cmd
}

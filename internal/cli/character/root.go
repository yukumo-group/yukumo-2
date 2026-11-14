package character

import (
	"fmt"

	"yukumo-script/internal/cli/character/voice"
	svccharacter "yukumo-script/internal/service/character"
	svcvoice "yukumo-script/internal/service/character/voice"

	"github.com/spf13/cobra"
)

type Options struct {
	Interactive bool
}

func NewCmd(chars *svccharacter.Service, voices *svcvoice.Service, opts Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "character",
		Short: "Manage characters",
		Long: `Manage characters and their voices.

  character <name> voice list|get|add|edit|delete|duplicate ...`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) >= 2 && args[1] == "voice" {
				return voice.Dispatch(cmd, voices, voice.Options{Interactive: opts.Interactive}, args[0], args[2:])
			}
			if len(args) == 0 {
				return cmd.Help()
			}
			return fmt.Errorf("unknown character args %v", args)
		},
	}
	get := newGetCmd(chars, opts)
	if opts.Interactive {
		get.Hidden = true
	}
	cmd.AddCommand(
		newListCmd(chars),
		get,
		newAddCmd(chars, opts),
		newEditCmd(chars, voices, opts),
		newDeleteCmd(chars, opts),
	)
	return cmd
}

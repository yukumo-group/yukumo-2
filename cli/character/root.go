package character

import (
	"fmt"

	"yukumo-script/cli/character/voice"
	"yukumo-script/cli/internal/cobraio"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

type Options struct {
	Interactive bool
}

func NewCmd(app *api.App, opts Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "character",
		Short: "Manage characters",
		Long: `Manage characters and their voices.

  character <name> voice list|get|add|edit|delete|duplicate ...`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) >= 2 && args[1] == "voice" {
				return voice.Dispatch(cobraio.FromCmd(cmd), app, voice.Options{Interactive: opts.Interactive}, args[0], args[2:])
			}
			if len(args) == 0 {
				return cmd.Help()
			}
			return fmt.Errorf("unknown character args %v", args)
		},
	}
	get := newGetCmd(app, opts)
	if opts.Interactive {
		get.Hidden = true
	}
	cmd.AddCommand(
		newListCmd(app),
		get,
		newAddCmd(app, opts),
		newEditCmd(app, opts),
		newDeleteCmd(app, opts),
	)
	return cmd
}

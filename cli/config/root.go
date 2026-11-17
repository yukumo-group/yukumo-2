package config

import (
	"yukumo-script/cli/internal/cobraio"
	"yukumo-script/cli/internal/cobraui"
	"yukumo-script/cli/ui"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

type Options struct {
	Interactive bool
}

func NewCmd(app *api.App, opts Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Read and update application settings",
	}
	get := newGetCmd(app, opts)
	if opts.Interactive {
		get.Hidden = true
	}
	cmd.AddCommand(newListCmd(app), get, newEditCmd(app, opts))
	if opts.Interactive {
		cmd.RunE = func(cmd *cobra.Command, _ []string) error {
			return ui.Browse(cobraio.FromCmd(cmd), cobraui.FromCobra(cmd))
		}
	}
	return cmd
}

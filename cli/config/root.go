package config

import (
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
	return cmd
}

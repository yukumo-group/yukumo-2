package config

import (
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func NewCmd(app *api.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Read and update application settings",
	}
	cmd.AddCommand(newListCmd(app), newGetCmd(app), newEditCmd(app))
	return cmd
}

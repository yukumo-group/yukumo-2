package config

import (
	"yukumo-script/cli/action"
	"yukumo-script/cli/internal/cobraio"
	"yukumo-script/cli/render"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func newEditCmd(app *api.App) *cobra.Command {
	return &cobra.Command{
		Use:   "edit <key> <value>",
		Short: "Update a single setting",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := cobraio.FromCmd(cmd)
			result, err := action.ApplyConfigEdit(s.Ctx, app, args[0], args[1], true)
			if err != nil {
				return err
			}
			if result.Warning != "" {
				_ = render.Warning(s.Err, result.Warning)
			}
			return render.Line(s.Out, "%s: %s\n", args[0], result.Value)
		},
	}
}

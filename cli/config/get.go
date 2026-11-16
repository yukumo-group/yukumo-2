package config

import (
	"yukumo-script/cli/action"
	"yukumo-script/cli/internal/cobraio"
	"yukumo-script/cli/render"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func newGetCmd(app *api.App, opts Options) *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Print one setting by key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = opts
			s := cobraio.FromCmd(cmd)
			value, err := action.GetConfigField(s.Ctx, app, args[0])
			if err != nil {
				return err
			}
			return render.Line(s.Out, "%s\n", value)
		},
	}
}

package character

import (
	"yukumo-script/cli/action"
	"yukumo-script/cli/internal/cobraio"
	"yukumo-script/cli/render"
	"yukumo-script/cli/ui"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func newGetCmd(app *api.App, opts Options) *cobra.Command {
	return &cobra.Command{
		Use:   "get [name]",
		Short: "Print one character",
		Args: func(cmd *cobra.Command, args []string) error {
			if opts.Interactive {
				return cobra.MaximumNArgs(1)(cmd, args)
			}
			return cobra.ExactArgs(1)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			s := cobraio.FromCmd(cmd)
			name := argAt(args, 0)
			if opts.Interactive {
				var err error
				name, err = ui.PickCharacter(s, app, name)
				if err != nil {
					return err
				}
			}
			c, err := action.GetCharacter(s.Ctx, app, name)
			if err != nil {
				return err
			}
			return render.YAML(s.Out, c)
		},
	}
}

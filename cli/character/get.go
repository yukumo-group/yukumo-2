package character

import (
	"yukumo-script/cli/action"
	"yukumo-script/cli/internal/cobraio"
	"yukumo-script/cli/render"
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
			name := ""
			if len(args) >= 1 {
				name = args[0]
			}
			if opts.Interactive {
				if err := ensureName(app, &name); err != nil {
					return err
				}
			}
			s := cobraio.FromCmd(cmd)
			c, err := action.GetCharacter(s.Ctx, app, name)
			if err != nil {
				return err
			}
			return render.YAML(s.Out, c)
		},
	}
}

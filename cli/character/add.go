package character

import (
	"yukumo-script/cli/action"
	"yukumo-script/cli/internal/cobraio"
	"yukumo-script/cli/render"
	"yukumo-script/cli/streams"
	"yukumo-script/cli/ui"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func newAddCmd(app *api.App, opts Options) *cobra.Command {
	return &cobra.Command{
		Use:   "add [name] [description]",
		Short: "Create a character",
		Args: func(cmd *cobra.Command, args []string) error {
			if opts.Interactive {
				return cobra.MaximumNArgs(2)(cmd, args)
			}
			return cobra.RangeArgs(1, 2)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			s := cobraio.FromCmd(cmd)
			if opts.Interactive {
				return ui.AddCharacter(s, app, argAt(args, 0), argAt(args, 1), len(args) >= 2)
			}
			return runOneShotAdd(s, app, args)
		},
	}
}

func runOneShotAdd(s streams.Streams, app *api.App, args []string) error {
	name := argAt(args, 0)
	description := argAt(args, 1)
	created, err := action.AddCharacter(s.Ctx, app, name, description)
	if err != nil {
		return err
	}
	return render.Line(s.Out, "added %q\n", created.Name)
}

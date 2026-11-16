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

func newDeleteCmd(app *api.App, opts Options) *cobra.Command {
	return &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete a character",
		Args: func(cmd *cobra.Command, args []string) error {
			if opts.Interactive {
				return cobra.MaximumNArgs(1)(cmd, args)
			}
			return cobra.ExactArgs(1)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			s := cobraio.FromCmd(cmd)
			if opts.Interactive {
				return ui.DeleteCharacter(s, app, argAt(args, 0))
			}
			return runOneShotDelete(s, app, args)
		},
	}
}

func runOneShotDelete(s streams.Streams, app *api.App, args []string) error {
	name := argAt(args, 0)
	if err := action.DeleteCharacter(s.Ctx, app, name); err != nil {
		return err
	}
	return render.Line(s.Out, "deleted %q\n", name)
}

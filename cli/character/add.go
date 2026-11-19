package character

import (
	"yukumo-script/cli/action"
	"yukumo-script/cli/internal/cobraio"
	"yukumo-script/cli/render"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func newAddCmd(app *api.App) *cobra.Command {
	return &cobra.Command{
		Use:   "add <name> [description]",
		Short: "Create a character",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := cobraio.FromCmd(cmd)
			description := ""
			if len(args) >= 2 {
				description = args[1]
			}
			created, err := action.AddCharacter(s.Ctx, app, args[0], description)
			if err != nil {
				return err
			}
			return render.Line(s.Out, "added %q\n", created.Name)
		},
	}
}

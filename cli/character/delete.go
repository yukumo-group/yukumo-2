package character

import (
	"yukumo-script/cli/action"
	"yukumo-script/cli/internal/cobraio"
	"yukumo-script/cli/render"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func newDeleteCmd(app *api.App) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a character",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := cobraio.FromCmd(cmd)
			if err := action.DeleteCharacter(s.Ctx, app, args[0]); err != nil {
				return err
			}
			return render.Line(s.Out, "deleted %q\n", args[0])
		},
	}
}

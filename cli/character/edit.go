package character

import (
	"yukumo-script/cli/action"
	"yukumo-script/cli/internal/cobraio"
	"yukumo-script/cli/render"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func newEditCmd(app *api.App) *cobra.Command {
	return &cobra.Command{
		Use:   "edit <name> <field> <value>",
		Short: "Update a character field",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := cobraio.FromCmd(cmd)
			_, result, err := action.ApplyCharacterEdit(s.Ctx, app, args[0], args[1], args[2], true)
			if err != nil {
				return err
			}
			if result.Warning != "" {
				_ = render.Warning(s.Err, result.Warning)
			}
			return render.Line(s.Out, "%s: %s\n", args[1], result.Value)
		},
	}
}

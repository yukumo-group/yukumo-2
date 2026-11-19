package character

import (
	"yukumo-script/cli/action"
	"yukumo-script/cli/internal/cobraio"
	"yukumo-script/cli/render"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func newGetCmd(app *api.App) *cobra.Command {
	return &cobra.Command{
		Use:   "get <name>",
		Short: "Print one character",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := cobraio.FromCmd(cmd)
			c, err := action.GetCharacter(s.Ctx, app, args[0])
			if err != nil {
				return err
			}
			return render.YAML(s.Out, c)
		},
	}
}

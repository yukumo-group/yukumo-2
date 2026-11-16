package character

import (
	"yukumo-script/cli/action"
	"yukumo-script/cli/internal/cobraio"
	"yukumo-script/cli/render"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func newListCmd(app *api.App) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Print all characters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			s := cobraio.FromCmd(cmd)
			chars, err := action.ListCharacters(s.Ctx, app)
			if err != nil {
				return err
			}
			return render.YAML(s.Out, chars)
		},
	}
}

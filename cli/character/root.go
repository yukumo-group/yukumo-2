package character

import (
	"fmt"

	"yukumo-script/cli/character/voice"
	"yukumo-script/cli/internal/cobraio"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func NewCmd(app *api.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "character",
		Short: "Manage characters",
		Long: `Manage characters and their voices.

  character <name> voice list|get|add|edit|delete|duplicate ...`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) >= 2 && args[1] == "voice" {
				return voice.Dispatch(cobraio.FromCmd(cmd), app, args[0], args[2:])
			}
			if len(args) == 0 {
				return cmd.Help()
			}
			return fmt.Errorf("unknown character args %v", args)
		},
	}
	cmd.AddCommand(
		newListCmd(app),
		newGetCmd(app),
		newAddCmd(app),
		newEditCmd(app),
		newDeleteCmd(app),
	)
	return cmd
}

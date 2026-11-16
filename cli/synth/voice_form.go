package synth

import (
	"fmt"

	"yukumo-script/cli/internal/cobraio"
	"yukumo-script/cli/ui"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func newVoiceCmd(app *api.App, opts Options) *cobra.Command {
	return &cobra.Command{
		Use:   "voice",
		Short: "Synthesize from a character voice",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !opts.Interactive {
				return fmt.Errorf("interactive only; use: synth -v character:voice TEXT")
			}
			return ui.SynthVoiceForm(cobraio.FromCmd(cmd), app)
		},
	}
}

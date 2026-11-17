package cli

import (
	"yukumo-script/cli/internal/cobraio"
	"yukumo-script/cli/internal/cobraui"
	"yukumo-script/cli/ui"

	"github.com/spf13/cobra"
)

func enableInteractiveRoot(root *cobra.Command) {
	spec := cobraui.FromCobra(root)
	root.RunE = func(cmd *cobra.Command, _ []string) error {
		return ui.Browse(cobraio.FromCmd(cmd), spec)
	}
}

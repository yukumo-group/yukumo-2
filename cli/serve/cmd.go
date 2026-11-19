package serve

import (
	"yukumo-script/cli/internal/cobraio"
	"yukumo-script/core/api"
	"yukumo-script/core/api/rpc"

	"github.com/spf13/cobra"
)

func NewCmd(app *api.App) *cobra.Command {
	return &cobra.Command{
		Use:           "serve",
		Short:         "JSON-RPC 2.0 over stdin/stdout (NDJSON)",
		Hidden:        true,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			s := cobraio.FromCmd(cmd)
			return Loop(s.Ctx, rpc.Bind(app), s.In, s.Out, s.Err)
		},
	}
}

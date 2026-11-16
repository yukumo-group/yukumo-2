package config

import (
	"fmt"

	"yukumo-script/cli/schema"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func newGetCmd(app *api.App, opts Options) *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Print one setting by key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = opts
			cfg, err := app.Config()
			if err != nil {
				return err
			}
			entry, err := schema.Lookup(schema.ConfigFields(), args[0])
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), entry.Get(cfg))
			return err
		},
	}
}

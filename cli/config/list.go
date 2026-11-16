package config

import (
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newListCmd(app *api.App) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Print all settings",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := app.Config()
			if err != nil {
				return err
			}
			enc := yaml.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent(2)
			defer func() { _ = enc.Close() }()
			return enc.Encode(cfg)
		},
	}
}

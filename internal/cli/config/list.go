package config

import (
	svcconfig "yukumo-script/internal/service/config"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newListCmd(svc *svcconfig.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Print all settings",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := svc.Get()
			if err != nil {
				return err
			}
			enc := yaml.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent(2)
			defer enc.Close()
			return enc.Encode(cfg)
		},
	}
}

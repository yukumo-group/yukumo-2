package config

import (
	"fmt"

	svcconfig "yukumo-script/internal/service/config"

	"github.com/spf13/cobra"
)

func newGetCmd(svc *svcconfig.Service, opts Options) *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Print one setting by key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = opts
			cfg, err := svc.Get()
			if err != nil {
				return err
			}
			entry, err := lookupField(args[0])
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), entry.Get(cfg))
			return err
		},
	}
}

package config

import (
	svcconfig "yukumo-script/internal/service/config"

	"github.com/spf13/cobra"
)

type Options struct {
	Interactive bool
}

func NewCmd(svc *svcconfig.Service, opts Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Read and update application settings",
	}
	get := newGetCmd(svc, opts)
	if opts.Interactive {
		get.Hidden = true
	}
	cmd.AddCommand(newListCmd(svc), get, newEditCmd(svc, opts))
	return cmd
}

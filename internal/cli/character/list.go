package character

import (
	svccharacter "yukumo-script/internal/service/character"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newListCmd(svc *svccharacter.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Print all characters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			chars, err := svc.List()
			if err != nil {
				return err
			}
			enc := yaml.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent(2)
			defer enc.Close()
			return enc.Encode(chars)
		},
	}
}

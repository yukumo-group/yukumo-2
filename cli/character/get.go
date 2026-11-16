package character

import (
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newGetCmd(app *api.App, opts Options) *cobra.Command {
	return &cobra.Command{
		Use:   "get [name]",
		Short: "Print one character",
		Args: func(cmd *cobra.Command, args []string) error {
			if opts.Interactive {
				return cobra.MaximumNArgs(1)(cmd, args)
			}
			return cobra.ExactArgs(1)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name := ""
			if len(args) >= 1 {
				name = args[0]
			}
			if opts.Interactive {
				if err := ensureName(app, &name); err != nil {
					return err
				}
			}
			c, err := app.GetCharacter(name)
			if err != nil {
				return err
			}
			enc := yaml.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent(2)
			defer func() { _ = enc.Close() }()
			return enc.Encode(c)
		},
	}
}

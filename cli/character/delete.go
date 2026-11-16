package character

import (
	"errors"
	"fmt"

	"yukumo-script/cli/prompt"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func newDeleteCmd(app *api.App, opts Options) *cobra.Command {
	return &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete a character",
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
				ok := false
				if err := prompt.Confirm(fmt.Sprintf("Delete character %q?", name), &ok); err != nil {
					if errors.Is(err, prompt.ErrAborted) {
						return nil
					}
					return err
				}
				if !ok {
					return nil
				}
			}
			if err := app.DeleteCharacter(name); err != nil {
				return err
			}
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "deleted %q\n", name)
			return err
		},
	}
}

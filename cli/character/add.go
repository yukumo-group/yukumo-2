package character

import (
	"errors"
	"fmt"

	"yukumo-script/cli/prompt"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func newAddCmd(app *api.App, opts Options) *cobra.Command {
	return &cobra.Command{
		Use:   "add [name] [description]",
		Short: "Create a character",
		Args: func(cmd *cobra.Command, args []string) error {
			if opts.Interactive {
				return cobra.MaximumNArgs(2)(cmd, args)
			}
			return cobra.RangeArgs(1, 2)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c := api.Character{}
			if len(args) >= 1 {
				c.Name = args[0]
			}
			if len(args) >= 2 {
				c.Description = args[1]
			}
			if opts.Interactive {
				if err := prompt.EnsureString("Name", &c.Name); err != nil {
					if errors.Is(err, prompt.ErrAborted) {
						return nil
					}
					return err
				}
				if len(args) < 2 {
					if err := prompt.String("Description (optional)", &c.Description); err != nil {
						if errors.Is(err, prompt.ErrAborted) {
							return nil
						}
						return err
					}
				}
			}
			if err := app.AddCharacter(c); err != nil {
				return err
			}
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "added %q\n", c.Name)
			return err
		},
	}
}

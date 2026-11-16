package config

import (
	"fmt"

	"yukumo-script/cli/action"
	"yukumo-script/cli/internal/cobraio"
	"yukumo-script/cli/render"
	"yukumo-script/cli/streams"
	"yukumo-script/cli/ui"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func newEditCmd(app *api.App, opts Options) *cobra.Command {
	return &cobra.Command{
		Use:   "edit [key] [value]",
		Short: "Update a single setting",
		Args: func(cmd *cobra.Command, args []string) error {
			if opts.Interactive {
				return cobra.MaximumNArgs(2)(cmd, args)
			}
			return cobra.ExactArgs(2)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			s := cobraio.FromCmd(cmd)
			if opts.Interactive && len(args) < 2 {
				return ui.EditConfig(s, app, argAt(args, 0), argAt(args, 1))
			}
			return runOneShotEdit(s, app, args)
		},
	}
}

func runOneShotEdit(s streams.Streams, app *api.App, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("key and value are required")
	}
	result, err := action.ApplyConfigEdit(s.Ctx, app, args[0], args[1], true)
	if err != nil {
		return err
	}
	if result.Warning != "" {
		_ = render.Warning(s.Err, result.Warning)
	}
	return render.Line(s.Out, "%s: %s\n", args[0], result.Value)
}

func argAt(args []string, i int) string {
	if i < len(args) {
		return args[i]
	}
	return ""
}

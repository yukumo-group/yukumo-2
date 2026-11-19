package cli

import (
	"context"

	clicharacter "yukumo-script/cli/character"
	cliconfig "yukumo-script/cli/config"
	clisynth "yukumo-script/cli/synth"
	"yukumo-script/core/api"
	"yukumo-script/core/api/cliopen"

	"github.com/spf13/cobra"
)

func Execute() error {
	app, err := cliopen.Open(cliopen.Options{})
	if err != nil {
		return err
	}
	defer func() { _ = app.Close(context.Background()) }()

	return newRoot(app).ExecuteContext(context.Background())
}

func newRoot(app *api.App) *cobra.Command {
	root := &cobra.Command{
		Use:           "yukumo-script",
		Short:         "Yukumo script and voice tooling",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(cliconfig.NewCmd(app))
	root.AddCommand(clicharacter.NewCmd(app))
	root.AddCommand(clisynth.NewCmd(app))
	return root
}

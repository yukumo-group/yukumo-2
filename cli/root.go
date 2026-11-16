package cli

import (
	"context"

	clicharacter "yukumo-script/cli/character"
	cliconfig "yukumo-script/cli/config"
	"yukumo-script/cli/prompt"
	clisynth "yukumo-script/cli/synth"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func Execute() error {
	app, err := api.Open(api.Options{})
	if err != nil {
		return err
	}
	defer func() { _ = app.Close(context.Background()) }()

	cfg, err := app.Config()
	if err != nil {
		return err
	}
	prompt.SetAccessible(cfg.CLI.WithAccessibility)

	root := newRoot(app, cfg.CLI.DefaultInteractive)
	if cfg.CLI.DefaultInteractive {
		enableInteractiveBrowse(root)
	}
	return root.ExecuteContext(context.Background())
}

func newRoot(app *api.App, interactive bool) *cobra.Command {
	root := &cobra.Command{
		Use:           "yukumo-script",
		Short:         "Yukumo script and voice tooling",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(cliconfig.NewCmd(app, cliconfig.Options{Interactive: interactive}))
	root.AddCommand(clicharacter.NewCmd(app, clicharacter.Options{Interactive: interactive}))
	root.AddCommand(clisynth.NewCmd(app, clisynth.Options{Interactive: interactive}))
	return root
}

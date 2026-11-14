package cli

import (
	"fmt"
	"os"
	"path/filepath"

	cliconfig "yukumo-script/internal/cli/config"
	"yukumo-script/internal/cli/prompt"
	svcconfig "yukumo-script/internal/service/config"
	"yukumo-script/internal/store"

	"github.com/spf13/cobra"
)

func Execute() error {
	root, err := newRoot()
	if err != nil {
		return err
	}
	return root.Execute()
}

func newRoot() (*cobra.Command, error) {
	cfgDir, err := configDir()
	if err != nil {
		return nil, err
	}
	fsStore := store.NewFS(cfgDir)
	cfgSvc := svcconfig.New(fsStore)

	cfg, err := cfgSvc.Get()
	if err != nil {
		return nil, err
	}
	prompt.SetAccessible(cfg.CLI.Accessibility())

	interactive := cfg.CLI.Interactive()
	opts := cliconfig.Options{Interactive: interactive}
	root := &cobra.Command{
		Use:           "yukumo-script",
		Short:         "Yukumo script and voice tooling",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(cliconfig.NewCmd(cfgSvc, opts))
	if interactive {
		enableInteractiveBrowse(root)
	}
	return root, nil
}

func configDir() (string, error) {
	if v := os.Getenv("YUKUMO_CONFIG_DIR"); v != "" {
		return v, nil
	}
	base, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve config dir: %w", err)
	}
	return filepath.Join(base, "yukumo-script"), nil
}

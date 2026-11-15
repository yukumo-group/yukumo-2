package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	clicharacter "yukumo-script/internal/cli/character"
	cliconfig "yukumo-script/internal/cli/config"
	"yukumo-script/internal/cli/prompt"
	"yukumo-script/internal/runtime/aq"
	svccharacter "yukumo-script/internal/service/character"
	svcvoice "yukumo-script/internal/service/character/voice"
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
	cfgSvc.EnsurePaths()
	prompt.SetAccessible(cfg.CLI.Accessibility())

	aqPool := aq.NewPool()

	interactive := cfg.CLI.Interactive()
	cfgOpts := cliconfig.Options{Interactive: interactive}
	charOpts := clicharacter.Options{Interactive: interactive}
	charSvc := svccharacter.New(cfgSvc)
	voiceSvc := svcvoice.New(charSvc, aqPool)
	root := &cobra.Command{
		Use:           "yukumo-script",
		Short:         "Yukumo script and voice tooling",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.PersistentPostRun = func(*cobra.Command, []string) {
		_ = aqPool.Close(context.Background())
	}
	root.AddCommand(cliconfig.NewCmd(cfgSvc, cfgOpts))
	root.AddCommand(clicharacter.NewCmd(charSvc, voiceSvc, charOpts))
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

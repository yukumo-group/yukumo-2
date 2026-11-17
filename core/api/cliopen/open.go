package cliopen

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"yukumo-script/core/api"
	"yukumo-script/core/internal/service/config"
	"yukumo-script/core/internal/store/fs"
	"yukumo-script/runtime/aq"
)

// Options controls how the filesystem-backed application is opened.
// CharacterPath and OutputPath are optional overrides for tests/embedding.
type Options struct {
	ConfigDir     string
	CharacterPath string
	OutputPath    string
}

// Open is the CLI platform composition root: FS stores, OS paths and aq.Pool.
func Open(opts Options) (*api.App, error) {
	dir, err := resolveConfigDir(opts.ConfigDir)
	if err != nil {
		return nil, err
	}
	cfgStore := fs.New(dir)
	defaults, err := domainDefaults()
	if err != nil {
		return nil, err
	}
	cfgSvc := config.New(cfgStore, defaults)
	cfg, err := cfgSvc.Get()
	if err != nil {
		return nil, err
	}

	if opts.CharacterPath != "" {
		cfg.Paths.CharacterPath = opts.CharacterPath
	}
	if opts.OutputPath != "" {
		cfg.Paths.OutputPath = opts.OutputPath
	}
	charStore := fs.New(cfg.Paths.CharacterPath)
	artifacts := fs.New(cfg.Paths.OutputPath)
	outputs := fs.Output{}

	pool := aq.NewPool()
	engine := aqEngine{pool: pool}
	app, err := api.New(api.Deps{
		ConfigStore:    cfgStore,
		CharacterStore: charStore,
		Artifacts:      artifacts,
		Outputs:        outputs,
		Engine:         engine,
		Defaults:       apiDefaults(defaults),
	})
	if err != nil {
		_ = pool.Close(context.Background())
		return nil, err
	}
	if err := ensurePaths(cfg.Paths.CharacterPath, cfg.Paths.OutputPath); err != nil {
		_ = app.Close(context.Background())
		return nil, err
	}
	if err := app.CleanIndex(); err != nil {
		_ = app.Close(context.Background())
		return nil, fmt.Errorf("clean synth index: %w", err)
	}
	return app, nil
}

// EnsurePaths creates the configured character and output directories.
func EnsurePaths(app *api.App) error {
	cfg, err := app.Config()
	if err != nil {
		return err
	}
	return ensurePaths(cfg.Paths.CharacterPath, cfg.Paths.OutputPath)
}

func ensurePaths(paths ...string) error {
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if err := os.MkdirAll(p, 0o755); err != nil {
			return fmt.Errorf("create config path %q: %w", p, err)
		}
	}
	return nil
}

func resolveConfigDir(override string) (string, error) {
	if override != "" {
		return override, nil
	}
	if v := os.Getenv("YUKUMO_CONFIG_DIR"); v != "" {
		return v, nil
	}
	base, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve config dir: %w", err)
	}
	return filepath.Join(base, "yukumo-script"), nil
}

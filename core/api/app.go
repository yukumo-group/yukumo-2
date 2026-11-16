package api

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"yukumo-script/core/internal/service/character"
	svcvoice "yukumo-script/core/internal/service/character/voice"
	"yukumo-script/core/internal/service/config"
	"yukumo-script/core/internal/service/synth"
	"yukumo-script/core/internal/store"
	"yukumo-script/runtime/aq"
)

// Options controls how the application is opened.
type Options struct {
	ConfigDir string
}

// App is the use-case facade used by CLI and future GUIs.
type App struct {
	config *config.Service
	chars  *character.Service
	voices *svcvoice.Service
	synth  *synth.Service
	pool   *aq.Pool
}

// Open wires the filesystem-backed application. A future web adapter can add
// its own constructor without changing the use-case methods.
func Open(opts Options) (*App, error) {
	dir, err := resolveConfigDir(opts.ConfigDir)
	if err != nil {
		return nil, err
	}
	cfgSvc := config.New(store.NewFS(dir))
	cfgSvc.EnsurePaths()

	pool := aq.NewPool()
	charSvc := character.New(cfgSvc)
	voiceSvc := svcvoice.New(charSvc, pool)
	synthSvc := synth.New(cfgSvc, voiceSvc, aqEngine{pool: pool})
	if err := synthSvc.CleanIndex(); err != nil {
		_ = pool.Close(context.Background())
		return nil, fmt.Errorf("clean synth index: %w", err)
	}
	return &App{config: cfgSvc, chars: charSvc, voices: voiceSvc, synth: synthSvc, pool: pool}, nil
}

func (a *App) Close(ctx context.Context) error {
	if a == nil || a.pool == nil {
		return nil
	}
	return a.pool.Close(ctx)
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

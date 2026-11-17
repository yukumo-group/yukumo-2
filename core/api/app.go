package api

import (
	"context"
	"errors"
	"fmt"

	"yukumo-script/core/internal/service/character"
	svcvoice "yukumo-script/core/internal/service/character/voice"
	"yukumo-script/core/internal/service/config"
	"yukumo-script/core/internal/service/synth"
	"yukumo-script/core/port"
)

// Deps are the host-supplied ports. ConfigStore and CharacterStore are
// required; Artifacts, Outputs and Engine are optional.
type Deps struct {
	ConfigStore    port.Store
	CharacterStore port.Store
	Artifacts      port.ArtifactStore
	Outputs        port.OutputWriter
	Engine         port.Engine
	Defaults       Config
}

// App is the use-case facade used by CLI, wasm hosts and future GUIs.
type App struct {
	config *config.Service
	chars  *character.Service
	voices *svcvoice.Service
	synth  *synth.Service
	engine port.Engine
	stores []port.Store
}

// New wires the use cases around host-supplied ports. It performs no OS IO.
func New(deps Deps) (*App, error) {
	if deps.ConfigStore == nil {
		return nil, fmt.Errorf("api: config store is required")
	}
	if deps.CharacterStore == nil {
		return nil, fmt.Errorf("api: character store is required")
	}
	cfgSvc := config.New(deps.ConfigStore, fromConfig(deps.Defaults))
	charSvc := character.New(cfgSvc, deps.CharacterStore)
	voiceSvc := svcvoice.New(charSvc, deps.Engine)
	synthSvc := synth.New(cfgSvc, voiceSvc, deps.Engine, deps.Artifacts, deps.Outputs)
	app := &App{
		config: cfgSvc,
		chars:  charSvc,
		voices: voiceSvc,
		synth:  synthSvc,
		engine: deps.Engine,
		stores: []port.Store{deps.ConfigStore, deps.CharacterStore},
	}
	if deps.Artifacts != nil {
		app.stores = append(app.stores, deps.Artifacts)
	}
	return app, nil
}

// CleanIndex drops synth cache entries whose artifact is missing.
func (a *App) CleanIndex() error {
	if a == nil || a.synth == nil {
		return nil
	}
	return a.synth.CleanIndex()
}

// Close closes the engine and flushes stores that implement port.Flusher.
func (a *App) Close(ctx context.Context) error {
	if a == nil {
		return nil
	}
	var errs []error
	if c, ok := a.engine.(port.Closer); ok && c != nil {
		errs = append(errs, c.Close(ctx))
	}
	for _, s := range a.stores {
		if f, ok := s.(port.Flusher); ok {
			errs = append(errs, f.Flush(ctx))
		}
	}
	return errors.Join(errs...)
}

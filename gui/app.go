package main

import (
	"context"
	"errors"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
	coreapi "yukumo-script/core/api"
	"yukumo-script/core/api/cliopen"
)

// App is the Wails transport adapter over core/api. The same method surface is
// what a future wasm host is expected to expose to the frontend.
type App struct {
	mu          sync.RWMutex
	opMu        sync.Mutex
	core        *coreapi.App
	player      *wavPlayer
	openOptions cliopen.Options
}

// ServiceStartup opens the OS-backed core application.
func (a *App) ServiceStartup(context.Context, application.ServiceOptions) error {
	core, err := cliopen.Open(a.openOptions)
	if err != nil {
		return err
	}
	a.mu.Lock()
	a.core = core
	a.player = newWAVPlayer()
	a.mu.Unlock()
	return nil
}

// ServiceShutdown closes the engine and flushes stores.
func (a *App) ServiceShutdown() error {
	a.opMu.Lock()
	defer a.opMu.Unlock()
	a.mu.Lock()
	core := a.core
	player := a.player
	a.core = nil
	a.player = nil
	a.mu.Unlock()
	if player != nil {
		player.Stop()
	}
	if core == nil {
		return nil
	}
	return core.Close(context.Background())
}

func (a *App) current() (*coreapi.App, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if a.core == nil {
		return nil, errors.New("yukumo backend is not started")
	}
	return a.core, nil
}

func (a *App) mutate(fn func(*coreapi.App) error) error {
	core, err := a.current()
	if err != nil {
		return err
	}
	a.opMu.Lock()
	defer a.opMu.Unlock()
	return fn(core)
}

func (a *App) ListCharacters() ([]coreapi.Character, error) {
	core, err := a.current()
	if err != nil {
		return nil, err
	}
	return core.ListCharacters()
}

func (a *App) GetCharacter(name string) (coreapi.Character, error) {
	core, err := a.current()
	if err != nil {
		return coreapi.Character{}, err
	}
	return core.GetCharacter(name)
}

func (a *App) AddCharacter(c coreapi.Character) error {
	return a.mutate(func(core *coreapi.App) error { return core.AddCharacter(c) })
}

func (a *App) UpdateCharacter(oldName string, c coreapi.Character) error {
	return a.mutate(func(core *coreapi.App) error { return core.UpdateCharacter(oldName, c) })
}

func (a *App) DeleteCharacter(name string) error {
	return a.mutate(func(core *coreapi.App) error { return core.DeleteCharacter(name) })
}

func (a *App) ListVoices(charName string) ([]coreapi.Voice, error) {
	core, err := a.current()
	if err != nil {
		return nil, err
	}
	return core.ListVoices(charName)
}

func (a *App) GetVoice(charName, voiceName string) (coreapi.Voice, error) {
	core, err := a.current()
	if err != nil {
		return coreapi.Voice{}, err
	}
	return core.GetVoice(charName, voiceName)
}

func (a *App) AddVoice(ctx context.Context, charName string, v coreapi.Voice) error {
	return a.mutate(func(core *coreapi.App) error { return core.AddVoice(ctx, charName, v) })
}

func (a *App) UpdateVoice(ctx context.Context, charName, oldName string, v coreapi.Voice) error {
	return a.mutate(func(core *coreapi.App) error {
		return core.UpdateVoice(ctx, charName, oldName, v)
	})
}

func (a *App) DeleteVoice(charName, voiceName string) error {
	return a.mutate(func(core *coreapi.App) error { return core.DeleteVoice(charName, voiceName) })
}

func (a *App) DuplicateVoice(ctx context.Context, charName, srcName, newName string) (string, error) {
	var out string
	err := a.mutate(func(core *coreapi.App) error {
		name, err := core.DuplicateVoice(ctx, charName, srcName, newName)
		out = name
		return err
	})
	return out, err
}

func (a *App) Config() (coreapi.Config, error) {
	core, err := a.current()
	if err != nil {
		return coreapi.Config{}, err
	}
	return core.Config()
}

func (a *App) SaveConfig(c coreapi.Config) error {
	return a.mutate(func(core *coreapi.App) error { return core.SaveConfig(c) })
}

func (a *App) PresetNames(ctx context.Context, version int) ([]string, error) {
	core, err := a.current()
	if err != nil {
		return nil, err
	}
	return core.PresetNames(ctx, version)
}

func (a *App) HasPreset(ctx context.Context, version, preset int) (bool, error) {
	core, err := a.current()
	if err != nil {
		return false, err
	}
	return core.HasPreset(ctx, version, preset)
}

// SynthInput is the GUI's bytes-only synthesis request. The transport never
// writes WAV files: the frontend plays or downloads the returned bytes.
type SynthInput struct {
	Voice       coreapi.Voice `json:"voice"`
	Text        string        `json:"text"`
	Version     int           `json:"version"`
	SpeedFactor float64       `json:"speedFactor"`
}

// PlayWAV plays a preview through a native OS audio player.
func (a *App) PlayWAV(data []byte) error {
	a.mu.RLock()
	player := a.player
	a.mu.RUnlock()
	if player == nil {
		return errors.New("yukumo backend is not started")
	}
	return player.Play(data)
}

// StopPlayback stops the current preview.
func (a *App) StopPlayback() {
	a.mu.RLock()
	player := a.player
	a.mu.RUnlock()
	if player != nil {
		player.Stop()
	}
}

func (a *App) Synthesize(ctx context.Context, input SynthInput) ([]byte, error) {
	core, err := a.current()
	if err != nil {
		return nil, err
	}
	res, err := core.Synthesize(ctx, coreapi.SynthRequest{
		Voice:       &input.Voice,
		VoiceName:   input.Voice.Name,
		Text:        input.Text,
		Version:     input.Version,
		SpeedFactor: input.SpeedFactor,
		NoWrite:     true,
	})
	if err != nil {
		return nil, err
	}
	return res.WAV, nil
}

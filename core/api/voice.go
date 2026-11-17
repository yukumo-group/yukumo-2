package api

import (
	"context"
	"fmt"
)

func (a *App) ListVoices(charName string) ([]Voice, error) {
	list, err := a.voices.List(charName)
	if err != nil {
		return nil, err
	}
	out := make([]Voice, len(list))
	for i, v := range list {
		out[i] = toVoice(v)
	}
	return out, nil
}

func (a *App) GetVoice(charName, voiceName string) (Voice, error) {
	v, err := a.voices.Get(charName, voiceName)
	if err != nil {
		return Voice{}, err
	}
	return toVoice(v), nil
}

func (a *App) AddVoice(ctx context.Context, charName string, v Voice) error {
	return a.voices.Add(ctx, charName, fromVoice(v))
}

func (a *App) UpdateVoice(ctx context.Context, charName, oldName string, v Voice) error {
	return a.voices.Update(ctx, charName, oldName, fromVoice(v))
}

func (a *App) DeleteVoice(charName, voiceName string) error {
	return a.voices.Delete(charName, voiceName)
}

func (a *App) DuplicateVoice(ctx context.Context, charName, srcName, newName string) (string, error) {
	return a.voices.Duplicate(ctx, charName, srcName, newName)
}

func (a *App) PresetNames(ctx context.Context, version int) ([]string, error) {
	catalog := a.voices.Presets()
	if catalog == nil {
		return nil, fmt.Errorf("%w: preset catalog unavailable", ErrInvalid)
	}
	return catalog.PresetNames(ctx, int32(version))
}

func (a *App) HasPreset(ctx context.Context, version, preset int) (bool, error) {
	catalog := a.voices.Presets()
	if catalog == nil {
		return false, fmt.Errorf("%w: preset catalog unavailable", ErrInvalid)
	}
	return catalog.HasPreset(ctx, int32(version), int32(preset))
}

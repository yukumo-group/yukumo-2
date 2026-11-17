//go:build js

package main

import (
	"context"

	"yukumo-script/core/port"
)

type fakeEngine struct{}

func (fakeEngine) Synthesize(context.Context, string, port.EngineParam) ([]byte, error) {
	wav := make([]byte, 44)
	copy(wav[0:4], "RIFF")
	copy(wav[8:12], "WAVE")
	copy(wav[12:16], "fmt ")
	copy(wav[36:40], "data")
	return wav, nil
}

func (fakeEngine) HasPreset(context.Context, int32, int32) (bool, error) { return true, nil }

func (fakeEngine) PresetNames(context.Context, int32) ([]string, error) {
	return []string{"fake"}, nil
}

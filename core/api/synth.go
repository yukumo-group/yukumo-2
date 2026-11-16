package api

import (
	"context"

	"yukumo-script/core/internal/service/synth"
)

type SynthRequest struct {
	VoiceRef      string
	VoicePathSpec string
	Text          string
	Version       int
	SpeedFactor   float64
	Output        string
	Preview       bool
	ReuseCache    *bool
}

type SynthResult struct {
	WAV       []byte
	Path      string
	Wrote     bool
	FromCache bool
}

func (a *App) Synthesize(ctx context.Context, req SynthRequest) (SynthResult, error) {
	res, err := a.synth.Synthesize(ctx, synth.Request{
		VoiceRef:      req.VoiceRef,
		VoicePathSpec: req.VoicePathSpec,
		Text:          req.Text,
		Version:       req.Version,
		SpeedFactor:   req.SpeedFactor,
		Output:        req.Output,
		Preview:       req.Preview,
		ReuseCache:    req.ReuseCache,
	})
	if err != nil {
		return SynthResult{}, err
	}
	return SynthResult{WAV: res.WAV, Path: res.Path, Wrote: res.Wrote, FromCache: res.FromCache}, nil
}

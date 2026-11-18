package api

import (
	"context"

	"yukumo-script/core/internal/domain"
	"yukumo-script/core/internal/service/synth"
)

type SynthRequest struct {
	VoiceRef    string  `json:"voiceRef,omitempty"`  // character:voice (mutually exclusive with Voice/VoiceFile)
	VoiceFile   []byte  `json:"voiceFile,omitempty"` // raw Character YAML bytes
	Voice       *Voice  `json:"voice,omitempty"`     // inline voice (mutually exclusive with VoiceRef/VoiceFile)
	VoiceName   string  `json:"voiceName,omitempty"` // voice name fallback for VoiceFile / inline Voice
	VoiceBase   string  `json:"voiceBase,omitempty"` // fallback character name when the voice has none
	Text        string  `json:"text"`
	Version     int     `json:"version"`
	SpeedFactor float64 `json:"speedFactor"`
	Output      string  `json:"output,omitempty"`
	Preview     bool    `json:"preview"`
	NoWrite     bool    `json:"noWrite"`              // return WAV bytes without writing output or cache
	ReuseCache  *bool   `json:"reuseCache,omitempty"` // nil means use the config value
}

type SynthResult struct {
	WAV       []byte
	Path      string
	Wrote     bool
	FromCache bool
}

func (a *App) Synthesize(ctx context.Context, req SynthRequest) (SynthResult, error) {
	res, err := a.synth.Synthesize(ctx, synth.Request{
		VoiceRef:    req.VoiceRef,
		VoiceFile:   req.VoiceFile,
		Voice:       fromVoicePtr(req.Voice),
		VoiceName:   req.VoiceName,
		VoiceBase:   req.VoiceBase,
		Text:        req.Text,
		Version:     req.Version,
		SpeedFactor: req.SpeedFactor,
		Output:      req.Output,
		Preview:     req.Preview,
		NoWrite:     req.NoWrite,
		ReuseCache:  req.ReuseCache,
	})
	if err != nil {
		return SynthResult{}, err
	}
	return SynthResult{WAV: res.WAV, Path: res.Path, Wrote: res.Wrote, FromCache: res.FromCache}, nil
}

func fromVoicePtr(v *Voice) *domain.Voice {
	if v == nil {
		return nil
	}
	out := fromVoice(*v)
	return &out
}

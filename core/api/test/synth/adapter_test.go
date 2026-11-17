package synth_test

import (
	"context"

	"yukumo-script/core/port"
	"yukumo-script/runtime/aq"
)

type aqEngine struct {
	pool *aq.Pool
}

func (e aqEngine) Synthesize(ctx context.Context, text string, p port.EngineParam) ([]byte, error) {
	return e.pool.Synthesize(ctx, text, aq.Param{
		Version: int32(p.Version),
		Preset:  int32(p.Preset),
		Speed:   int32(p.Speed),
		Phont:   p.Phont,
		Base:    int32(p.Base),
		Volume:  int32(p.Volume),
		Pitch:   int32(p.Pitch),
		Accent:  int32(p.Accent),
		Lmd:     int32(p.Lmd),
		Fsc:     int32(p.Fsc),
	})
}

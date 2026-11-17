package port

import "context"

// EngineParam mirrors third-party/aq_wasm.h AqSynthParam.
type EngineParam struct {
	Version int
	Preset  int
	Speed   int
	Phont   []byte
	Base    int
	Volume  int
	Pitch   int
	Accent  int
	Lmd     int
	Fsc     int
}

// Synthesizer is the minimal engine surface used by synth.
type Synthesizer interface {
	Synthesize(ctx context.Context, text string, p EngineParam) ([]byte, error)
}

// PresetCatalog is the minimal engine surface used by voice validation.
type PresetCatalog interface {
	HasPreset(ctx context.Context, version, preset int32) (bool, error)
	PresetNames(ctx context.Context, version int32) ([]string, error)
}

// Engine is the full host engine contract.
type Engine interface {
	Synthesizer
	PresetCatalog
}

// Closer is optionally implemented by engines that own resources.
type Closer interface {
	Close(ctx context.Context) error
}

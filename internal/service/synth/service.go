package synth

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"yukumo-script/internal/domain"
	"yukumo-script/internal/runtime/aq"
)

type Config interface {
	Get() (domain.Config, error)
}

type Voices interface {
	Get(charName, voiceName string) (domain.Voice, error)
}

type Engine interface {
	Synthesize(ctx context.Context, text string, param aq.Param) ([]byte, error)
}

type Service struct {
	config Config
	voices Voices
	engine Engine
}

func New(cfg Config, voices Voices, engine Engine) *Service {
	return &Service{config: cfg, voices: voices, engine: engine}
}

type Request struct {
	VoiceRef    string  // "character.voice"
	Text        string
	Version     int     // 0 = auto 1→2→10
	SpeedFactor float64 // default 1.0; multiplies voice param speed
	Output      string  // explicit -o path
	Preview     bool
}

type Result struct {
	WAV   []byte
	Path  string // set when a file was written
	Wrote bool
}

func (s *Service) Synthesize(ctx context.Context, req Request) (Result, error) {
	text := strings.TrimSpace(req.Text)
	if text == "" {
		return Result{}, fmt.Errorf("%w: text is required", domain.ErrInvalid)
	}
	charName, voiceName, err := ParseVoiceRef(req.VoiceRef)
	if err != nil {
		return Result{}, err
	}
	v, err := s.voices.Get(charName, voiceName)
	if err != nil {
		return Result{}, err
	}
	p, err := SelectParam(v, req.Version)
	if err != nil {
		return Result{}, err
	}

	factor := req.SpeedFactor
	if factor <= 0 {
		factor = 1
	}
	speed := effectiveSpeed(p.Speed, factor)
	wav, err := s.engine.Synthesize(ctx, text, toAQParam(p, speed))
	if err != nil {
		return Result{}, err
	}

	cfg, err := s.config.Get()
	if err != nil {
		return Result{}, err
	}

	out := Result{WAV: wav}
	path, write, err := s.resolveOutputPath(cfg, req, charName, voiceName)
	if err != nil {
		return Result{}, err
	}
	if !write {
		return out, nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return Result{}, err
	}
	if err := os.WriteFile(path, wav, 0o644); err != nil {
		return Result{}, err
	}
	out.Path = path
	out.Wrote = true
	return out, nil
}

func (s *Service) resolveOutputPath(cfg domain.Config, req Request, charName, voiceName string) (path string, write bool, err error) {
	if strings.TrimSpace(req.Output) != "" {
		return req.Output, true, nil
	}
	if !cfg.Synth.WriteAlways() && req.Preview {
		return "", false, nil
	}
	dir := strings.TrimSpace(cfg.Paths.OutputPath)
	if dir == "" {
		dir = domain.DefaultOutputPath()
	}
	format := strings.TrimSpace(cfg.Synth.OutputFormat)
	if format == "" {
		format = domain.DefaultOutputFormat
	}
	base := FormatOutputName(format, charName, voiceName)
	path, err = UniquePath(dir, base, ".wav")
	return path, true, err
}

func effectiveSpeed(base int, factor float64) int {
	if base == 0 {
		base = 100
	}
	v := int(math.Round(float64(base) * factor))
	if v < 50 {
		return 50
	}
	if v > 300 {
		return 300
	}
	return v
}

func toAQParam(p domain.SynthParam, speed int) aq.Param {
	return aq.Param{
		Version: int32(p.Version),
		Preset:  int32(p.Preset),
		Speed:   int32(speed),
		Phont:   p.PhontData,
		Base:    int32(p.Base),
		Volume:  int32(p.Volume),
		Pitch:   int32(p.Pitch),
		Accent:  int32(p.Accent),
		Lmd:     int32(p.Lmd),
		Fsc:     int32(p.Fsc),
	}
}

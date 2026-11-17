package synth

import (
	"context"
	"fmt"
	"math"
	"strings"

	"yukumo-script/core/internal/domain"
	"yukumo-script/core/port"
)

type Config interface {
	Get() (domain.Config, error)
}

type Voices interface {
	Get(charName, voiceName string) (domain.Voice, error)
}

type Service struct {
	config    Config
	voices    Voices
	engine    port.Synthesizer
	artifacts port.ArtifactStore
	outputs   port.OutputWriter
}

func New(cfg Config, voices Voices, engine port.Synthesizer, artifacts port.ArtifactStore, outputs port.OutputWriter) *Service {
	return &Service{config: cfg, voices: voices, engine: engine, artifacts: artifacts, outputs: outputs}
}

type Request struct {
	VoiceRef    string // "character:voice" (mutually exclusive with VoiceFile)
	VoiceFile   []byte // raw Character YAML bytes
	VoiceName   string // required with VoiceFile
	VoiceBase   string // fallback character name when the YAML has none
	Text        string
	Version     int     // 0 = auto 1→2→10
	SpeedFactor float64 // default 1.0
	Output      string  // explicit -o path
	Preview     bool
	ReuseCache  *bool // nil → use config synth.reuse_cache
}

type Result struct {
	WAV       []byte
	Path      string
	Wrote     bool
	FromCache bool
}

func (s *Service) CleanIndex() error {
	if s.artifacts == nil {
		return nil
	}
	return cleanIndex(s.artifacts)
}

func (s *Service) Synthesize(ctx context.Context, req Request) (Result, error) {
	text := strings.TrimSpace(req.Text)
	if text == "" {
		return Result{}, fmt.Errorf("%w: text is required", domain.ErrInvalid)
	}
	hasRef := strings.TrimSpace(req.VoiceRef) != ""
	hasFile := req.VoiceFile != nil
	if hasRef == hasFile {
		return Result{}, fmt.Errorf("%w: provide exactly one of voice or voice file", domain.ErrInvalid)
	}

	var (
		charName, voiceName string
		voice               domain.Voice
		fileBytes           []byte
	)
	if hasFile {
		c, v, raw, err := LoadVoiceBytes(req.VoiceFile, req.VoiceName, req.VoiceBase)
		if err != nil {
			return Result{}, err
		}
		charName, voiceName, voice, fileBytes = c.Name, v.Name, v, raw
	} else {
		var err error
		charName, voiceName, err = ParseVoiceRef(req.VoiceRef)
		if err != nil {
			return Result{}, err
		}
		voice, err = s.voices.Get(charName, voiceName)
		if err != nil {
			return Result{}, err
		}
	}

	p, err := SelectParam(voice, req.Version)
	if err != nil {
		return Result{}, err
	}
	factor := req.SpeedFactor
	if factor <= 0 {
		factor = 1
	}
	speed := effectiveSpeed(p.Speed, factor)
	hash := TaskHash(text, speed, p, fileBytes)

	cfg, err := s.config.Get()
	if err != nil {
		return Result{}, err
	}
	reuse := cfg.Synth.CacheReuse()
	if req.ReuseCache != nil {
		reuse = *req.ReuseCache
	}

	target, err := s.resolveTarget(cfg, req, charName, voiceName)
	if err != nil {
		return Result{}, err
	}

	if reuse && s.artifacts != nil {
		key, ok, err := lookupCache(s.artifacts, hash)
		if err != nil {
			return Result{}, err
		}
		if ok {
			data, err := s.artifacts.Read(key)
			if err != nil {
				_ = putCacheRemove(s.artifacts, hash)
			} else if !target.write {
				return Result{WAV: data, FromCache: true}, nil
			} else {
				if err := s.writeTarget(target, data); err != nil {
					return Result{}, err
				}
				return Result{WAV: data, Path: s.targetPath(target), Wrote: true, FromCache: true}, nil
			}
		}
	}

	if s.engine == nil {
		return Result{}, fmt.Errorf("%w: engine unavailable", domain.ErrInvalid)
	}
	wav, err := s.engine.Synthesize(ctx, text, toEngineParam(p, speed))
	if err != nil {
		return Result{}, err
	}
	out := Result{WAV: wav}
	if !target.write {
		return out, nil
	}
	if err := s.writeTarget(target, wav); err != nil {
		return Result{}, err
	}
	out.Path = s.targetPath(target)
	out.Wrote = true
	if err := putCache(s.artifacts, hash, target.key, target.external != "", wav); err != nil {
		return out, err
	}
	return out, nil
}

type target struct {
	external string // -o path; empty means artifact key
	key      string
	write    bool
}

func (s *Service) resolveTarget(cfg domain.Config, req Request, charName, voiceName string) (target, error) {
	if strings.TrimSpace(req.Output) != "" {
		return target{external: req.Output, write: true}, nil
	}
	if !cfg.Synth.WriteAlways() && req.Preview {
		return target{}, nil
	}
	if s.artifacts == nil {
		return target{}, fmt.Errorf("%w: artifact store unavailable", domain.ErrInvalid)
	}
	format := strings.TrimSpace(cfg.Synth.OutputFormat)
	if format == "" {
		format = domain.DefaultOutputFormat
	}
	base := FormatOutputName(format, charName, voiceName)
	key, err := uniqueKey(s.artifacts, base, ".wav")
	if err != nil {
		return target{}, err
	}
	return target{key: key, write: true}, nil
}

func (s *Service) writeTarget(t target, data []byte) error {
	if t.external != "" {
		if s.outputs == nil {
			return fmt.Errorf("%w: output writer unavailable", domain.ErrInvalid)
		}
		return s.outputs.Write(t.external, data)
	}
	return s.artifacts.Write(t.key, data)
}

func (s *Service) targetPath(t target) string {
	if t.external != "" {
		return t.external
	}
	return s.artifacts.Path(t.key)
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

func toEngineParam(p domain.SynthParam, speed int) port.EngineParam {
	return port.EngineParam{
		Version: p.Version,
		Preset:  p.Preset,
		Speed:   speed,
		Phont:   p.PhontData,
		Base:    p.Base,
		Volume:  p.Volume,
		Pitch:   p.Pitch,
		Accent:  p.Accent,
		Lmd:     p.Lmd,
		Fsc:     p.Fsc,
	}
}

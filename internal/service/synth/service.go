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
	VoiceRef      string  // "character:voice" (mutually exclusive with VoicePathSpec)
	VoicePathSpec string  // "/path/character.yaml:voice"
	Text          string
	Version       int     // 0 = auto 1→2→10
	SpeedFactor   float64 // default 1.0
	Output        string  // explicit -o path
	Preview       bool
	ReuseCache    *bool // nil → use config synth.reuse_cache
}

type Result struct {
	WAV     []byte
	Path    string
	Wrote   bool
	FromCache bool
}

func (s *Service) CleanIndex() error {
	cfg, err := s.config.Get()
	if err != nil {
		return err
	}
	return CleanIndex(outputDir(cfg))
}

func (s *Service) Synthesize(ctx context.Context, req Request) (Result, error) {
	text := strings.TrimSpace(req.Text)
	if text == "" {
		return Result{}, fmt.Errorf("%w: text is required", domain.ErrInvalid)
	}
	hasRef := strings.TrimSpace(req.VoiceRef) != ""
	hasPath := strings.TrimSpace(req.VoicePathSpec) != ""
	if hasRef == hasPath {
		return Result{}, fmt.Errorf("%w: provide exactly one of --voice or --voice-path", domain.ErrInvalid)
	}

	var (
		charName, voiceName string
		voice               domain.Voice
		fileBytes           []byte
	)
	if hasPath {
		yamlPath, vn, err := ParseVoicePathSpec(req.VoicePathSpec)
		if err != nil {
			return Result{}, err
		}
		c, v, raw, err := LoadVoiceFile(yamlPath, vn)
		if err != nil {
			return Result{}, err
		}
		charName, voiceName, voice, fileBytes = c.Name, v.Name, v, raw
		if charName == "" {
			charName = baseNameNoExt(yamlPath)
		}
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
	dir := outputDir(cfg)
	reuse := cfg.Synth.CacheReuse()
	if req.ReuseCache != nil {
		reuse = *req.ReuseCache
	}

	dest, write, err := s.resolveOutputPath(cfg, req, charName, voiceName)
	if err != nil {
		return Result{}, err
	}

	if reuse {
		cached, ok, err := lookupCache(dir, hash)
		if err != nil {
			return Result{}, err
		}
		if ok {
			data, err := os.ReadFile(cached)
			if err != nil {
				putCacheRemove(dir, hash) // best-effort; lookup already cleaned missing
				// fall through to synthesize
			} else if !write {
				return Result{WAV: data, FromCache: true}, nil
			} else {
				if err := copyFile(cached, dest); err != nil {
					return Result{}, err
				}
				return Result{WAV: data, Path: dest, Wrote: true, FromCache: true}, nil
			}
		}
	}

	wav, err := s.engine.Synthesize(ctx, text, toAQParam(p, speed))
	if err != nil {
		return Result{}, err
	}
	out := Result{WAV: wav}
	if !write {
		return out, nil
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return Result{}, err
	}
	if err := os.WriteFile(dest, wav, 0o644); err != nil {
		return Result{}, err
	}
	out.Path = dest
	out.Wrote = true
	if err := putCache(dir, hash, filepath.Base(dest)); err != nil {
		return out, err
	}
	return out, nil
}

func putCacheRemove(outputDir, hash string) {
	idx, err := loadIndex(outputDir)
	if err != nil {
		return
	}
	delete(idx, hash)
	_ = saveIndex(outputDir, idx)
}

func (s *Service) resolveOutputPath(cfg domain.Config, req Request, charName, voiceName string) (path string, write bool, err error) {
	if strings.TrimSpace(req.Output) != "" {
		return req.Output, true, nil
	}
	if !cfg.Synth.WriteAlways() && req.Preview {
		return "", false, nil
	}
	dir := outputDir(cfg)
	format := strings.TrimSpace(cfg.Synth.OutputFormat)
	if format == "" {
		format = domain.DefaultOutputFormat
	}
	base := FormatOutputName(format, charName, voiceName)
	path, err = UniquePath(dir, base, ".wav")
	return path, true, err
}

func outputDir(cfg domain.Config) string {
	dir := strings.TrimSpace(cfg.Paths.OutputPath)
	if dir == "" {
		return domain.DefaultOutputPath()
	}
	return dir
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

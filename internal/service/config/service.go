package config

import (
	"fmt"
	"os"
	"strings"

	"yukumo-script/internal/domain"
	"yukumo-script/internal/store"
)

const fileName = "config.yaml"

type Service struct {
	store store.Store
}

func New(s store.Store) *Service {
	return &Service{store: s}
}

func (s *Service) Get() (domain.Config, error) {
	cfg, err := store.LoadYAMLOrDefault(s.store, fileName, domain.DefaultConfig())
	if err != nil {
		return cfg, err
	}
	applyPathDefaults(&cfg)
	return cfg, nil
}

func (s *Service) Save(cfg domain.Config) error {
	applyPathDefaults(&cfg)
	ensurePaths(cfg)
	return store.SaveYAML(s.store, fileName, cfg)
}

func (s *Service) Edit(mutate func(*domain.Config)) (domain.Config, error) {
	cfg, err := s.Get()
	if err != nil {
		return domain.Config{}, err
	}
	mutate(&cfg)
	if err := s.Save(cfg); err != nil {
		return domain.Config{}, err
	}
	return cfg, nil
}

// EnsurePaths creates configured directories; panics if creation fails (e.g. permission denied).
func (s *Service) EnsurePaths() {
	cfg, err := s.Get()
	if err != nil {
		panic(fmt.Sprintf("load config for path ensure: %v", err))
	}
	ensurePaths(cfg)
}

func applyPathDefaults(cfg *domain.Config) {
	if strings.TrimSpace(cfg.Paths.CharacterPath) == "" {
		cfg.Paths.CharacterPath = domain.DefaultCharacterPath()
	}
	if strings.TrimSpace(cfg.Paths.OutputPath) == "" {
		cfg.Paths.OutputPath = domain.DefaultOutputPath()
	}
	if strings.TrimSpace(cfg.Synth.OutputFormat) == "" {
		cfg.Synth.OutputFormat = domain.DefaultOutputFormat
	}
	if cfg.Synth.AlwaysWrite == nil {
		on := true
		cfg.Synth.AlwaysWrite = &on
	}
}

func ensurePaths(cfg domain.Config) {
	for _, p := range []string{cfg.Paths.CharacterPath, cfg.Paths.OutputPath} {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if err := os.MkdirAll(p, 0o755); err != nil {
			panic(fmt.Sprintf("create config path %q: %v", p, err))
		}
	}
}

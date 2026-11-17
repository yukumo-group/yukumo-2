package config

import (
	"strings"

	"yukumo-script/core/internal/domain"
	"yukumo-script/core/internal/store"
	"yukumo-script/core/port"
)

const fileName = "config.yaml"

type Service struct {
	store    port.Store
	defaults domain.Config
}

func New(s port.Store, defaults domain.Config) *Service {
	return &Service{store: s, defaults: defaults}
}

func (s *Service) Get() (domain.Config, error) {
	cfg, err := store.LoadYAMLOrDefault(s.store, fileName, s.defaults)
	if err != nil {
		return cfg, err
	}
	applyDefaults(&cfg, s.defaults)
	return cfg, nil
}

func (s *Service) Save(cfg domain.Config) error {
	applyDefaults(&cfg, s.defaults)
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

func applyDefaults(cfg *domain.Config, def domain.Config) {
	if strings.TrimSpace(cfg.Paths.CharacterPath) == "" {
		cfg.Paths.CharacterPath = def.Paths.CharacterPath
	}
	if strings.TrimSpace(cfg.Paths.OutputPath) == "" {
		cfg.Paths.OutputPath = def.Paths.OutputPath
	}
	if strings.TrimSpace(cfg.Synth.OutputFormat) == "" {
		cfg.Synth.OutputFormat = def.Synth.OutputFormat
	}
	if cfg.Synth.AlwaysWrite == nil {
		cfg.Synth.AlwaysWrite = boolPtr(def.Synth.AlwaysWrite, true)
	}
	if cfg.Synth.ReuseCache == nil {
		cfg.Synth.ReuseCache = boolPtr(def.Synth.ReuseCache, true)
	}
	if cfg.CLI.DefaultInteractive == nil {
		cfg.CLI.DefaultInteractive = boolPtr(def.CLI.DefaultInteractive, true)
	}
	if cfg.CLI.WithAccessibility == nil {
		cfg.CLI.WithAccessibility = boolPtr(def.CLI.WithAccessibility, false)
	}
}

func boolPtr(v *bool, def bool) *bool {
	b := def
	if v != nil {
		b = *v
	}
	return &b
}

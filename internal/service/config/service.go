package config

import (
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
	return store.LoadYAMLOrDefault(s.store, fileName, domain.DefaultConfig())
}

func (s *Service) Save(cfg domain.Config) error {
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

package character

import (
	"context"
	"fmt"
	"strings"

	"yukumo-script/internal/domain"
	svcconfig "yukumo-script/internal/service/config"
	"yukumo-script/internal/store"
)

type PresetCatalog interface {
	HasPreset(ctx context.Context, version, preset int32) (bool, error)
	PresetNames(ctx context.Context, version int32) ([]string, error)
}

type Service struct {
	config  *svcconfig.Service
	store   store.Store // optional override for tests
	presets PresetCatalog
}

func New(cfg *svcconfig.Service, presets PresetCatalog) *Service {
	return &Service{config: cfg, presets: presets}
}

func NewWithStore(s store.Store, presets PresetCatalog) *Service {
	return &Service{store: s, presets: presets}
}

func (s *Service) Presets() PresetCatalog { return s.presets }

func (s *Service) backend() (store.Store, error) {
	if s.store != nil {
		return s.store, nil
	}
	cfg, err := s.config.Get()
	if err != nil {
		return nil, err
	}
	root := strings.TrimSpace(cfg.Paths.CharacterPath)
	if root == "" {
		return nil, fmt.Errorf("%w: paths.character_path is empty", domain.ErrInvalid)
	}
	return store.NewFS(root), nil
}

func (s *Service) initDefaultVoice() bool {
	if s.config == nil {
		return false
	}
	cfg, err := s.config.Get()
	if err != nil {
		return false
	}
	return cfg.Characters.InitDefaultVoice
}

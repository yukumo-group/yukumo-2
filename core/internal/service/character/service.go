package character

import (
	"fmt"
	"strings"

	"yukumo-script/core/internal/domain"
	svcconfig "yukumo-script/core/internal/service/config"
	"yukumo-script/core/internal/store"
)

type Service struct {
	config *svcconfig.Service
	store  store.Store // optional override for tests
}

func New(cfg *svcconfig.Service) *Service {
	return &Service{config: cfg}
}

func NewWithStore(s store.Store) *Service {
	return &Service{store: s}
}

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

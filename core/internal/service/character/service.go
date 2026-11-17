package character

import (
	svcconfig "yukumo-script/core/internal/service/config"
	"yukumo-script/core/port"
)

type Service struct {
	config *svcconfig.Service
	store  port.Store
}

func New(cfg *svcconfig.Service, s port.Store) *Service {
	return &Service{config: cfg, store: s}
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

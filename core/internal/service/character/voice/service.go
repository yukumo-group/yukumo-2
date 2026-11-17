package voice

import (
	"yukumo-script/core/internal/domain"
	"yukumo-script/core/port"
)

type Characters interface {
	Get(name string) (domain.Character, error)
	Update(oldName string, c domain.Character) error
}

type Service struct {
	chars   Characters
	presets port.PresetCatalog
}

func New(chars Characters, presets port.PresetCatalog) *Service {
	return &Service{chars: chars, presets: presets}
}

func (s *Service) Presets() port.PresetCatalog { return s.presets }

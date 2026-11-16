package voice

import (
	"context"

	"yukumo-script/core/internal/domain"
)

type Characters interface {
	Get(name string) (domain.Character, error)
	Update(oldName string, c domain.Character) error
}

type PresetCatalog interface {
	HasPreset(ctx context.Context, version, preset int32) (bool, error)
	PresetNames(ctx context.Context, version int32) ([]string, error)
}

type Service struct {
	chars   Characters
	presets PresetCatalog
}

func New(chars Characters, presets PresetCatalog) *Service {
	return &Service{chars: chars, presets: presets}
}

func (s *Service) Presets() PresetCatalog { return s.presets }

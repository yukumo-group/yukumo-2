package voice

import (
	"context"
	"fmt"

	"yukumo-script/core/internal/domain"
)

func (s *Service) validate(ctx context.Context, v domain.Voice) error {
	if len(v.Params) > 3 {
		return fmt.Errorf("%w: at most 3 params (versions 1/2/10)", domain.ErrInvalid)
	}
	seen := map[int]bool{}
	for _, p := range v.Params {
		switch p.Version {
		case domain.Version1, domain.Version2, domain.Version10:
		default:
			return fmt.Errorf("%w: version must be 1, 2, or 10", domain.ErrInvalid)
		}
		if seen[p.Version] {
			return fmt.Errorf("%w: duplicate version %d", domain.ErrInvalid, p.Version)
		}
		seen[p.Version] = true
		if err := s.validatePreset(ctx, p.Version, p.Preset); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) validatePreset(ctx context.Context, version, preset int) error {
	if s.presets == nil {
		return nil
	}
	ok, err := s.presets.HasPreset(ctx, int32(version), int32(preset))
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%w: preset %d not valid for version %d", domain.ErrInvalid, preset, version)
	}
	return nil
}

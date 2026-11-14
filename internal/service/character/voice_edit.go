package character

import (
	"context"
	"fmt"
	"strings"

	"yukumo-script/internal/domain"
)

func (s *Service) UpdateVoice(ctx context.Context, charName, oldName string, v domain.Voice) error {
	oldName = strings.TrimSpace(oldName)
	v.Name = strings.TrimSpace(v.Name)
	if oldName == "" || v.Name == "" {
		return fmt.Errorf("%w: voice name is required", domain.ErrInvalid)
	}
	if err := s.validateVoice(ctx, v); err != nil {
		return err
	}
	c, err := s.Get(charName)
	if err != nil {
		return err
	}
	idx := voiceIndex(c.Voices, oldName)
	if idx < 0 {
		return fmt.Errorf("%w: voice %q", domain.ErrNotFound, oldName)
	}
	if v.Name != oldName && voiceIndex(c.Voices, v.Name) >= 0 {
		return fmt.Errorf("%w: voice %q", domain.ErrExists, v.Name)
	}
	v.Normalize()
	c.Voices[idx] = v
	return s.Update(charName, c)
}

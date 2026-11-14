package character

import (
	"context"
	"fmt"
	"strings"

	"yukumo-script/internal/domain"
)

func (s *Service) AddVoice(ctx context.Context, charName string, v domain.Voice) error {
	v.Name = strings.TrimSpace(v.Name)
	if v.Name == "" {
		return fmt.Errorf("%w: voice name is required", domain.ErrInvalid)
	}
	if err := s.validateVoice(ctx, v); err != nil {
		return err
	}
	c, err := s.Get(charName)
	if err != nil {
		return err
	}
	if voiceIndex(c.Voices, v.Name) >= 0 {
		return fmt.Errorf("%w: voice %q", domain.ErrExists, v.Name)
	}
	v.Normalize()
	c.Voices = append(c.Voices, v)
	return s.Update(charName, c)
}

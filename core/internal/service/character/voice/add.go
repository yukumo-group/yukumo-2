package voice

import (
	"context"
	"fmt"
	"strings"

	"yukumo-script/core/internal/domain"
)

func (s *Service) Add(ctx context.Context, charName string, v domain.Voice) error {
	v.Name = strings.TrimSpace(v.Name)
	if v.Name == "" {
		return fmt.Errorf("%w: voice name is required", domain.ErrInvalid)
	}
	if err := s.validate(ctx, v); err != nil {
		return err
	}
	c, err := s.chars.Get(charName)
	if err != nil {
		return err
	}
	if indexOf(c.Voices, v.Name) >= 0 {
		return fmt.Errorf("%w: voice %q", domain.ErrExists, v.Name)
	}
	v.Normalize()
	c.Voices = append(c.Voices, v)
	return s.chars.Update(charName, c)
}

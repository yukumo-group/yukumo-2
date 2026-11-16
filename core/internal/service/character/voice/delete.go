package voice

import (
	"fmt"
	"strings"

	"yukumo-script/core/internal/domain"
)

func (s *Service) Delete(charName, voiceName string) error {
	voiceName = strings.TrimSpace(voiceName)
	if voiceName == "" {
		return fmt.Errorf("%w: voice name is required", domain.ErrInvalid)
	}
	c, err := s.chars.Get(charName)
	if err != nil {
		return err
	}
	idx := indexOf(c.Voices, voiceName)
	if idx < 0 {
		return fmt.Errorf("%w: voice %q", domain.ErrNotFound, voiceName)
	}
	c.Voices = append(c.Voices[:idx], c.Voices[idx+1:]...)
	return s.chars.Update(charName, c)
}

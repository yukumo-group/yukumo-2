package character

import (
	"fmt"
	"strings"

	"yukumo-script/internal/domain"
)

func (s *Service) DeleteVoice(charName, voiceName string) error {
	voiceName = strings.TrimSpace(voiceName)
	if voiceName == "" {
		return fmt.Errorf("%w: voice name is required", domain.ErrInvalid)
	}
	c, err := s.Get(charName)
	if err != nil {
		return err
	}
	idx := voiceIndex(c.Voices, voiceName)
	if idx < 0 {
		return fmt.Errorf("%w: voice %q", domain.ErrNotFound, voiceName)
	}
	c.Voices = append(c.Voices[:idx], c.Voices[idx+1:]...)
	return s.Update(charName, c)
}

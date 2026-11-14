package voice

import (
	"fmt"
	"strings"

	"yukumo-script/internal/domain"
)

func (s *Service) List(charName string) ([]domain.Voice, error) {
	c, err := s.chars.Get(charName)
	if err != nil {
		return nil, err
	}
	return c.Voices, nil
}

func (s *Service) Get(charName, voiceName string) (domain.Voice, error) {
	voiceName = strings.TrimSpace(voiceName)
	if voiceName == "" {
		return domain.Voice{}, fmt.Errorf("%w: voice name is required", domain.ErrInvalid)
	}
	c, err := s.chars.Get(charName)
	if err != nil {
		return domain.Voice{}, err
	}
	for _, v := range c.Voices {
		if v.Name == voiceName {
			return v, nil
		}
	}
	return domain.Voice{}, fmt.Errorf("%w: voice %q", domain.ErrNotFound, voiceName)
}

func indexOf(voices []domain.Voice, name string) int {
	for i, v := range voices {
		if v.Name == name {
			return i
		}
	}
	return -1
}

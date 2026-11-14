package character

import (
	"fmt"
	"strings"

	"yukumo-script/internal/domain"
	"yukumo-script/internal/store"
)

func (s *Service) Add(c domain.Character) error {
	c.Name = strings.TrimSpace(c.Name)
	if c.Name == "" {
		return fmt.Errorf("%w: name is required", domain.ErrInvalid)
	}
	file, err := FileName(c.Name)
	if err != nil {
		return err
	}
	b, err := s.backend()
	if err != nil {
		return err
	}
	exists, err := b.Exists(file)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: character %q", domain.ErrExists, c.Name)
	}
	if len(c.Voices) == 0 && s.initDefaultVoice() {
		c.Voices = []domain.Voice{{
			Name:   "default",
			Params: []domain.SynthParam{domain.DefaultSynthParam(domain.Version1)},
		}}
	}
	for i := range c.Voices {
		c.Voices[i].Normalize()
	}
	return store.SaveYAML(b, file, c)
}

package character

import (
	"errors"
	"fmt"
	"strings"

	"yukumo-script/core/internal/domain"
	"yukumo-script/core/internal/store"
)

func (s *Service) Get(name string) (domain.Character, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.Character{}, fmt.Errorf("%w: name is required", domain.ErrInvalid)
	}
	file, err := FileName(name)
	if err != nil {
		return domain.Character{}, err
	}
	b, err := s.backend()
	if err != nil {
		return domain.Character{}, err
	}
	c, err := store.LoadYAML[domain.Character](b, file)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return domain.Character{}, fmt.Errorf("%w: character %q", domain.ErrNotFound, name)
		}
		return domain.Character{}, err
	}
	if c.Name == "" {
		c.Name = name
	}
	return c, nil
}

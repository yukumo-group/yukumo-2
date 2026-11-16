package character

import (
	"errors"
	"fmt"
	"strings"

	"yukumo-script/core/internal/domain"
	"yukumo-script/core/internal/store"
)

func (s *Service) Delete(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("%w: name is required", domain.ErrInvalid)
	}
	file, err := FileName(name)
	if err != nil {
		return err
	}
	b, err := s.backend()
	if err != nil {
		return err
	}
	if err := b.Delete(file); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("%w: character %q", domain.ErrNotFound, name)
		}
		return err
	}
	return nil
}

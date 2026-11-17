package character

import (
	"errors"
	"fmt"
	"strings"

	"yukumo-script/core/internal/domain"
	"yukumo-script/core/internal/store"
)

func (s *Service) Update(oldName string, c domain.Character) error {
	oldName = strings.TrimSpace(oldName)
	c.Name = strings.TrimSpace(c.Name)
	if oldName == "" || c.Name == "" {
		return fmt.Errorf("%w: name is required", domain.ErrInvalid)
	}
	oldFile, err := FileName(oldName)
	if err != nil {
		return err
	}
	newFile, err := FileName(c.Name)
	if err != nil {
		return err
	}
	if _, err := store.LoadYAML[domain.Character](s.store, oldFile); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("%w: character %q", domain.ErrNotFound, oldName)
		}
		return err
	}

	if oldFile != newFile {
		exists, err := s.store.Exists(newFile)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("%w: character %q", domain.ErrExists, c.Name)
		}
	}

	for i := range c.Voices {
		c.Voices[i].Normalize()
	}
	if err := store.SaveYAML(s.store, newFile, c); err != nil {
		return err
	}
	if oldFile != newFile {
		if err := s.store.Delete(oldFile); err != nil {
			return err
		}
	}
	return nil
}

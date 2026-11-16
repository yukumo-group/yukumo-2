package character

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"yukumo-script/core/internal/domain"
	"yukumo-script/core/internal/store"
)

func (s *Service) List() ([]domain.Character, error) {
	b, err := s.backend()
	if err != nil {
		return nil, err
	}
	names, err := b.List("")
	if err != nil {
		return nil, err
	}
	out := make([]domain.Character, 0, len(names))
	for _, name := range names {
		if !strings.EqualFold(filepath.Ext(name), ".yaml") && !strings.EqualFold(filepath.Ext(name), ".yml") {
			continue
		}
		c, err := store.LoadYAML[domain.Character](b, name)
		if err != nil {
			return nil, fmt.Errorf("load %s: %w", name, err)
		}
		if c.Name == "" {
			c.Name = strings.TrimSuffix(name, filepath.Ext(name))
		}
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out, nil
}

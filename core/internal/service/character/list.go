package character

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"yukumo-script/core/internal/domain"
	"yukumo-script/core/internal/store"
)

func (s *Service) List() ([]domain.Character, error) {
	names, err := s.store.List("")
	if err != nil {
		return nil, err
	}
	out := make([]domain.Character, 0, len(names))
	for _, name := range names {
		if !strings.EqualFold(path.Ext(name), ".yaml") && !strings.EqualFold(path.Ext(name), ".yml") {
			continue
		}
		c, err := store.LoadYAML[domain.Character](s.store, name)
		if err != nil {
			return nil, fmt.Errorf("load %s: %w", name, err)
		}
		if c.Name == "" {
			c.Name = strings.TrimSuffix(name, path.Ext(name))
		}
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out, nil
}

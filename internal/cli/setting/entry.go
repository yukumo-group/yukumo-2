package setting

import (
	"fmt"

	"yukumo-script/internal/cli/prompt"
)

type Entry[T any] struct {
	Field Field
	Get   func(T) string
	Set   func(*T, string) error
}

func Keys[T any](entries []Entry[T]) []string {
	keys := make([]string, len(entries))
	for i, e := range entries {
		keys[i] = e.Field.Key
	}
	return keys
}

func Lookup[T any](entries []Entry[T], key string) (Entry[T], error) {
	for _, e := range entries {
		if e.Field.Key == key {
			return e, nil
		}
	}
	return Entry[T]{}, fmt.Errorf("unknown key %q (want %v)", key, Keys(entries))
}

func Metas[T any](entries []Entry[T], v T) map[string]prompt.LeafMeta {
	out := make(map[string]prompt.LeafMeta, len(entries))
	for _, e := range entries {
		out[e.Field.Key] = prompt.LeafMeta{
			Value:       e.Get(v),
			Description: e.Field.Description,
		}
	}
	return out
}

package setting

import (
	"yukumo-script/cli/prompt"
	"yukumo-script/cli/schema"
)

// Metas converts schema entries into prompt leaf metadata.
func Metas[T any](entries []schema.Entry[T], v T) map[string]prompt.LeafMeta {
	out := make(map[string]prompt.LeafMeta, len(entries))
	for _, e := range entries {
		out[e.Field.Key] = prompt.LeafMeta{
			Value:       e.Get(v),
			Description: e.Field.Description,
		}
	}
	return out
}

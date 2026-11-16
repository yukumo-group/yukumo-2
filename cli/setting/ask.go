package setting

import (
	"fmt"
	"io"
	"strconv"

	"yukumo-script/cli/prompt"
	"yukumo-script/cli/schema"
)

func AskKey(keys []string, values map[string]prompt.LeafMeta, key *string) error {
	return prompt.EnsureTreeKey("Field", keys, values, key)
}

func BrowseKeys(
	keys []string,
	values func() map[string]prompt.LeafMeta,
	onLeaf func(key string) error,
	startPrefix string,
	focusKey string,
) error {
	return prompt.BrowseKeys("Settings", keys, values, onLeaf, startPrefix, focusKey)
}

func Ask(f schema.Field, current string, notify io.Writer) (string, error) {
	seg := askSegment(f)
	switch f.Kind {
	case schema.KindString:
		v := current
		if err := prompt.String(seg, &v); err != nil {
			return "", err
		}
		return v, nil
	case schema.KindDir:
		v := current
		if err := prompt.Dir(seg, &v); err != nil {
			return "", err
		}
		return v, nil
	case schema.KindFile:
		v := current
		if err := prompt.File(seg, &v); err != nil {
			return "", err
		}
		return v, nil
	case schema.KindInt:
		v, err := strconv.Atoi(current)
		if err != nil {
			v = f.Min
		}
		if err := prompt.RangedInt(seg, &v, f.Min, f.Max, notify); err != nil {
			return "", err
		}
		return strconv.Itoa(v), nil
	case schema.KindBool:
		v, err := prompt.ParseBool(current)
		if err != nil {
			v = false
		}
		prompt.ToggleBool(&v)
		return strconv.FormatBool(v), nil
	default:
		return "", fmt.Errorf("unsupported field kind %d", f.Kind)
	}
}

func askSegment(f schema.Field) string {
	if prompt.Crumb() != "" {
		return ""
	}
	return f.Title
}

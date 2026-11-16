package setting

import (
	"fmt"
	"io"
	"strconv"

	"yukumo-script/cli/prompt"
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

func Ask(f Field, current string, notify io.Writer) (string, error) {
	seg := askSegment(f)
	switch f.Kind {
	case KindString:
		v := current
		if err := prompt.String(seg, &v); err != nil {
			return "", err
		}
		return v, nil
	case KindDir:
		v := current
		if err := prompt.Dir(seg, &v); err != nil {
			return "", err
		}
		return v, nil
	case KindFile:
		v := current
		if err := prompt.File(seg, &v); err != nil {
			return "", err
		}
		return v, nil
	case KindInt:
		v, err := strconv.Atoi(current)
		if err != nil {
			v = f.Min
		}
		if err := prompt.RangedInt(seg, &v, f.Min, f.Max, notify); err != nil {
			return "", err
		}
		return strconv.Itoa(v), nil
	case KindBool:
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

func askSegment(f Field) string {
	if prompt.Crumb() != "" {
		return ""
	}
	return f.Title
}

func Normalize(f Field, raw string, notify io.Writer) (string, error) {
	switch f.Kind {
	case KindString:
		return raw, nil
	case KindDir, KindFile:
		return prompt.ExpandPath(raw), nil
	case KindInt:
		v, clamped, err := prompt.ParseRangedInt(raw, f.Min, f.Max)
		if err != nil {
			return "", err
		}
		if clamped && notify != nil {
			_, _ = fmt.Fprintf(notify, "clamped %s → %d (range [%d, %d])\n", raw, v, f.Min, f.Max)
		}
		return strconv.Itoa(v), nil
	case KindBool:
		v, err := prompt.ParseBool(raw)
		if err != nil {
			return "", err
		}
		return strconv.FormatBool(v), nil
	default:
		return "", fmt.Errorf("unsupported field kind %d", f.Kind)
	}
}

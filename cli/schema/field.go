// Package schema defines the pure field metadata shared by shell and
// action code. It has no cobra dependency.
package schema

import (
	"fmt"
	"strconv"

	"yukumo-script/cli/parse"
)

// Kind describes how a field is parsed and edited.
type Kind uint8

const (
	KindString Kind = iota
	KindDir
	KindFile
	KindInt
	KindBool
)

// Field describes one editable setting.
type Field struct {
	Key         string
	Title       string
	Description string
	Kind        Kind
	Min         int
	Max         int
}

// Entry binds a Field to getter/setter functions for a concrete value.
type Entry[T any] struct {
	Field Field
	Get   func(T) string
	Set   func(*T, string) error
}

// Keys returns the field keys in order.
func Keys[T any](entries []Entry[T]) []string {
	keys := make([]string, len(entries))
	for i, e := range entries {
		keys[i] = e.Field.Key
	}
	return keys
}

// Lookup finds an entry by key.
func Lookup[T any](entries []Entry[T], key string) (Entry[T], error) {
	for _, e := range entries {
		if e.Field.Key == key {
			return e, nil
		}
	}
	return Entry[T]{}, fmt.Errorf("unknown key %q (want %v)", key, Keys(entries))
}

// Normalize parses raw according to f. It returns a warning when an integer
// value was clamped; callers decide where to write that warning.
func Normalize(f Field, raw string) (value, warning string, err error) {
	switch f.Kind {
	case KindString:
		return raw, "", nil
	case KindDir, KindFile:
		return parse.ExpandPath(raw), "", nil
	case KindInt:
		v, clamped, err := parse.ParseRangedInt(raw, f.Min, f.Max)
		if err != nil {
			return "", "", err
		}
		if clamped {
			warning = fmt.Sprintf("clamped %s → %d (range [%d, %d])\n", raw, v, f.Min, f.Max)
		}
		return strconv.Itoa(v), warning, nil
	case KindBool:
		v, err := parse.ParseBool(raw)
		if err != nil {
			return "", "", err
		}
		return strconv.FormatBool(v), "", nil
	default:
		return "", "", fmt.Errorf("unsupported field kind %d", f.Kind)
	}
}

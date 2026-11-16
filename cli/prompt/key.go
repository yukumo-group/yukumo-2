package prompt

import (
	"fmt"
	"slices"

	"github.com/charmbracelet/huh"
)

func Key(title string, keys []string, value *string) error {
	if value == nil {
		return fmt.Errorf("key value is nil")
	}
	return run(huh.NewSelect[string]().
		Title(Title(title)).
		Options(huh.NewOptions(keys...)...).
		Value(value))
}

func EnsureKey(title string, keys []string, value *string) error {
	if value == nil {
		return fmt.Errorf("key value is nil")
	}
	if *value != "" {
		if !slices.Contains(keys, *value) {
			return fmt.Errorf("unknown key %q", *value)
		}
		return nil
	}
	return Key(title, keys, value)
}

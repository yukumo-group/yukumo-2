package prompt

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
)

func ParseBool(raw string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "true", "1", "yes", "y", "on":
		return true, nil
	case "false", "0", "no", "n", "off":
		return false, nil
	default:
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return false, fmt.Errorf("invalid bool %q", raw)
		}
		return v, nil
	}
}

func Bool(title string, value *bool) error {
	if value == nil {
		return fmt.Errorf("bool value is nil")
	}
	return run(huh.NewSelect[bool]().
		Title(Title(title)).
		Options(
			huh.NewOption("true", true),
			huh.NewOption("false", false),
		).
		Value(value))
}

// ToggleBool flips *value in place (settings-style select-to-invert).
func ToggleBool(value *bool) {
	if value != nil {
		*value = !*value
	}
}

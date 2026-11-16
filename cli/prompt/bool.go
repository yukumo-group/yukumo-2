package prompt

import (
	"fmt"

	"yukumo-script/cli/parse"

	"github.com/charmbracelet/huh"
)

func ParseBool(raw string) (bool, error) {
	return parse.ParseBool(raw)
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

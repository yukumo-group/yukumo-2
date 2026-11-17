package prompt

import "github.com/charmbracelet/huh"

func String(segment string, value *string) error {
	return run(huh.NewInput().
		Title(Title(segment)).
		Value(value))
}

func EnsureString(segment string, value *string) error {
	if value == nil || *value != "" {
		return nil
	}
	return String(segment, value)
}

func Confirm(segment string, value *bool) error {
	return run(huh.NewConfirm().
		Title(Title(segment)).
		Value(value))
}

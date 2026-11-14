package prompt

import "github.com/charmbracelet/huh"

func String(title string, value *string) error {
	return run(huh.NewInput().
		Title(title).
		Value(value))
}

func EnsureString(title string, value *string) error {
	if value == nil || *value != "" {
		return nil
	}
	return String(title, value)
}

func Confirm(title string, value *bool) error {
	return run(huh.NewConfirm().
		Title(title).
		Value(value))
}

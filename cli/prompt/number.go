package prompt

import (
	"fmt"
	"io"
	"strconv"

	"github.com/charmbracelet/huh"
)

func ClampInt(v, min, max int) (int, bool) {
	switch {
	case v < min:
		return min, true
	case v > max:
		return max, true
	default:
		return v, false
	}
}

func ParseRangedInt(raw string, min, max int) (int, bool, error) {
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, false, fmt.Errorf("invalid number %q", raw)
	}
	clamped, changed := ClampInt(v, min, max)
	return clamped, changed, nil
}

func RangedInt(title string, value *int, min, max int, notify io.Writer) error {
	if value == nil {
		return fmt.Errorf("number value is nil")
	}
	raw := strconv.Itoa(*value)
	if err := run(huh.NewInput().
		Title(Title(title)).
		Description(fmt.Sprintf("[%d, %d]", min, max)).
		Value(&raw).
		Validate(func(s string) error {
			_, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("invalid number")
			}
			return nil
		})); err != nil {
		return err
	}
	v, clamped, err := ParseRangedInt(raw, min, max)
	if err != nil {
		return err
	}
	if clamped && notify != nil {
		fmt.Fprintf(notify, "clamped %s → %d (range [%d, %d])\n", raw, v, min, max)
	}
	*value = v
	return nil
}

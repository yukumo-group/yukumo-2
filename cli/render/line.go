package render

import (
	"fmt"
	"io"
)

// Line writes a formatted line.
func Line(w io.Writer, format string, args ...any) error {
	_, err := fmt.Fprintf(w, format, args...)
	return err
}

// Warning writes a non-fatal warning line.
func Warning(w io.Writer, s string) error {
	_, err := io.WriteString(w, s)
	return err
}

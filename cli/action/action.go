// Package action implements CLI use cases. It returns domain values and does
// not write output, play audio, or depend on cobra/huh.
package action

// EditResult is the user-facing result of a field edit. Warning is a
// non-fatal normalization message that the caller may write to stderr.
type EditResult struct {
	Value   string
	Warning string
}

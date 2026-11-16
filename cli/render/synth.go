package render

import (
	"fmt"
	"io"
	"time"

	"yukumo-script/core/api"
)

// SynthResultText formats the one-shot synth result line.
func SynthResultText(res api.SynthResult, elapsed time.Duration) string {
	ms := ElapsedText(elapsed)
	switch {
	case res.Wrote && res.FromCache:
		return fmt.Sprintf("%s (cached) %s", res.Path, ms)
	case res.Wrote:
		return fmt.Sprintf("%s %s", res.Path, ms)
	case res.FromCache:
		return fmt.Sprintf("(cached) %s", ms)
	default:
		return ms
	}
}

// SynthStatusText formats the interactive synth status label.
func SynthStatusText(res api.SynthResult, elapsed time.Duration) string {
	ms := ElapsedText(elapsed)
	switch {
	case res.Wrote && res.FromCache:
		return res.Path + " (cached) " + ms
	case res.Wrote:
		return res.Path + " " + ms
	case res.FromCache:
		return "ok (cached, not written) " + ms
	default:
		return "ok (not written) " + ms
	}
}

// SynthResult writes the one-shot synth result line when there is one.
func SynthResult(w io.Writer, res api.SynthResult, elapsed time.Duration) error {
	line := SynthResultText(res, elapsed)
	if line == "" {
		return nil
	}
	_, err := fmt.Fprintln(w, line)
	return err
}

// ElapsedText formats a duration for synth output.
func ElapsedText(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d)/float64(time.Millisecond))
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

package synth

import (
	"context"
	"fmt"
	"os"
	"time"

	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

func cmdContext(cmd *cobra.Command) context.Context {
	if cmd != nil {
		if ctx := cmd.Context(); ctx != nil {
			return ctx
		}
	}
	return context.Background()
}

// runRequest synthesizes, prints the output path (with timing) when written, and optionally previews.
// When asyncPreview is true, playback starts in the background so the UI can continue.
func runRequest(cmd *cobra.Command, app *api.App, req api.SynthRequest, asyncPreview bool) (api.SynthResult, time.Duration, error) {
	start := time.Now()
	res, err := app.Synthesize(cmdContext(cmd), req)
	elapsed := time.Since(start)
	if err != nil {
		return res, elapsed, err
	}
	line := formatResultLine(res, elapsed)
	if line != "" {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), line)
	}
	if req.Preview {
		if err := previewResult(res, asyncPreview); err != nil {
			return res, elapsed, err
		}
	}
	return res, elapsed, nil
}

func formatResultLine(res api.SynthResult, elapsed time.Duration) string {
	ms := formatElapsed(elapsed)
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

func formatElapsed(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d)/float64(time.Millisecond))
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

func previewResult(res api.SynthResult, async bool) error {
	path := res.Path
	cleanup := ""
	if path == "" {
		f, err := os.CreateTemp("", "yukumo-synth-*.wav")
		if err != nil {
			return err
		}
		cleanup = f.Name()
		if _, err := f.Write(res.WAV); err != nil {
			_ = f.Close()
			_ = os.Remove(cleanup)
			return err
		}
		if err := f.Close(); err != nil {
			_ = os.Remove(cleanup)
			return err
		}
		path = cleanup
	}
	if async {
		playWAVAsync(path, cleanup)
		return nil
	}
	defer func() {
		if cleanup != "" {
			_ = os.Remove(cleanup)
		}
	}()
	return playWAV(path)
}

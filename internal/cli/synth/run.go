package synth

import (
	"context"
	"fmt"
	"os"

	svcsynth "yukumo-script/internal/service/synth"

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

// runRequest synthesizes, prints the output path when written, and optionally previews.
// When asyncPreview is true, playback starts in the background so the UI can continue.
func runRequest(cmd *cobra.Command, svc *svcsynth.Service, req svcsynth.Request, asyncPreview bool) (svcsynth.Result, error) {
	res, err := svc.Synthesize(cmdContext(cmd), req)
	if err != nil {
		return res, err
	}
	if res.Wrote {
		if res.FromCache {
			fmt.Fprintf(cmd.OutOrStdout(), "%s (cached)\n", res.Path)
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), res.Path)
		}
	} else if res.FromCache {
		fmt.Fprintln(cmd.OutOrStdout(), "(cached)")
	}
	if req.Preview {
		if err := previewResult(res, asyncPreview); err != nil {
			return res, err
		}
	}
	return res, nil
}

func previewResult(res svcsynth.Result, async bool) error {
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

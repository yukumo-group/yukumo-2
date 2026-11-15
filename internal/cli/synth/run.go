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
func runRequest(cmd *cobra.Command, svc *svcsynth.Service, req svcsynth.Request) (svcsynth.Result, error) {
	res, err := svc.Synthesize(cmdContext(cmd), req)
	if err != nil {
		return res, err
	}
	if res.Wrote {
		fmt.Fprintln(cmd.OutOrStdout(), res.Path)
	}
	if req.Preview {
		if err := previewResult(res); err != nil {
			return res, err
		}
	}
	return res, nil
}

func previewResult(res svcsynth.Result) error {
	path := res.Path
	if path != "" {
		return playWAV(path)
	}
	f, err := os.CreateTemp("", "yukumo-synth-*.wav")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	if _, err := f.Write(res.WAV); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return playWAV(tmp)
}

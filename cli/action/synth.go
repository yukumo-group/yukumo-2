package action

import (
	"context"
	"time"

	"yukumo-script/core/api"
)

// Synthesize runs the synth use case and reports elapsed time. It does not
// print or play the result.
func Synthesize(ctx context.Context, app *api.App, req api.SynthRequest) (api.SynthResult, time.Duration, error) {
	start := time.Now()
	res, err := app.Synthesize(ctx, req)
	return res, time.Since(start), err
}

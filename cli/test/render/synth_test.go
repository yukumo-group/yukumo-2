package render_test

import (
	"testing"
	"time"

	"yukumo-script/cli/render"
	"yukumo-script/core/api"
)

func TestFormatResultLine(t *testing.T) {
	elapsed := 12*time.Millisecond + 340*time.Microsecond
	tests := []struct {
		name string
		res  api.SynthResult
		want string
	}{
		{
			name: "wrote",
			res:  api.SynthResult{Wrote: true, Path: "/tmp/a.wav"},
			want: "/tmp/a.wav 12.34ms",
		},
		{
			name: "wrote cached",
			res:  api.SynthResult{Wrote: true, FromCache: true, Path: "/tmp/a.wav"},
			want: "/tmp/a.wav (cached) 12.34ms",
		},
		{
			name: "cached not written",
			res:  api.SynthResult{FromCache: true},
			want: "(cached) 12.34ms",
		},
		{
			name: "not written",
			res:  api.SynthResult{},
			want: "12.34ms",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := render.SynthResultText(tt.res, elapsed); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatElapsed(t *testing.T) {
	if got := render.ElapsedText(1500 * time.Millisecond); got != "1.50s" {
		t.Fatalf("got %q, want 1.50s", got)
	}
}

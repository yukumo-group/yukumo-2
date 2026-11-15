package synth

import (
	"testing"
	"time"

	svcsynth "yukumo-script/internal/service/synth"
)

func TestFormatResultLine(t *testing.T) {
	elapsed := 12*time.Millisecond + 340*time.Microsecond
	tests := []struct {
		name string
		res  svcsynth.Result
		want string
	}{
		{
			name: "wrote",
			res:  svcsynth.Result{Wrote: true, Path: "/tmp/a.wav"},
			want: "/tmp/a.wav 12.34ms",
		},
		{
			name: "wrote cached",
			res:  svcsynth.Result{Wrote: true, FromCache: true, Path: "/tmp/a.wav"},
			want: "/tmp/a.wav (cached) 12.34ms",
		},
		{
			name: "cached not written",
			res:  svcsynth.Result{FromCache: true},
			want: "(cached) 12.34ms",
		},
		{
			name: "not written",
			res:  svcsynth.Result{},
			want: "12.34ms",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatResultLine(tt.res, elapsed); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatElapsed(t *testing.T) {
	if got := formatElapsed(1500 * time.Millisecond); got != "1.50s" {
		t.Fatalf("got %q, want 1.50s", got)
	}
}

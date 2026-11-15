package synth_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"yukumo-script/internal/service/synth"
)

func TestFormatOutputName(t *testing.T) {
	got := synth.FormatOutputName("c-v-ymd", "remilia", "default")
	wantPrefix := "remilia-default-" + time.Now().Format("20060102")
	if got != wantPrefix {
		t.Fatalf("got %q want %q", got, wantPrefix)
	}
	got = synth.FormatOutputName("{c}_{v}_{ymd}", "A B", "v1")
	if !strings.HasPrefix(got, "A_B_v1_") {
		t.Fatalf("braced/sanitize: %q", got)
	}
}

func TestUniquePath(t *testing.T) {
	dir := t.TempDir()
	p1, err := synth.UniquePath(dir, "out", ".wav")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p1, []byte("a"), 0o644); err != nil {
		t.Fatal(err)
	}
	p2, err := synth.UniquePath(dir, "out", ".wav")
	if err != nil {
		t.Fatal(err)
	}
	if p2 != filepath.Join(dir, "out-1.wav") {
		t.Fatalf("got %q", p2)
	}
}

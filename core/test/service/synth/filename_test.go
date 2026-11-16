package synth_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"yukumo-script/core/internal/domain"
	"yukumo-script/core/internal/service/synth"
)

func TestFormatOutputNameDefault(t *testing.T) {
	got := synth.FormatOutputName("", "remilia", "default")
	now := time.Now()
	pat := fmt.Sprintf(`^remilia:default_%04d-%02d-%02d:\d{2}-\d{2}-\d{2}_\d{3}$`,
		now.Year(), int(now.Month()), now.Day())
	if !regexp.MustCompile(pat).MatchString(got) {
		t.Fatalf("got %q want ~%s", got, pat)
	}
	if domain.DefaultOutputFormat != "{character}:{voice}_{yyyy}-{MM}-{dd}:{hh}-{mm}-{ss}_{ms}" {
		t.Fatalf("default format: %q", domain.DefaultOutputFormat)
	}
	got2 := synth.FormatOutputName(domain.DefaultOutputFormat, "A B", "v1")
	if !regexp.MustCompile(`^A_B:v1_\d{4}-\d{2}-\d{2}:\d{2}-\d{2}-\d{2}_\d{3}$`).MatchString(got2) {
		t.Fatalf("sanitize/braced: %q", got2)
	}
}

func TestFormatOutputNameWidths(t *testing.T) {
	got := synth.FormatOutputName("{yyyy}-{MM}-{dd} {hh}:{mm}:{ss}", "c", "v")
	if !regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`).MatchString(got) {
		t.Fatalf("got %q", got)
	}
	y2 := synth.FormatOutputName("{yy}", "c", "v")
	if !regexp.MustCompile(`^\d{2}$`).MatchString(y2) {
		t.Fatalf("yy: %q", y2)
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

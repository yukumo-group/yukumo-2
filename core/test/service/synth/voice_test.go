package synth_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"yukumo-script/core/internal/domain"
	"yukumo-script/core/internal/service/synth"
)

func TestParseVoiceRef(t *testing.T) {
	c, v, err := synth.ParseVoiceRef("remilia:default")
	if err != nil || c != "remilia" || v != "default" {
		t.Fatalf("got %q %q %v", c, v, err)
	}
	if _, _, err := synth.ParseVoiceRef("remilia.default"); !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("dot ref should be invalid: %v", err)
	}
	if _, _, err := synth.ParseVoiceRef("nocolon"); !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("expected invalid: %v", err)
	}
}

func TestParseVoicePathSpec(t *testing.T) {
	path, voice, err := synth.ParseVoicePathSpec(`/tmp/foo:bar.yaml:default`)
	if err != nil || path != `/tmp/foo:bar.yaml` || voice != "default" {
		t.Fatalf("got %q %q %v", path, voice, err)
	}
}

func TestLoadVoiceFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "remilia.yaml")
	content := "name: remilia\nvoices:\n  - name: default\n    params:\n      - version: 1\n        preset: 0\n        speed: 100\n"
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	c, v, raw, err := synth.LoadVoiceFile(p, "default")
	if err != nil {
		t.Fatal(err)
	}
	if c.Name != "remilia" || v.Name != "default" || len(raw) == 0 {
		t.Fatalf("got %+v %+v len=%d", c, v, len(raw))
	}
}

func TestSelectParam(t *testing.T) {
	v := domain.Voice{
		Name: "default",
		Params: []domain.SynthParam{
			{Version: domain.Version2, Preset: 1, Speed: 100},
			{Version: domain.Version10, Preset: 0, Speed: 100},
		},
	}
	p, err := synth.SelectParam(v, 0)
	if err != nil || p.Version != domain.Version2 {
		t.Fatalf("auto prefer 2 over 10 when no v1: %+v %v", p, err)
	}
	p, err = synth.SelectParam(v, domain.Version10)
	if err != nil || p.Version != domain.Version10 {
		t.Fatalf("explicit: %+v %v", p, err)
	}
	_, err = synth.SelectParam(domain.Voice{Name: "x"}, 0)
	if !errors.Is(err, domain.ErrInvalid) || !strings.Contains(err.Error(), "params empty") {
		t.Fatalf("empty params: %v", err)
	}
	_, err = synth.SelectParam(v, domain.Version1)
	if !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("missing version: %v", err)
	}
}

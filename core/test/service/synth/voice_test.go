package synth_test

import (
	"errors"
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

func TestLoadVoiceBytes(t *testing.T) {
	content := []byte("name: remilia\nvoices:\n  - name: default\n    params:\n      - version: 1\n        preset: 0\n        speed: 100\n")
	c, v, raw, err := synth.LoadVoiceBytes(content, "default", "fallback")
	if err != nil {
		t.Fatal(err)
	}
	if c.Name != "remilia" || v.Name != "default" || len(raw) == 0 {
		t.Fatalf("got %+v %+v len=%d", c, v, len(raw))
	}
	// no name in YAML falls back to the hint
	content = []byte("voices:\n  - name: default\n    params:\n      - version: 1\n")
	c, _, _, err = synth.LoadVoiceBytes(content, "default", "fallback")
	if err != nil || c.Name != "fallback" {
		t.Fatalf("fallback = %+v %v", c, err)
	}
}

func TestSelectParam(t *testing.T) {
	v := domain.Voice{
		Name:  "default",
		Speed: 100,
		Params: []domain.SynthParam{
			{Version: domain.Version2, Preset: 1},
			{Version: domain.Version10, Preset: 0},
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

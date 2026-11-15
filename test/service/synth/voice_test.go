package synth_test

import (
	"errors"
	"strings"
	"testing"

	"yukumo-script/internal/domain"
	"yukumo-script/internal/service/synth"
)

func TestParseVoiceRef(t *testing.T) {
	c, v, err := synth.ParseVoiceRef("remilia.default")
	if err != nil || c != "remilia" || v != "default" {
		t.Fatalf("got %q %q %v", c, v, err)
	}
	if _, _, err := synth.ParseVoiceRef("nodot"); !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("expected invalid: %v", err)
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

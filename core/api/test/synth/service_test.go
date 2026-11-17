package synth_test

import (
	"context"
	"testing"

	"yukumo-script/core/internal/domain"
	svccharacter "yukumo-script/core/internal/service/character"
	svcvoice "yukumo-script/core/internal/service/character/voice"
	svcconfig "yukumo-script/core/internal/service/config"
	"yukumo-script/core/internal/service/synth"
	"yukumo-script/core/internal/store"
	"yukumo-script/runtime/aq"
)

type allowAllPresets struct{}

func (allowAllPresets) HasPreset(context.Context, int32, int32) (bool, error) {
	return true, nil
}
func (allowAllPresets) PresetNames(context.Context, int32) ([]string, error) {
	return []string{"p0"}, nil
}

func newSynthService(t *testing.T, alwaysWrite bool) (*synth.Service, *store.Memory) {
	t.Helper()
	ctx := context.Background()
	cfgMem := store.NewMemory()
	charsMem := store.NewMemory()
	artifacts := store.NewMemory()
	cfgSvc := svcconfig.New(cfgMem, domain.DefaultConfig(domain.PathConfig{}))
	if _, err := cfgSvc.Edit(func(c *domain.Config) {
		c.Synth.OutputFormat = "{character}-{voice}"
		c.Synth.SetAlwaysWrite(alwaysWrite)
	}); err != nil {
		t.Fatal(err)
	}
	chars := svccharacter.New(cfgSvc, charsMem)
	voices := svcvoice.New(chars, allowAllPresets{})
	if err := chars.Add(domain.Character{Name: "remilia"}); err != nil {
		t.Fatal(err)
	}
	if err := voices.Add(ctx, "remilia", domain.Voice{
		Name:   "default",
		Params: []domain.SynthParam{{Version: domain.Version1, Preset: 0, Speed: 100}},
	}); err != nil {
		t.Fatal(err)
	}
	pool := aq.NewPool()
	t.Cleanup(func() { _ = pool.Close(ctx) })
	return synth.New(cfgSvc, voices, aqEngine{pool: pool}, artifacts, nil), artifacts
}

func TestSynthesizeVoiceWritesArtifact(t *testing.T) {
	ctx := context.Background()
	svc, artifacts := newSynthService(t, true)
	res, err := svc.Synthesize(ctx, synth.Request{
		VoiceRef: "remilia:default",
		Text:     "こんにちは",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Wrote || res.Path == "" {
		t.Fatalf("expected write: %+v", res)
	}
	ok, err := artifacts.Exists(res.Path)
	if err != nil || !ok {
		t.Fatalf("artifact %q missing: %v", res.Path, err)
	}
	if len(res.WAV) < 44 {
		t.Fatalf("short wav: %d", len(res.WAV))
	}
}

func TestSynthesizePreviewSkipWrite(t *testing.T) {
	ctx := context.Background()
	svc, artifacts := newSynthService(t, false)
	res, err := svc.Synthesize(ctx, synth.Request{
		VoiceRef: "remilia:default",
		Text:     "あ",
		Preview:  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Wrote {
		t.Fatalf("preview should skip write: %+v", res)
	}
	keys, err := artifacts.List("")
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 0 {
		t.Fatalf("artifact store should be empty, got %v", keys)
	}
}

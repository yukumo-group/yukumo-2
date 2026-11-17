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

func TestTaskHashStableAndDistinct(t *testing.T) {
	p := domain.SynthParam{Version: 1, Preset: 0, Speed: 100}
	a := synth.TaskHash("あ", 100, p, nil)
	b := synth.TaskHash("あ", 100, p, nil)
	if a != b || a == "" {
		t.Fatalf("unstable hash %q %q", a, b)
	}
	if synth.TaskHash("い", 100, p, nil) == a {
		t.Fatal("text should change hash")
	}
	if synth.TaskHash("あ", 100, p, []byte("file")) == a {
		t.Fatal("file bytes should change hash")
	}
}

func TestReuseCacheWritesNewArtifact(t *testing.T) {
	ctx := context.Background()
	cfgMem := store.NewMemory()
	charsMem := store.NewMemory()
	artifacts := store.NewMemory()
	cfgSvc := svcconfig.New(cfgMem, domain.DefaultConfig(domain.PathConfig{}))
	if _, err := cfgSvc.Edit(func(c *domain.Config) {
		c.Synth.OutputFormat = "{character}-{voice}"
		c.Synth.SetAlwaysWrite(true)
		c.Synth.SetReuseCache(true)
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
	svc := synth.New(cfgSvc, voices, aqEngine{pool: pool}, artifacts, nil)

	r1, err := svc.Synthesize(ctx, synth.Request{VoiceRef: "remilia:default", Text: "こんにちは"})
	if err != nil {
		t.Fatal(err)
	}
	if !r1.Wrote || r1.FromCache {
		t.Fatalf("first: %+v", r1)
	}
	r2, err := svc.Synthesize(ctx, synth.Request{VoiceRef: "remilia:default", Text: "こんにちは"})
	if err != nil {
		t.Fatal(err)
	}
	if !r2.FromCache || !r2.Wrote || r2.Path == r1.Path {
		t.Fatalf("second should copy to new key: first=%q second=%+v", r1.Path, r2)
	}
	for _, key := range []string{r1.Path, r2.Path} {
		ok, err := artifacts.Exists(key)
		if err != nil || !ok {
			t.Fatalf("artifact %q missing: %v", key, err)
		}
	}
}

func TestCleanIndexDropsMissing(t *testing.T) {
	mem := store.NewMemory()
	if err := mem.Write("index.yaml", []byte("deadbeef: missing.wav\n")); err != nil {
		t.Fatal(err)
	}
	svc := synth.New(nil, nil, nil, mem, nil)
	if err := svc.CleanIndex(); err != nil {
		t.Fatal(err)
	}
	data, err := mem.Read("index.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "{}\n" {
		t.Fatalf("expected cleaned index, got %q", data)
	}
}

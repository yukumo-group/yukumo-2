package synth_test

import (
	"context"
	"os"
	"path/filepath"
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

func TestReuseCacheCopiesNewFile(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	charsDir := filepath.Join(dir, "characters")
	outDir := filepath.Join(dir, "out")
	_ = os.MkdirAll(charsDir, 0o755)

	mem := store.NewMemory()
	cfgSvc := svcconfig.New(mem)
	if _, err := cfgSvc.Edit(func(c *domain.Config) {
		c.Paths.CharacterPath = charsDir
		c.Paths.OutputPath = outDir
		c.Synth.OutputFormat = "{character}-{voice}"
		c.Synth.SetAlwaysWrite(true)
		c.Synth.SetReuseCache(true)
	}); err != nil {
		t.Fatal(err)
	}
	chars := svccharacter.New(cfgSvc)
	voices := svcvoice.New(chars, allowAllPresets{})
	_ = chars.Add(domain.Character{Name: "remilia"})
	_ = voices.Add(ctx, "remilia", domain.Voice{
		Name:   "default",
		Params: []domain.SynthParam{{Version: domain.Version1, Preset: 0, Speed: 100}},
	})
	pool := aq.NewPool()
	t.Cleanup(func() { _ = pool.Close(ctx) })
	svc := synth.New(cfgSvc, voices, aqEngine{pool: pool})

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
		t.Fatalf("second should copy to new path: first=%q second=%+v", r1.Path, r2)
	}
	if _, err := os.Stat(r1.Path); err != nil {
		t.Fatal("original cache source should remain", err)
	}
	if _, err := os.Stat(r2.Path); err != nil {
		t.Fatal(err)
	}
}

func TestCleanIndexDropsMissing(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out")
	_ = os.MkdirAll(out, 0o755)
	idxPath := filepath.Join(out, "index.yaml")
	if err := os.WriteFile(idxPath, []byte("deadbeef: missing.wav\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := synth.CleanIndex(out); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(idxPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) > 5 && string(data) != "{}\n" && string(data) != "{}\n\n" {
		// empty map marshals as "{}\n"
		if string(data) != "{}\n" {
			t.Fatalf("expected cleaned index, got %q", data)
		}
	}
}

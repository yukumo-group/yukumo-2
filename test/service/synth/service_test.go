package synth_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"yukumo-script/internal/domain"
	svccharacter "yukumo-script/internal/service/character"
	svcvoice "yukumo-script/internal/service/character/voice"
	svcconfig "yukumo-script/internal/service/config"
	"yukumo-script/internal/service/synth"
	"yukumo-script/internal/runtime/aq"
	"yukumo-script/internal/store"
)

type allowAllPresets struct{}

func (allowAllPresets) HasPreset(context.Context, int32, int32) (bool, error) {
	return true, nil
}
func (allowAllPresets) PresetNames(context.Context, int32) ([]string, error) {
	return []string{"p0"}, nil
}

func TestSynthesizeVoiceWritesFile(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	charsDir := filepath.Join(dir, "characters")
	outDir := filepath.Join(dir, "out")
	if err := os.MkdirAll(charsDir, 0o755); err != nil {
		t.Fatal(err)
	}

	mem := store.NewMemory()
	cfgSvc := svcconfig.New(mem)
	if _, err := cfgSvc.Edit(func(c *domain.Config) {
		c.Paths.CharacterPath = charsDir
		c.Paths.OutputPath = outDir
		c.Synth.OutputFormat = "{character}-{voice}"
		c.Synth.SetAlwaysWrite(true)
	}); err != nil {
		t.Fatal(err)
	}

	chars := svccharacter.New(cfgSvc)
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
	svc := synth.New(cfgSvc, voices, pool)

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
	if _, err := os.Stat(res.Path); err != nil {
		t.Fatal(err)
	}
	if len(res.WAV) < 44 {
		t.Fatalf("short wav: %d", len(res.WAV))
	}
}

func TestSynthesizePreviewSkipWrite(t *testing.T) {
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
		c.Synth.SetAlwaysWrite(false)
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
	svc := synth.New(cfgSvc, voices, pool)

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
	entries, _ := os.ReadDir(outDir)
	if len(entries) != 0 {
		t.Fatalf("output dir should be empty, got %v", entries)
	}
}

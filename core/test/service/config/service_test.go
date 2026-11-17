package config_test

import (
	"path/filepath"
	"testing"

	"yukumo-script/core/internal/domain"
	svcconfig "yukumo-script/core/internal/service/config"
	"yukumo-script/core/internal/store"
)

func TestGetDefaultAndEdit(t *testing.T) {
	dir := t.TempDir()
	chars := filepath.Join(dir, "characters")
	out := filepath.Join(dir, "output")
	defaults := domain.DefaultConfig(domain.PathConfig{CharacterPath: chars, OutputPath: out})
	mem := store.NewMemory()
	svc := svcconfig.New(mem, defaults)

	cfg, err := svc.Get()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Paths.CharacterPath != chars ||
		cfg.Paths.OutputPath != out ||
		cfg.Characters.InitDefaultVoice {
		t.Fatalf("unexpected default: %+v", cfg)
	}
	if !cfg.CLI.Interactive() || cfg.CLI.Accessibility() {
		t.Fatalf("unexpected cli defaults: interactive=%v accessibility=%v",
			cfg.CLI.Interactive(), cfg.CLI.Accessibility())
	}
	if cfg.Synth.OutputFormat != domain.DefaultOutputFormat || !cfg.Synth.WriteAlways() || !cfg.Synth.CacheReuse() {
		t.Fatalf("unexpected synth defaults: %+v", cfg.Synth)
	}

	cfg, err = svc.Edit(func(c *domain.Config) {
		c.Paths.CharacterPath = chars
		c.Paths.OutputPath = out
		c.Characters.InitDefaultVoice = true
		c.CLI.SetInteractive(false)
		c.CLI.SetAccessibility(true)
		c.Synth.OutputFormat = "{character}-{voice}"
		c.Synth.SetAlwaysWrite(false)
	})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.CLI.Interactive() || !cfg.CLI.Accessibility() {
		t.Fatalf("cli after edit: %+v", cfg.CLI)
	}
	if cfg.Synth.OutputFormat != "{character}-{voice}" || cfg.Synth.WriteAlways() {
		t.Fatalf("synth after edit: %+v", cfg.Synth)
	}
	got, err := svc.Get()
	if err != nil {
		t.Fatal(err)
	}
	if got.Paths.CharacterPath != chars || got.Paths.OutputPath != out || !got.Characters.InitDefaultVoice {
		t.Fatalf("persisted = %+v", got)
	}
	if got.CLI.Interactive() || !got.CLI.Accessibility() {
		t.Fatalf("persisted cli = interactive=%v accessibility=%v",
			got.CLI.Interactive(), got.CLI.Accessibility())
	}
	if got.Synth.OutputFormat != "{character}-{voice}" || got.Synth.WriteAlways() {
		t.Fatalf("persisted synth = %+v", got.Synth)
	}
}

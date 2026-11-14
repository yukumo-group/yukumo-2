package config_test

import (
	"testing"

	"yukumo-script/internal/domain"
	svcconfig "yukumo-script/internal/service/config"
	"yukumo-script/internal/store"
)

func TestGetDefaultAndEdit(t *testing.T) {
	mem := store.NewMemory()
	svc := svcconfig.New(mem)

	cfg, err := svc.Get()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Paths.CharacterPath != "" || cfg.Paths.DataPath != "" || cfg.Characters.InitDefaultVoice {
		t.Fatalf("unexpected default: %+v", cfg)
	}
	if !cfg.CLI.Interactive() || cfg.CLI.Accessibility() {
		t.Fatalf("unexpected cli defaults: interactive=%v accessibility=%v",
			cfg.CLI.Interactive(), cfg.CLI.Accessibility())
	}

	cfg, err = svc.Edit(func(c *domain.Config) {
		c.Paths.CharacterPath = "/chars"
		c.Paths.DataPath = "/data"
		c.Characters.InitDefaultVoice = true
		c.CLI.SetInteractive(false)
		c.CLI.SetAccessibility(true)
	})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.CLI.Interactive() || !cfg.CLI.Accessibility() {
		t.Fatalf("cli after edit: %+v", cfg.CLI)
	}

	got, err := svc.Get()
	if err != nil {
		t.Fatal(err)
	}
	if got.Paths.CharacterPath != "/chars" || got.Paths.DataPath != "/data" || !got.Characters.InitDefaultVoice {
		t.Fatalf("persisted = %+v", got)
	}
	if got.CLI.Interactive() || !got.CLI.Accessibility() {
		t.Fatalf("persisted cli = interactive=%v accessibility=%v",
			got.CLI.Interactive(), got.CLI.Accessibility())
	}
}

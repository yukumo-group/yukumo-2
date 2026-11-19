package api_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"yukumo-script/core/api"
	"yukumo-script/core/internal/store"
)

func TestNewWithMemoryPorts(t *testing.T) {
	cfgStore := store.NewMemory()
	charStore := store.NewMemory()
	app, err := api.New(api.Deps{
		ConfigStore:    cfgStore,
		CharacterStore: charStore,
		Artifacts:      store.NewMemory(),
		Defaults:       api.DefaultConfig(),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = app.Close(context.Background()) })

	cfg, err := app.Config()
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Synth.AlwaysWrite || !cfg.Synth.ReuseCache {
		t.Fatalf("defaults: %+v", cfg.Synth)
	}
	if err := app.SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}
	if err := app.AddCharacter(api.Character{Name: "remilia", Description: "example"}); err != nil {
		t.Fatal(err)
	}
	list, err := app.ListCharacters()
	if err != nil || len(list) != 1 || list[0].Name != "remilia" {
		t.Fatalf("list = %+v, %v", list, err)
	}
}

func TestNewRequiresStores(t *testing.T) {
	if _, err := api.New(api.Deps{CharacterStore: store.NewMemory()}); err == nil {
		t.Fatal("expected config store error")
	}
	if _, err := api.New(api.Deps{ConfigStore: store.NewMemory()}); err == nil {
		t.Fatal("expected character store error")
	}
}

func TestCharacterJSONContract(t *testing.T) {
	data, err := json.Marshal(api.Character{
		Name:        "remilia",
		Description: "example",
		Voices: []api.Voice{{
			Name:   "default",
			Speed:  100,
			Params: []api.SynthParam{{Version: api.Version1, Preset: 0}},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	for _, want := range []string{
		`"name":"remilia"`,
		`"description":"example"`,
		`"voices":[`,
		`"params":[`,
		`"version":1`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("json %s missing %s", got, want)
		}
	}
	if strings.Contains(got, `"Name"`) || strings.Contains(got, `"PhontData"`) {
		t.Fatalf("json should use camelCase: %s", got)
	}
}

package action_test

import (
	"context"
	"path/filepath"
	"testing"

	"yukumo-script/cli/action"
	"yukumo-script/core/api"
	"yukumo-script/core/api/cliopen"
)

func newApp(t *testing.T) *api.App {
	t.Helper()
	dir := t.TempDir()
	app, err := cliopen.Open(cliopen.Options{
		ConfigDir:     dir,
		CharacterPath: filepath.Join(dir, "characters"),
		OutputPath:    filepath.Join(dir, "output"),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = app.Close(context.Background()) })
	cfg, err := app.Config()
	if err != nil {
		t.Fatal(err)
	}
	cfg.Paths.CharacterPath = filepath.Join(dir, "characters")
	cfg.Paths.OutputPath = filepath.Join(dir, "output")
	if err := app.SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}
	return app
}

func TestConfigAndCharacterUseCases(t *testing.T) {
	ctx := context.Background()
	app := newApp(t)

	if _, err := action.ApplyConfigEdit(ctx, app, "synth.reuse_cache", "false", true); err != nil {
		t.Fatal(err)
	}
	got, err := action.GetConfigField(ctx, app, "synth.reuse_cache")
	if err != nil {
		t.Fatal(err)
	}
	if got != "false" {
		t.Fatalf("reuse_cache = %q, want false", got)
	}

	created, err := action.AddCharacter(ctx, app, "remilia", "example")
	if err != nil {
		t.Fatal(err)
	}
	if created.Name != "remilia" {
		t.Fatalf("created name = %q", created.Name)
	}
	updated, result, err := action.ApplyCharacterEdit(ctx, app, "remilia", "description", "updated", true)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Description != "updated" || result.Value != "updated" {
		t.Fatalf("updated=%+v result=%+v", updated, result)
	}
	if _, err := action.GetCharacter(ctx, app, "remilia"); err != nil {
		t.Fatal(err)
	}
	if err := action.DeleteCharacter(ctx, app, "remilia"); err != nil {
		t.Fatal(err)
	}
}
